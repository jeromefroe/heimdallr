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

// Package controller defines a controller for heimdallr checks.
package controller

import (
	"fmt"

	"github.com/jeromefroe/heimdallr/pkg/apis/heimdallr/v1alpha1"
	"github.com/jeromefroe/heimdallr/pkg/pingdom"

	"go.uber.org/zap"
)

// Controller watches for heimdallr checks and translates them into calls to Pingdom.
type Controller struct {
	client pingdomClient
	logger *zap.Logger
}

// New creates a new controller.
func New(client *pingdom.Client, logger *zap.Logger) *Controller {
	return new(client, logger)
}

func new(client pingdomClient, logger *zap.Logger) *Controller {
	return &Controller{
		client: client,
		logger: logger,
	}
}

// OnAdd handles new HTTP checks.
func (c *Controller) OnAdd(obj interface{}) {
	chk, ok := obj.(*v1alpha1.HTTPCheck)
	if !ok {
		c.logUnexpected("OnAdd", obj)
		return
	}

	err := c.client.UpdateHTTPCheck(*chk)
	if err != nil {
		chk.Status.State = fmt.Sprintf("failed to create check: %v", err)
		c.logger.Error("unexpected error encountered adding check", zap.Error(err))
		return
	}

	chk.Status.State = "successfully created check"
	c.logger.Info("OnAdd successful", zap.String("name", chk.Name))
}

// OnUpdate handles updates HTTP checks.
func (c *Controller) OnUpdate(oldObj, newObj interface{}) {
	oldChk, ok := oldObj.(*v1alpha1.HTTPCheck)
	if !ok {
		c.logUnexpected("OnUpdate", oldObj)
	}

	newChk, ok := newObj.(*v1alpha1.HTTPCheck)
	if !ok {
		c.logUnexpected("OnUpdate", newObj)
	}

	if oldChk == newChk {
		c.logger.Info(
			"new and old checks are identical so no further work is required",
			zap.String("name", newChk.Name),
		)
		return
	}

	err := c.client.UpdateHTTPCheck(*newChk)
	if err != nil {
		newChk.Status.State = fmt.Sprintf("failed to update check: %v", err)
		c.logger.Error("unexpected error encountered updating check", zap.Error(err))
		return
	}

	newChk.Status.State = "successfully updated check"
	c.logger.Info("OnUpdate successful", zap.String("name", newChk.Name))
}

// OnDelete handles deleted HTTP checks.
func (c *Controller) OnDelete(obj interface{}) {
	chk, ok := obj.(*v1alpha1.HTTPCheck)
	if !ok {
		c.logUnexpected("OnDelete", obj)
		return
	}

	err := c.client.DeleteHTTPCheck(*chk)
	if err != nil {
		chk.Status.State = fmt.Sprintf("failed to delete check: %v", err)
		c.logger.Error("unexpected error encountered deleting check", zap.Error(err))
		return
	}

	chk.Status.State = "successfully deleted check"
	c.logger.Info("OnDelete successful", zap.String("name", chk.Name))
}

func (c Controller) logUnexpected(fn string, obj interface{}) {
	c.logger.Error(
		"unexpected object received",
		zap.String("function", fn),
		zap.Any("object", obj),
		zap.String("type", fmt.Sprintf("%T", obj)),
	)
}
