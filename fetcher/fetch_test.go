package fetcher

import (
	"errors"
	"golang-patterns/internal/sleep"
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
	type functions struct {
		producerFn ProducerFn
		consumerFn func(*sampleStruct) ConsumerFn
	}
	type fields struct {
		functions []functions
	}
	type args struct {
		target *sampleStruct
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *sampleStruct
		wantErr bool
	}{
		{
			name: "Should bind properly",
			fields: fields{
				functions: []functions{
					{
						producerFn: func() (interface{}, error) {
							return "some value", nil
						},
						consumerFn: func(target *sampleStruct) ConsumerFn {
							return func(source interface{}) {
								target.SomeField = source.(string)
							}
						},
					},
					{
						producerFn: func() (interface{}, error) {
							return 15, nil

						},
						consumerFn: func(target *sampleStruct) ConsumerFn {
							return func(source interface{}) {
								target.AnotherField = source.(int)
							}
						},
					},
				},
			},
			args: args{
				target: &sampleStruct{},
			},
			want: &sampleStruct{
				SomeField:    "some value",
				AnotherField: 15,
			},
			wantErr: false,
		},
		{
			name: "Should handle error properly",
			fields: fields{
				functions: []functions{
					{
						producerFn: func() (interface{}, error) {
							sleep.Random()
							return "", errors.New("some error")
						},
						consumerFn: func(target *sampleStruct) ConsumerFn {
							return func(source interface{}) {
								sleep.Random()
								target.SomeField = source.(string)
							}
						},
					},
					{
						producerFn: func() (interface{}, error) {
							sleep.Random()
							return 15, nil

						},
						consumerFn: func(target *sampleStruct) ConsumerFn {
							return func(source interface{}) {
								sleep.Random()
								target.AnotherField = source.(int)
							}
						},
					},
				},
			},
			args: args{
				target: &sampleStruct{},
			},
			want:    &sampleStruct{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetcher := New()

			for _, f := range tt.fields.functions {
				fetcher.With(f.producerFn, f.consumerFn(tt.args.target))
			}

			if err := fetcher.Fetch(tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.args.target, tt.want) {
				t.Errorf("Fetch() = %v, want %v", tt.args.target, tt.want)
			}
		})
	}
}
