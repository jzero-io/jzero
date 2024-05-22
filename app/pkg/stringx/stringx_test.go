package stringx

import "testing"

func TestToCamel(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				s: "ping-pong",
			},
			want: "pingPong",
		},

		{
			name: "test2",
			args: args{
				s: "ping/pong",
			},
			want: "pingPong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToCamel(tt.args.s); got != tt.want {
				t.Errorf("ToCamel() = %v, want %v", got, tt.want)
			}
		})
	}
}
