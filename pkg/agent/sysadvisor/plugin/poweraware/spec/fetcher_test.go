/*
Copyright 2022 The Katalyst Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package spec

import (
	"context"
	"reflect"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kubewharf/katalyst-core/pkg/metaserver/agent/node"
)

type mockNodeFetcher struct {
	node.NodeFetcher
}

func (m mockNodeFetcher) GetNode(ctx context.Context) (*v1.Node, error) {
	return &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"foo/power-alert":       "s0",
				"foo/power-budget":      "128",
				"foo/power-internal-op": "8",
				"foo/power-alert-time":  "2024-06-01T19:15:58Z",
			},
		},
	}, nil
}

var _ node.NodeFetcher = &mockNodeFetcher{}

func Test_specFetcherByNodeAnnotation_GetPowerSpec(t *testing.T) {
	t.Parallel()
	type fields struct {
		nodeFetcher node.NodeFetcher
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *PowerSpec
		wantErr bool
	}{
		{
			name: "happy path no error",
			fields: fields{
				nodeFetcher: &mockNodeFetcher{},
			},
			args: args{
				ctx: context.TODO(),
			},
			want: &PowerSpec{
				Alert:      PowerAlertS0,
				Budget:     128,
				InternalOp: InternalOpNoop,
				AlertTime:  time.Date(2024, time.June, 1, 19, 15, 58, 0, time.UTC),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := specFetcherByNodeAnnotation{
				nodeFetcher:         tt.fields.nodeFetcher,
				annotationKeyPrefix: "foo",
			}
			got, err := s.GetPowerSpec(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPowerSpec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPowerSpec() got = %v, want %v", got, tt.want)
			}
		})
	}
}
