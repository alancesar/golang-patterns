package fetcher

import (
	"reflect"
	"sync"
	"testing"
)

type sampleStruct struct {
	sync.Mutex
	SomeField    string
	AnotherField int
}

func TestFetcher_Fetch(t *testing.T) {
	type fields struct {
		dispatchers []<-chan incoming
	}
	type args struct {
		locker *sampleStruct
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *sampleStruct
	}{
		{
			name:   "Should bind properly",
			fields: fields{},
			args: args{
				locker: &sampleStruct{},
			},
			want: &sampleStruct{
				SomeField:    "some value",
				AnotherField: 15,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			New().With(func() interface{} {
				return "some value"
			}, func(value interface{}) {
				tt.args.locker.SomeField = value.(string)
			}).With(func() interface{} {
				return 15
			}, func(value interface{}) {
				tt.args.locker.AnotherField = value.(int)
			}).Fetch(tt.args.locker)

			if !reflect.DeepEqual(tt.args.locker, tt.want) {
				t.Errorf("Fetch() = %v, want %v", tt.args.locker, tt.want)
			}
		})
	}
}
