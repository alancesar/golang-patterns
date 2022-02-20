package newfetcher

import (
	"golang.org/x/sync/errgroup"
	"reflect"
	"sync"
	"testing"
)

type testData struct {
	sync.Mutex
	Value string
}

func TestNew(t *testing.T) {
	type args struct {
		provider ProviderFn[string]
		callback CallbackFn[*testData, string]
	}
	tests := []struct {
		name    string
		args    args
		want    *testData
		wantErr bool
	}{
		{
			name: "Should fetch all properties properly",
			args: args{
				provider: func() (string, error) {
					return "Some value", nil
				},
				callback: func(target *testData, source string) error {
					target.Value = source
					return nil
				},
			},
			want: &testData{
				Value: "Some value",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var group errgroup.Group
			target := testData{}

			group.Go(func() error {
				return New[*testData](&target, tt.args.provider, tt.args.callback)
			})

			if err := group.Wait(); (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(&target, tt.want) {
				t.Errorf("Work() = %v, want %v", &target, tt.want)
			}
		})
	}
}
