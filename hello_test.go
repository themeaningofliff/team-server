package main

import (
	"testing"
)

func Test_add(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name       string
		args       args
		wantResult int
	}{
		{"Sample Test", args{1, 3}, 4},
		{"Zero Test", args{0, 0}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := add(tt.args.x, tt.args.y); gotResult != tt.wantResult {
				t.Errorf("add() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

// func Test_main(t *testing.T) {
// 	tests := []struct {
// 		name string
// 	}{
// 	// TODO: Add test cases.
// 	}
// 	for range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			main()
// 		})
// 	}
// }
