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

package controller

import (
	"errors"
	"testing"

	"github.com/jeromefroe/heimdallr/pkg/apis/heimdallr/v1alpha1"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestOnAdd(t *testing.T) {
	mCtrl := gomock.NewController(t)
	defer mCtrl.Finish()

	check := v1alpha1.HTTPCheck{
		ObjectMeta: metav1.ObjectMeta{
			Name: "check",
		},
	}

	cli := NewMockpingdomClient(mCtrl)
	cli.EXPECT().UpdateHTTPCheck(check).Return(nil)

	ctrl := new(cli, zap.NewNop())
	ctrl.OnAdd(&check)

	assert.Contains(t, check.Status.State, "success")
}

func TestOnAddError(t *testing.T) {
	mCtrl := gomock.NewController(t)
	defer mCtrl.Finish()

	check := v1alpha1.HTTPCheck{
		ObjectMeta: metav1.ObjectMeta{
			Name: "check",
		},
	}

	cli := NewMockpingdomClient(mCtrl)
	cli.EXPECT().UpdateHTTPCheck(check).Return(errors.New("bad requests"))

	ctrl := new(cli, zap.NewNop())
	ctrl.OnAdd(&check)

	assert.Contains(t, check.Status.State, "fail")
}

func TestOnUpdate(t *testing.T) {
	mCtrl := gomock.NewController(t)
	defer mCtrl.Finish()

	var (
		newCheck = v1alpha1.HTTPCheck{
			ObjectMeta: metav1.ObjectMeta{
				Name: "check",
			},
			Spec: v1alpha1.HTTPCheckSpec{
				EnableTLS: false,
			},
		}
		oldCheck = v1alpha1.HTTPCheck{
			ObjectMeta: metav1.ObjectMeta{
				Name: "check",
			},
			Spec: v1alpha1.HTTPCheckSpec{
				EnableTLS: true,
			},
		}
	)

	cli := NewMockpingdomClient(mCtrl)
	cli.EXPECT().UpdateHTTPCheck(newCheck).Return(nil)

	ctrl := new(cli, zap.NewNop())
	ctrl.OnUpdate(&oldCheck, &newCheck)

	assert.Contains(t, newCheck.Status.State, "success")
}

func TestOnUpdateError(t *testing.T) {
	mCtrl := gomock.NewController(t)
	defer mCtrl.Finish()

	var (
		newCheck = v1alpha1.HTTPCheck{
			ObjectMeta: metav1.ObjectMeta{
				Name: "check",
			},
			Spec: v1alpha1.HTTPCheckSpec{
				EnableTLS: false,
			},
		}
		oldCheck = v1alpha1.HTTPCheck{
			ObjectMeta: metav1.ObjectMeta{
				Name: "check",
			},
			Spec: v1alpha1.HTTPCheckSpec{
				EnableTLS: true,
			},
		}
	)

	cli := NewMockpingdomClient(mCtrl)
	cli.EXPECT().UpdateHTTPCheck(newCheck).Return(errors.New("bad request"))

	ctrl := new(cli, zap.NewNop())
	ctrl.OnUpdate(&oldCheck, &newCheck)

	assert.Contains(t, newCheck.Status.State, "fail")
}

func TestOnDelete(t *testing.T) {
	mCtrl := gomock.NewController(t)
	defer mCtrl.Finish()

	check := v1alpha1.HTTPCheck{
		ObjectMeta: metav1.ObjectMeta{
			Name: "check",
		},
	}

	cli := NewMockpingdomClient(mCtrl)
	cli.EXPECT().DeleteHTTPCheck(check).Return(nil)

	ctrl := new(cli, zap.NewNop())
	ctrl.OnDelete(&check)

	assert.Contains(t, check.Status.State, "success")
}

func TestOnDeleteError(t *testing.T) {
	mCtrl := gomock.NewController(t)
	defer mCtrl.Finish()

	check := v1alpha1.HTTPCheck{
		ObjectMeta: metav1.ObjectMeta{
			Name: "check",
		},
	}

	cli := NewMockpingdomClient(mCtrl)
	cli.EXPECT().DeleteHTTPCheck(check).Return(errors.New("bad requests"))

	ctrl := new(cli, zap.NewNop())
	ctrl.OnDelete(&check)

	assert.Contains(t, check.Status.State, "fail")
}
