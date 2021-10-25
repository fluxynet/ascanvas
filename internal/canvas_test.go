package internal

import (
	"reflect"
	"testing"

	"github.com/fluxynet/ascanvas"
)

func TestCanvasFromText(t *testing.T) {
	type args struct {
		id   string
		name string
		s    string
	}
	tests := []struct {
		name     string
		args     args
		want     *ascanvas.Canvas
		wantGrid [][]string
	}{
		{
			name: "canvas 1",
			args: args{
				id:   "1",
				name: "canvas 1",
				s: `
      
 XXXX 
 X  X 
 X  X 
 XXXX `},
			want: &ascanvas.Canvas{
				Id:      "1",
				Name:    "canvas 1",
				Content: "       XXXX  X  X  X  X  XXXX ",
				Width:   6,
				Height:  5,
			},
			wantGrid: [][]string{
				{" ", " ", " ", " ", " ", " "},
				{" ", "X", "X", "X", "X", " "},
				{" ", "X", " ", " ", "X", " "},
				{" ", "X", " ", " ", "X", " "},
				{" ", "X", "X", "X", "X", " "},
			},
		},
		{
			name: "fixture 1.1",
			args: args{
				id:   "1",
				name: "fixture 1.1",
				s: `
                        
   @@@@@                
   @XXX@  XXXXXXXXXXXXXX
   @@@@@  XOOOOOOOOOOOOX
          XOOOOOOOOOOOOX
          XOOOOOOOOOOOOX
          XOOOOOOOOOOOOX
          XXXXXXXXXXXXXX`,
			},
			want: &ascanvas.Canvas{
				Id:      "1",
				Name:    "fixture 1.1",
				Content: `                           @@@@@                   @XXX@  XXXXXXXXXXXXXX   @@@@@  XOOOOOOOOOOOOX          XOOOOOOOOOOOOX          XOOOOOOOOOOOOX          XOOOOOOOOOOOOX          XXXXXXXXXXXXXX`,
				Width:   24,
				Height:  8,
			},
			wantGrid: [][]string{
				{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
				{" ", " ", " ", "@", "@", "@", "@", "@", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " ", " "},
				{" ", " ", " ", "@", "X", "X", "X", "@", " ", " ", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
				{" ", " ", " ", "@", "@", "@", "@", "@", " ", " ", "X", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "X"},
				{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", "X", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "X"},
				{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", "X", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "X"},
				{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", "X", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "O", "X"},
				{" ", " ", " ", " ", " ", " ", " ", " ", " ", " ", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X", "X"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanvasFromText(tt.args.id, tt.args.name, tt.args.s)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"got = \n%s want = %s",
					got.String(),
					tt.want.String(),
				)
			}

			if !reflect.DeepEqual(got.AsGrid(), tt.wantGrid) {
				t.Errorf("grid: \ngot  = %v\nwant = %v\n", got.AsGrid(), tt.wantGrid)
				return
			}
		})
	}
}
