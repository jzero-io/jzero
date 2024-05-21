package gen

import (
	"reflect"
	"testing"
)

func Test_separateCommonElements(t *testing.T) {
	type args struct {
		arrays [][]string
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "test",
			args: args{
				arrays: [][]string{
					{"1", "2", "3"},
					{"3", "5", "6"},
					{"7", "8", "9"},
				},
			},
			want: [][]string{
				{"1", "2"},
				{"5", "6"},
				{"7", "8", "9"},
				{"3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := separateCommonElements(tt.args.arrays...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("separateCommonElements() = %v, want %v", got, tt.want)
			}
		})
	}
}
