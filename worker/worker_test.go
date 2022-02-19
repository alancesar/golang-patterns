package worker

import (
	"context"
	"math/rand"
	"reflect"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorker_Work(t *testing.T) {
	type fields struct {
		fn     func(counter *uint64) Fn[int]
		buffer int
	}
	type args struct {
		ctx     context.Context
		items   []int
		counter uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint64
	}{
		{
			name: "Should execute all jobs properly",
			fields: fields{
				fn: func(counter *uint64) Fn[int] {
					return func(ctx context.Context, input int) {
						atomic.AddUint64(counter, 1)
						time.Sleep(time.Millisecond * time.Duration(input) * 100)
					}
				},
				buffer: 2,
			},
			args: args{
				ctx:     context.Background(),
				items:   []int{rand.Intn(50), rand.Intn(40), rand.Intn(30), rand.Intn(20), rand.Intn(10)},
				counter: 0,
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter := tt.args.counter
			w := New[int](tt.fields.fn(&counter), tt.fields.buffer)
			w.Work(tt.args.ctx, tt.args.items)

			if !reflect.DeepEqual(counter, tt.want) {
				t.Errorf("Work() = %v, want %v", tt.args.counter, tt.want)
			}
		})
	}
}
