package model

import (
	"reflect"
	"testing"
)

func TestGenBoardFromVals(t *testing.T) {
	type args struct {
		vals []byte
		idxs []int
	}
	tests := []struct {
		name      string
		args      args
		wantBoard []byte
		wantErr   bool
	}{
		{name: "Basic positive test", args: args{vals: []byte{0, 0, 1, 2, 3}, idxs: []int{1, 0, 2, 4, 3}}, wantBoard: []byte{0, 0, 1, 3, 2}, wantErr: false},
		{name: "Mismatch lengths", args: args{vals: []byte{0, 0, 1, 2}, idxs: []int{1, 0, 2, 4, 3}}, wantBoard: nil, wantErr: true},
		{name: "Index too large", args: args{vals: []byte{0, 0, 1, 2, 3}, idxs: []int{1, 0, 2, 5, 3}}, wantBoard: nil, wantErr: true},
		{name: "Index negative", args: args{vals: []byte{0, 0, 1, 2, 3}, idxs: []int{1, 0, -2, 4, 3}}, wantBoard: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBoard, err := GenBoardFromVals(tt.args.vals, tt.args.idxs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenBoardFromVals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBoard, tt.wantBoard) {
				t.Errorf("GenBoardFromVals() = %v, want %v", gotBoard, tt.wantBoard)
			}
		})
	}
}

func TestGenBoardForLevel(t *testing.T) {
	type args struct {
		level int
	}
	tests := []struct {
		name    string
		args    args
		wantMin int
		wantMax int
	}{
		{name: "Level 1", args: args{level: 1}, wantMin: 24, wantMax: 144},
		{name: "Level 2", args: args{level: 2}, wantMin: 144, wantMax: 432},
		{name: "Level 3", args: args{level: 3}, wantMin: 432, wantMax: 864},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board := GenBoardForLevel(tt.args.level)
			game := Game{Board: board}
			points := game.GetMaxScore()
			if points > tt.wantMax || points < tt.wantMin {
				t.Errorf("GenBoardForLevel() = %v, want min:%v, max:%v", points, tt.wantMin, tt.wantMax)
			}
		})
	}
}
