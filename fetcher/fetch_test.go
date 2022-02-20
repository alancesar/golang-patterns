package fetcher

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
		provider ProviderFn
		callback CallbackFn
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
				provider: func() (interface{}, error) {
					return "Some value", nil
				},
				callback: func(target sync.Locker, source interface{}) error {
					target.(*testData).Value = source.(string)
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
				return New(&target, tt.args.provider, tt.args.callback)
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
