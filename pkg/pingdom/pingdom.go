// Copyright (c) 2018 Jerome Froelich
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package pingdom provides a client for interacting with the Pingdom API.
package pingdom

import (
	"fmt"

	"github.com/jeromefroe/heimdallr/pkg/apis/heimdallr/v1alpha1"

	"github.com/russellcardullo/go-pingdom/pingdom"
)

// heimdallrTag is the tag added to every check to indicate that it is managed by heimdallr.
var heimdallrTag = "managed-by-heimdallr"

type httpCheck struct {
	id   int
	name string
	spec v1alpha1.HTTPCheckSpec
}

// Client is a Pingdom API Client.
type Client struct {
	userID     int
	idTag      string
	client     pingdomClient
	httpChecks map[string]httpCheck
}

// New creates a new Pingdom client.
func New(user, password, key string) (*Client, error) {
	var (
		client = pingdom.NewClient(user, password, key)
		shim   = newShimClient(client)
	)
	return new(user, heimdallrTag, shim)
}

func new(user, idTag string, client pingdomClient) (*Client, error) {
	users, err := client.Users().List()
	if err != nil {
		return nil, fmt.Errorf("failed to get list of users for account: %v", err)
	}

	var userID *int
	for _, userResp := range users {
		for _, email := range userResp.Email {
			if email.Address == user {
				userID = &userResp.Id
				break
			}
		}
	}

	if userID == nil {
		return nil, fmt.Errorf("failed to get ID of user %v", user)
	}

	c := &Client{
		userID:     *userID,
		idTag:      idTag,
		client:     client,
		httpChecks: make(map[string]httpCheck),
	}

	return c, c.sync()
}

// Sync fetches the current state of Pingdom.
func (c *Client) sync() error {
	list, err := c.client.Checks().List(map[string]string{
		"tags": c.idTag,
	})
	if err != nil {
		return fmt.Errorf("failed to get current list of heimdallr checks: %v", err)
	}

	for _, cr := range list {
		if check, ok := toHTTPCheck(cr); ok {
			c.httpChecks[cr.Name] = check
		}
	}
	return nil
}

// UpdateHTTPCheck updates an HTTP check, creating it if it does not exist.
func (c *Client) UpdateHTTPCheck(check v1alpha1.HTTPCheck) error {
	name := getName(check)
	pc := pingdom.HttpCheck{
		Name:                     name,
		UserIds:                  []int{c.userID},
		Hostname:                 check.Spec.Hostname,
		Resolution:               check.Spec.IntervalMinutes,
		Encryption:               check.Spec.EnableTLS,
		SendNotificationWhenDown: check.Spec.TriggerThreshold,
		NotifyAgainEvery:         check.Spec.RetriggerThreshold,
		NotifyWhenBackup:         check.Spec.NotifyWhenBackup,
		Tags:                     heimdallrTag,
	}

	hc, ok := c.httpChecks[name]
	if ok {
		_, err := c.client.Checks().Update(hc.id, &pc)
		if err != nil {
			return fmt.Errorf("failed to update check: %v", err)
		}
		hc.spec = check.Spec
	} else {
		res, err := c.client.Checks().Create(&pc)
		if err != nil {
			return fmt.Errorf("failed to create check: %v", err)
		}
		hc = httpCheck{
			id:   res.ID,
			name: name,
			spec: check.Spec,
		}
	}

	c.httpChecks[name] = hc
	return nil
}

// DeleteHTTPCheck deletes an HTTP check.
func (c *Client) DeleteHTTPCheck(check v1alpha1.HTTPCheck) error {
	name := getName(check)
	hc, exists := c.httpChecks[name]
	if !exists {
		return nil
	}

	_, err := c.client.Checks().Delete(hc.id)
	if err != nil {
		return fmt.Errorf("failed to delete check: %v", err)
	}

	delete(c.httpChecks, name)
	return nil
}

func toHTTPCheck(c pingdom.CheckResponse) (httpCheck, bool) {
	if c.Type.HTTP == nil {
		return httpCheck{}, false
	}

	return httpCheck{
		id:   c.ID,
		name: c.Name,
		spec: v1alpha1.HTTPCheckSpec{
			Hostname:           c.Hostname,
			IntervalMinutes:    c.Resolution,
			TriggerThreshold:   c.SendNotificationWhenDown,
			RetriggerThreshold: c.NotifyAgainEvery,
			NotifyWhenBackup:   c.NotifyWhenBackup,
			EnableTLS:          c.Type.HTTP.Encryption,
		},
	}, true
}

func getName(check v1alpha1.HTTPCheck) string {
	ns := check.ObjectMeta.Namespace
	if ns == "" {
		ns = "default"
	}

	return fmt.Sprintf("%s/%s", ns, check.ObjectMeta.Name)
}
