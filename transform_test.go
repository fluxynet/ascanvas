package ascanvas_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/fluxynet/ascanvas"
	"github.com/fluxynet/ascanvas/internal"
)

func TestTransformRectangle(t *testing.T) {
	tests := []struct {
		name    string
		canvas  *ascanvas.Canvas
		args    ascanvas.TransformRectangleArgs
		want    *ascanvas.Canvas
		wantErr error
	}{
		{
			name: "empty outline only",
			canvas: internal.CanvasFromText("1", "canvas 1", `
     
     
     
     
     `),
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 1,
					Y: 1,
				},
				Width:   4,
				Height:  4,
				Fill:    "",
				Outline: "X",
			},
			want: internal.CanvasFromText("1", "canvas 1", `
     
 XXXX
 X  X
 X  X
 XXXX`),
			wantErr: nil,
		},
		{
			name: "empty fill only",
			canvas: internal.CanvasFromText("2", "canvas 2", `
     
     
     
     
     `),
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 1,
					Y: 1,
				},
				Width:   4,
				Height:  4,
				Fill:    "X",
				Outline: "",
			},
			want: internal.CanvasFromText("2", "canvas 2", `
     
 XXXX
 XXXX
 XXXX
 XXXX`,
			),
			wantErr: nil,
		},
		{
			name: "empty outline and fill",
			canvas: internal.CanvasFromText("3", "canvas 3", `
     
     
     
     
     `,
			),
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 1,
					Y: 1,
				},
				Width:   4,
				Height:  4,
				Fill:    "O",
				Outline: "X",
			},
			want: internal.CanvasFromText("3", "canvas 3", `
     
 XXXX
 XOOX
 XOOX
 XXXX`,
			),
			wantErr: nil,
		},
		{
			name: "non-empty outline only",
			canvas: internal.CanvasFromText("1", "canvas 1", `
    Z
   > 
  ?  
 !   
.    `),
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 1,
					Y: 1,
				},
				Width:   4,
				Height:  4,
				Fill:    "",
				Outline: "X",
			},
			want: internal.CanvasFromText("1", "canvas 1", `
    Z
 XXXX
 X? X
 X  X
.XXXX`),
			wantErr: nil,
		},
		{
			name: "non-empty fill only",
			canvas: internal.CanvasFromText("2", "canvas 2", `
W    
@    
  @  
  @  
     `),
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 1,
					Y: 1,
				},
				Width:   4,
				Height:  4,
				Fill:    "X",
				Outline: "",
			},
			want: internal.CanvasFromText("2", "canvas 2", `
W    
@XXXX
 XXXX
 XXXX
 XXXX`),
			wantErr: nil,
		},
		{
			name: "non-empty outline and fill",
			canvas: internal.CanvasFromText("3", "canvas 3", `
   @*
     
     
   * 
   @ `),
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 1,
					Y: 1,
				},
				Width:   4,
				Height:  4,
				Fill:    "O",
				Outline: "X",
			},
			want: internal.CanvasFromText("3", "canvas 3", `
   @*
 XXXX
 XOOX
 XOOX
 XXXX`),
			wantErr: nil,
		},
		{
			name: "zero height outline only",
			canvas: internal.CanvasFromText("1", "canvas 1", `
    Z
   > 
  ?  
 !   
.    `),
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 1,
					Y: 1,
				},
				Width:   4,
				Height:  0,
				Fill:    "",
				Outline: "X",
			},
			want: internal.CanvasFromText("1", "canvas 1", `
    Z
 XXXX
  ?  
 !   
.    `),
			wantErr: nil,
		},
		{
			name: "zero height fill only",
			canvas: internal.CanvasFromText("2", "canvas 2", `
W    
@    
  @  
  @  
     `),
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 1,
					Y: 1,
				},
				Width:   4,
				Height:  0,
				Fill:    "X",
				Outline: "",
			},
			want: internal.CanvasFromText("2", "canvas 2", `
W    
@XXXX
  @  
  @  
     `),
			wantErr: nil,
		},
		{
			name: "zero height outline and fill",
			canvas: internal.CanvasFromText("3", "canvas 3", `
   @*
     
     
   * 
   @ `),
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 1,
					Y: 1,
				},
				Width:   4,
				Height:  0,
				Fill:    "O",
				Outline: "X",
			},
			want: internal.CanvasFromText("3", "canvas 3", `
   @*
 XXXX
     
   * 
   @ `),
			wantErr: nil,
		},
		{
			name: "zero width outline only",
			canvas: internal.CanvasFromText("1", "canvas 1", `
    Z
   > 
  ?  
 !   
.    `),
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 1,
					Y: 1,
				},
				Width:   0,
				Height:  4,
				Fill:    "",
				Outline: "X",
			},
			want: internal.CanvasFromText("1", "canvas 1", `
    Z
 X > 
 X?  
 X   
.X   `),
			wantErr: nil,
		},
		{
			name: "zero width fill only",
			canvas: internal.CanvasFromText("2", "canvas 2", `
W    
@    
  @  
  @  
 @   `),
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 1,
					Y: 1,
				},
				Width:   0,
				Height:  4,
				Fill:    "X",
				Outline: "",
			},
			want: internal.CanvasFromText("2", "canvas 2", `
W    
@X   
 X@  
 X@  
 X   `),
			wantErr: nil,
		},
		{
			name: "zero width outline and fill",
			canvas: internal.CanvasFromText("3", "canvas 3", `
   @*
     
     
 % * 
   @ `),
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 1,
					Y: 1,
				},
				Width:   0,
				Height:  4,
				Fill:    "O",
				Outline: "X",
			},
			want: internal.CanvasFromText("3", "canvas 3", `
   @*
 X   
 X   
 X * 
 X @ `),
			wantErr: nil,
		},
		{
			name: "fixture 1.1",
			canvas: internal.CanvasFromText("1", "fixture 1.1", `
                        
                        
                        
                        
                        
                        
                        
                        
                        `),
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 3,
					Y: 2,
				},
				Width:   5,
				Height:  3,
				Fill:    "X",
				Outline: "@",
			},
			want: internal.CanvasFromText("1", "fixture 1.1", `
                        
                        
   @@@@@                
   @XXX@                
   @@@@@                
                        
                        
                        
                        `),
		},
		{
			name: "fixture 1.2",
			canvas: internal.CanvasFromText("1", "fixture 1.2", `
                        
                        
   @@@@@                
   @XXX@                
   @@@@@                
                        
                        
                        
                        `),
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 10,
					Y: 3,
				},
				Width:   14,
				Height:  6,
				Fill:    "O",
				Outline: "X",
			},
			want: internal.CanvasFromText("1", "fixture 1.2", `
                        
                        
   @@@@@                
   @XXX@  XXXXXXXXXXXXXX
   @@@@@  XOOOOOOOOOOOOX
          XOOOOOOOOOOOOX
          XOOOOOOOOOOOOX
          XOOOOOOOOOOOOX
          XXXXXXXXXXXXXX`),
		},
		{
			name: "fixture 3.1",
			canvas: internal.CanvasFromText("1", "fixture 3.1", `
                            
                            
                            
                            
                            
                            
                            
                            
                            
                            
                            
                            `),
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 15,
					Y: 0,
				},
				Width:   7,
				Height:  6,
				Fill:    ".",
				Outline: "",
			},
			want: internal.CanvasFromText("1", "fixture 3.1", `
               .......      
               .......      
               .......      
               .......      
               .......      
               .......      
                            
                            
                            
                            
                            
                            `),
		},
		{
			name: "fixture 3.2",
			canvas: internal.CanvasFromText("1", "fixture 3.2", `
              .......       
              .......       
              .......       
              .......       
              .......       
              .......       
                            
                            
                            
                            
                            
                            `),
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 0,
					Y: 3,
				},
				Width:   8,
				Height:  4,
				Fill:    "",
				Outline: "O",
			},
			want: internal.CanvasFromText("1", "fixture 3.2", `
              .......       
              .......       
              .......       
OOOOOOOO      .......       
O      O      .......       
O      O      .......       
OOOOOOOO                    
                            
                            
                            
                            
                            `),
		},
		{
			name: "fixture 3.3",
			canvas: internal.CanvasFromText("1", "fixture 3.3", `
              .......       
              .......       
              .......       
OOOOOOOO      .......       
O      O      .......       
O      O      .......       
OOOOOOOO                    
                            
                            
                            
                            
                            `),
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 5,
					Y: 5,
				},
				Width:   5,
				Height:  3,
				Fill:    "X",
				Outline: "X",
			},
			want: internal.CanvasFromText("1", "fixture 3.3", `
              .......       
              .......       
              .......       
OOOOOOOO      .......       
O      O      .......       
O    XXXXX    .......       
OOOOOXXXXX                  
     XXXXX                  
                            
                            
                            
                            `),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ascanvas.TransformRectangle(tt.canvas, tt.args)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("TransformFloodfill() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("TransformFloodfill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.canvas.Id != tt.want.Id {
				t.Errorf("TransformRectangle() - Id\ngot  = %s\nwant = %s\n", tt.canvas.Id, tt.want.Id)
			}

			if tt.canvas.Name != tt.want.Name {
				t.Errorf("TransformRectangle() - Name\ngot  = %s\nwant = %s\n", tt.canvas.Name, tt.want.Name)
			}

			if tt.canvas.Content != tt.want.Content {
				t.Errorf("TransformRectangle() - Content\ngot  = %s\nwant = %s\n", tt.canvas.Content, tt.want.Content)
			}

			if tt.canvas.Width != tt.want.Width {
				t.Errorf("TransformRectangle() - Width\ngot  = %d\nwant = %d\n", tt.canvas.Width, tt.want.Width)
			}

			if tt.canvas.Height != tt.want.Height {
				t.Errorf("TransformRectangle() - Height\ngot  = %d\nwant = %d\n", tt.canvas.Height, tt.want.Height)
			}

			if !reflect.DeepEqual(tt.canvas, tt.want) {
				t.Errorf("TransformRectangle()\ngot:\n%s\nwant:\n%s", tt.canvas.String(), tt.want.String())
			}
		})
	}
}

