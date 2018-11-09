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

// +build integration

package pingdom

import (
	"os"
	"testing"

	"github.com/jeromefroe/heimdallr/pkg/apis/heimdallr/v1alpha1"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// To run the integration test:
//   PINGDOM_USERNAME=<username> PINGDOM_PASSWORD=<password> PINGDOM_APPKEY=<key> \
//     go test ./... -run ^TestIntegration$ -v -tags=integration

func TestIntegration(t *testing.T) {
	var (
		username = os.Getenv("PINGDOM_USERNAME")
		password = os.Getenv("PINGDOM_PASSWORD")
		appkey   = os.Getenv("PINGDOM_APPKEY")

		meta = metav1.ObjectMeta{
			Name:      "test",
			Namespace: "heimdallr",
		}
		spec = v1alpha1.HTTPCheckSpec{
			Hostname:           "froe.io",
			IntervalMinutes:    1,
			TriggerThreshold:   1,
			RetriggerThreshold: 10,
			NotifyWhenBackup:   true,
			EnableTLS:          true,
		}
		check = v1alpha1.HTTPCheck{
			ObjectMeta: meta,
			Spec:       spec,
		}
	)

	// This isn't very clean, but we override heimdallrTag for the test so we can
	// know which checks the test creates.
	heimdallrTag = "heimdallr-test"

	client, err := New(username, password, appkey)
	require.NoError(t, err)
	require.Len(t, client.httpChecks, 0)

	err = client.UpdateHTTPCheck(check)
	require.NoError(t, err)

	check.Spec.NotifyWhenBackup = false
	err = client.UpdateHTTPCheck(check)
	require.NoError(t, err)

	err = client.DeleteHTTPCheck(check)
	require.NoError(t, err)
}
