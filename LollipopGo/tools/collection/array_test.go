package collection

import (
	"reflect"
	"testing"
)

func TestDeleteInt32s(t *testing.T) {
	type args struct {
		array []int32
		elem  []int32
	}
	tests := []struct {
		name string
		args args
		want []int32
	}{
		{"", args{[]int32{1, 2, 2, 3, 3, 5, 7}, []int32{2, 3, 3}}, []int32{1, 5, 7}},
		{"", args{[]int32{1, 2, 3, 4, 5, 6}, []int32{1, 2, 3}}, []int32{4, 5, 6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeleteInt32s(tt.args.array, tt.args.elem...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteInt32s() = %v, want %v", got, tt.want)
			}
		})
	}
}
