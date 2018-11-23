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

package pingdom

import (
	"testing"

	"github.com/jeromefroe/heimdallr/pkg/apis/heimdallr/v1alpha1"

	"github.com/golang/mock/gomock"
	"github.com/russellcardullo/go-pingdom/pingdom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		id     = 42
		user   = "bob@example.come"
		users  = NewMockuserService(ctrl)
		checks = NewMockcheckService(ctrl)
		cli    = NewMockpingdomClient(ctrl)
	)

	users.EXPECT().List().Return([]pingdom.UsersResponse{
		{
			Id: id,
			Email: []pingdom.UserEmailResponse{
				{
					Address: user,
				},
			},
		},
	}, nil)

	checks.EXPECT().
		List(map[string]string{"tags": heimdallrTag, "include_tags": "true"}).
		Return([]pingdom.CheckResponse{
			{
				ID:   71,
				Name: "default/foo",
				Tags: []pingdom.CheckResponseTag{
					{
						Name: heimdallrTag,
					},
				},
			},
		}, nil)

	cli.EXPECT().Users().Return(users)
	cli.EXPECT().Checks().Return(checks)

	client, err := new(user, cli, zap.NewNop())
	require.NoError(t, err)
	assert.Equal(t, id, client.userID)

	assert.Len(t, client.httpChecks, 1)
}

func TestSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		checks = NewMockcheckService(ctrl)
		cli    = NewMockpingdomClient(ctrl)
	)

	checks.EXPECT().
		List(map[string]string{"tags": heimdallrTag, "include_tags": "true"}).
		Return([]pingdom.CheckResponse{
			{
				ID:                       71,
				Name:                     "default/foo",
				Hostname:                 "foo.io",
				Resolution:               10,
				SendNotificationWhenDown: 3,
				NotifyAgainEvery:         30,
				NotifyWhenBackup:         true,
				Tags: []pingdom.CheckResponseTag{
					{
						Name: heimdallrTag,
					},
				},
			},
			{
				ID:                       82,
				Name:                     "other/bar",
				Hostname:                 "bar.com",
				Resolution:               5,
				SendNotificationWhenDown: 2,
				NotifyAgainEvery:         8,
				NotifyWhenBackup:         false,
				Tags: []pingdom.CheckResponseTag{
					{
						Name: heimdallrTag,
					},
					{
						Name: tlsEnabledTag,
					},
				},
			},
		}, nil)

	cli.EXPECT().Checks().Return(checks)

	client := Client{
		client:     cli,
		httpChecks: make(map[string]httpCheck),
		logger:     zap.NewNop(),
	}
	require.NoError(t, client.sync())

	assert.Len(t, client.httpChecks, 2)

	expected := httpCheck{
		id:   71,
		name: "default/foo",
		spec: v1alpha1.HTTPCheckSpec{
			Hostname:           "foo.io",
			IntervalMinutes:    10,
			TriggerThreshold:   3,
			RetriggerThreshold: 30,
			NotifyWhenBackup:   true,
			EnableTLS:          false,
		},
	}
	assert.Equal(t, expected, client.httpChecks["default/foo"])

	expected = httpCheck{
		id:   82,
		name: "other/bar",
		spec: v1alpha1.HTTPCheckSpec{
			Hostname:           "bar.com",
			IntervalMinutes:    5,
			TriggerThreshold:   2,
			RetriggerThreshold: 8,
			NotifyWhenBackup:   false,
			EnableTLS:          true,
		},
	}
	assert.Equal(t, expected, client.httpChecks["other/bar"])
}

func TestUpdateHTTPCheckWithNewCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		id   = 42
		spec = v1alpha1.HTTPCheckSpec{
			Hostname:           "foo.io",
			IntervalMinutes:    10,
			TriggerThreshold:   3,
			RetriggerThreshold: 30,
			NotifyWhenBackup:   true,
			EnableTLS:          false,
		}
		check = v1alpha1.HTTPCheck{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "other",
			},
			Spec: spec,
		}
		name = "other/foo"

		checks = NewMockcheckService(ctrl)
		cli    = NewMockpingdomClient(ctrl)
	)

	checks.EXPECT().Create(gomock.Any()).Return(&pingdom.CheckResponse{
		ID: id,
	}, nil)

	cli.EXPECT().Checks().Return(checks)

	client := Client{
		client:     cli,
		httpChecks: map[string]httpCheck{},
		logger:     zap.NewNop(),
	}

	err := client.UpdateHTTPCheck(check)
	require.NoError(t, err)
	assert.Len(t, client.httpChecks, 1)

	expected := httpCheck{
		id:   id,
		name: name,
		spec: spec,
	}
	assert.Equal(t, expected, client.httpChecks[name])
}

func TestUpdateHTTPCheckWithExistingCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		id   = 42
		spec = v1alpha1.HTTPCheckSpec{
			Hostname:           "foo.io",
			IntervalMinutes:    10,
			TriggerThreshold:   3,
			RetriggerThreshold: 30,
			NotifyWhenBackup:   true,
			EnableTLS:          false,
		}
		check = v1alpha1.HTTPCheck{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "other",
			},
			Spec: spec,
		}
		name = "other/foo"

		checks = NewMockcheckService(ctrl)
		cli    = NewMockpingdomClient(ctrl)
	)

	checks.EXPECT().Update(id, gomock.Any())
	cli.EXPECT().Checks().Return(checks)

	client := Client{
		client: cli,
		httpChecks: map[string]httpCheck{
			name: {
				id:   id,
				name: name,
			},
		},
		logger: zap.NewNop(),
	}

	err := client.UpdateHTTPCheck(check)
	require.NoError(t, err)
	assert.Len(t, client.httpChecks, 1)

	expected := httpCheck{
		id:   id,
		name: name,
		spec: spec,
	}
	assert.Equal(t, expected, client.httpChecks[name])
}

func TestDeleteHTTPCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var (
		id     = 42
		checks = NewMockcheckService(ctrl)
		cli    = NewMockpingdomClient(ctrl)
	)

	checks.EXPECT().Delete(id)

	cli.EXPECT().Checks().Return(checks)

	client := Client{
		client: cli,
		httpChecks: map[string]httpCheck{
			"default/foo": {
				id: id,
			},
		},
		logger: zap.NewNop(),
	}

	check := v1alpha1.HTTPCheck{
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo",
		},
	}
	err := client.DeleteHTTPCheck(check)
	require.NoError(t, err)
	assert.Len(t, client.httpChecks, 0)
}