func TestTransformFloodfill(t *testing.T) {
	tests := []struct {
		name    string
		canvas  *ascanvas.Canvas
		args    ascanvas.TransformFloodfillArgs
		want    *ascanvas.Canvas
		wantErr error
	}{
		{
			name: "Bad Coords 1",
			canvas: internal.CanvasFromText("1", "canvas 1", `
     
 XXXX
 X  X
 X  X
 XXXX`),
			args: ascanvas.TransformFloodfillArgs{
				Start: ascanvas.Coordinates{
					X: -1,
					Y: 2,
				},
				Fill: "O",
			},
			want: internal.CanvasFromText("1", "canvas 1", `
     
 XXXX
 X  X
 X  X
 XXXX`),
			wantErr: ascanvas.ErrOutOfBounds,
		},
		{
			name: "Bad Coords 2",
			canvas: internal.CanvasFromText("2", "canvas 2", `
     
 XXXX
 X  X
 X  X
 XXXX`),
			args: ascanvas.TransformFloodfillArgs{
				Start: ascanvas.Coordinates{
					X: 3,
					Y: -2,
				},
				Fill: "O",
			},
			want: internal.CanvasFromText("2", "canvas 2", `
     
 XXXX
 X  X
 X  X
 XXXX`),
			wantErr: ascanvas.ErrOutOfBounds,
		},
		{
			name: "Bad Coords 3",
			canvas: internal.CanvasFromText("3", "canvas 3", `
     
 XXXX
 X  X
 X  X
 XXXX`),
			args: ascanvas.TransformFloodfillArgs{
				Start: ascanvas.Coordinates{
					X: 300,
					Y: 2,
				},
				Fill: "O",
			},
			want: internal.CanvasFromText("3", "canvas 3", `
     
 XXXX
 X  X
 X  X
 XXXX`),
			wantErr: ascanvas.ErrOutOfBounds,
		},
		{
			name: "Bad Coords 4",
			canvas: internal.CanvasFromText("4", "canvas 4", `
     
 XXXX
 X  X
 X  X
 XXXX`),
			args: ascanvas.TransformFloodfillArgs{
				Start: ascanvas.Coordinates{
					X: 3,
					Y: 200,
				},
				Fill: "O",
			},
			want: internal.CanvasFromText("4", "canvas 4", `
     
 XXXX
 X  X
 X  X
 XXXX`),
			wantErr: ascanvas.ErrOutOfBounds,
		},
		{
			name: "OK fill 1",
			canvas: internal.CanvasFromText("1", "canvas 1", `
      
 XXXX 
 X  X 
 X  X 
 XXXX `),
			args: ascanvas.TransformFloodfillArgs{
				Start: ascanvas.Coordinates{
					X: 3,
					Y: 2,
				},
				Fill: "O",
			},
			want: internal.CanvasFromText("1", "canvas 1", `
      
 XXXX 
 XOOX 
 XOOX 
 XXXX `),
		},
		{
			name: "OK fill 2",
			canvas: internal.CanvasFromText("2", "canvas 2", `
     
 XXXX
 X  X
 X  X
 XXXX`),
			args: ascanvas.TransformFloodfillArgs{
				Start: ascanvas.Coordinates{
					X: 3,
					Y: 3,
				},
				Fill: "O",
			},
			want: internal.CanvasFromText("2", "canvas 2", `
     
 XXXX
 XOOX
 XOOX
 XXXX`),
		},
		{
			name: "OK fill 3",
			canvas: internal.CanvasFromText("3", "canvas 3", `
      
 XXXX 
 X  X 
 X  X 
 XXXX `),
			args: ascanvas.TransformFloodfillArgs{
				Start: ascanvas.Coordinates{
					X: 2,
					Y: 3,
				},
				Fill: "O",
			},
			want: internal.CanvasFromText("3", "canvas 3", `
      
 XXXX 
 XOOX 
 XOOX 
 XXXX `),
		},
		{
			name: "OK fill 4",
			canvas: internal.CanvasFromText("4", "canvas 4", `
     
 XXXX
 X  X
 X  X
 XXXX`),
			args: ascanvas.TransformFloodfillArgs{
				Start: ascanvas.Coordinates{
					X: 0,
					Y: 0,
				},
				Fill: "O",
			},
			want: internal.CanvasFromText("4", "canvas 4", `
OOOOO
OXXXX
OX  X
OX  X
OXXXX`),
		},
		{
			name: "OK fill 5",
			canvas: internal.CanvasFromText("5", "canvas 5", `
       
 XXXX  
 X  X  
 X  X  
 XXXX  `),
			args: ascanvas.TransformFloodfillArgs{
				Start: ascanvas.Coordinates{
					X: 4,
					Y: 4,
				},
				Fill: "O",
			},
			want: internal.CanvasFromText("5", "canvas 5", `
       
 OOOO  
 O  O  
 O  O  
 OOOO  `),
		},
		{
			name: "fixture 3.4",
			canvas: internal.CanvasFromText("1", "fixture 3.4", `
              .......       
              .......       
              .......       
OOOOOOOO      .......       
O      O      .......       
O    XXXXX    .......       
OOOOOXXXXX                  
     XXXXX                  
                            
                            
                            
                            `),
			args: ascanvas.TransformFloodfillArgs{
				Start: ascanvas.Coordinates{
					X: 0,
					Y: 0,
				},
				Fill: "-",
			},
			want: internal.CanvasFromText("1", "fixture 3.4", `
--------------.......-------
--------------.......-------
--------------.......-------
OOOOOOOO------.......-------
O      O------.......-------
O    XXXXX----.......-------
OOOOOXXXXX------------------
-----XXXXX------------------
----------------------------
----------------------------
----------------------------
----------------------------`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ascanvas.TransformFloodfill(tt.canvas, tt.args)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("TransformFloodfill() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("TransformFloodfill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(tt.canvas, tt.want) {
				t.Errorf("TransformFloodfill()\ngot:\n%s\nwant:\n%s", tt.canvas.String(), tt.want.String())
			}
		})
	}
}
