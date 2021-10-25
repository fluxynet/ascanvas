package canvas_test

import (
	"database/sql"
	"net/http"
	"strings"
	"testing"

	"go.uber.org/zap/zaptest"
	_ "modernc.org/sqlite"

	"github.com/fluxynet/ascanvas"
	"github.com/fluxynet/ascanvas/broadcaster/memory"
	"github.com/fluxynet/ascanvas/internal"
	"github.com/fluxynet/ascanvas/repo/sequel"
	"github.com/fluxynet/ascanvas/web"
	"github.com/fluxynet/ascanvas/web/canvas"
)

func makeDb() *sql.DB {
	var db, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		panic("failed to open database connection: " + err.Error())
	} else if err = db.Ping(); err != nil {
		panic("failed to ping database: " + err.Error())
	} else if _, err = db.Exec(sequel.SQLiteSchemaInit); err != nil {
		panic("failed to initialize schema: " + err.Error())
	}

	return db
}

func makeWebCanvas(t *testing.T, db *sql.DB) canvas.WebCanvas {
	return canvas.WebCanvas{
		GetID: web.StaticIDGetter("1", nil),
		Service: &ascanvas.CanvasService{
			Repo:        &sequel.Repository{DB: db},
			BroadCaster: memory.New(),
			Logger:      zaptest.NewLogger(t),
			GenerateID:  ascanvas.StaticUUIDGenerator("1", nil),
			Broadcast:   ascanvas.SyncBroadcast,
		},
	}
}

func canvasCreate(t *testing.T, db *sql.DB) http.HandlerFunc {
	return makeWebCanvas(t, db).Create
}

func canvasRectangle(t *testing.T, db *sql.DB) http.HandlerFunc {
	return makeWebCanvas(t, db).Rectangle
}

func canvasFloodfill(t *testing.T, db *sql.DB) http.HandlerFunc {
	return makeWebCanvas(t, db).Floodfill
}

func TestWebCanvas(t *testing.T) {
	str24x9 := strings.Repeat(" ", 24*9)
	str21x8 := strings.Repeat(" ", 21*8)
	str28x12 := strings.Repeat(" ", 28*12)

	headerJSON := map[string][]string{
		"Content-Type": []string{web.ContentTypeJSON},
	}

	type handlerMaker func(t *testing.T, db *sql.DB) http.HandlerFunc

	type test struct {
		req          internal.HttpTest
		handlerMaker handlerMaker
	}

	tests := []struct {
		name  string
		tests []test
	}{
		{
			name: "Test fixture 1",
			tests: []test{
				{
					handlerMaker: canvasCreate,
					req: internal.HttpTest{
						Request: internal.HttpTestRequest{
							Path:   "/",
							Method: http.MethodPost,
							Body:   `{"name": "F1","fill": "","width":24,"height":9}`,
						},
						Want: internal.HttpTestWant{
							Status: http.StatusOK,
							Header: headerJSON,
							Body:   `{"id":"1","name":"F1","content":"` + str24x9 + `","width":24,"height":9}`,
						},
					},
				},
				{
					handlerMaker: canvasRectangle,
					req: internal.HttpTest{
						Request: internal.HttpTestRequest{
							Path:   "/",
							Method: http.MethodPatch,
							Body:   `{"top_left":{"x":3,"y":2},"width":5,"height":3,"fill":"X","outline":"@"}`,
						},
						Want: internal.HttpTestWant{
							Status: http.StatusOK,
							Header: headerJSON,
							Body: internal.CanvasJsonFromText(
								"1",
								"F1",
								`
                        
                        
   @@@@@                
   @XXX@                
   @@@@@                
                        
                        
                        
                        `,
							),
						},
					},
				},
				{
					handlerMaker: canvasRectangle,
					req: internal.HttpTest{
						Request: internal.HttpTestRequest{
							Path:   "/",
							Method: http.MethodPatch,
							Body:   `{"top_left":{"x":10,"y":3},"width":14,"height":6,"fill":"O","outline":"X"}`,
						},
						Want: internal.HttpTestWant{
							Status: http.StatusOK,
							Header: headerJSON,
							Body: internal.CanvasJsonFromText(
								"1",
								"F1",
								`
                        
                        
   @@@@@                
   @XXX@  XXXXXXXXXXXXXX
   @@@@@  XOOOOOOOOOOOOX
          XOOOOOOOOOOOOX
          XOOOOOOOOOOOOX
          XOOOOOOOOOOOOX
          XXXXXXXXXXXXXX`,
							),
						},
					},
				},
			},
		},
		{
			name: "Test fixture 2",
			tests: []test{
				{
					handlerMaker: canvasCreate,
					req: internal.HttpTest{
						Request: internal.HttpTestRequest{
							Path:   "/",
							Method: http.MethodPatch,
							Body:   `{"name": "F2","fill": "","width":21,"height":8}`,
						},
						Want: internal.HttpTestWant{
							Status: http.StatusOK,
							Header: headerJSON,
							Body:   `{"id":"1","name":"F2","content":"` + str21x8 + `","width":21,"height":8}`,
						},
					},
				},
				{
					handlerMaker: canvasRectangle,
					req: internal.HttpTest{
						Request: internal.HttpTestRequest{
							Path:   "/",
							Method: http.MethodPatch,
							Body:   `{"top_left":{"x":15,"y":0},"width":7,"height":6,"fill":".","outline":""}`,
						},
						Want: internal.HttpTestWant{
							Status: http.StatusOK,
							Header: headerJSON,
							Body: internal.CanvasJsonFromText(
								"1",
								"F2",
								`
               ......
               ......
               ......
               ......
               ......
               ......
                     
                     `,
							),
						},
					},
				},
				{
					handlerMaker: canvasRectangle,
					req: internal.HttpTest{
						Request: internal.HttpTestRequest{
							Path:   "/",
							Method: http.MethodPatch,
							Body:   `{"top_left":{"x":0,"y":3},"width":8,"height":4,"fill":"","outline":"O"}`,
						},
						Want: internal.HttpTestWant{
							Status: http.StatusOK,
							Header: headerJSON,
							Body: internal.CanvasJsonFromText(
								"1",
								"F2",
								`
               ......
               ......
               ......
OOOOOOOO       ......
O      O       ......
O      O       ......
OOOOOOOO             
                     `,
							),
						},
					},
				},
				{
					handlerMaker: canvasRectangle,
					req: internal.HttpTest{
						Request: internal.HttpTestRequest{
							Path:   "/",
							Method: http.MethodPatch,
							Body:   `{"top_left":{"x":5,"y":5},"width":5,"height":3,"fill":"X","outline":"X"}`,
						},
						Want: internal.HttpTestWant{
							Status: http.StatusOK,
							Header: headerJSON,
							Body: internal.CanvasJsonFromText(
								"1",
								"F2",
								`
               ......
               ......
               ......
OOOOOOOO       ......
O      O       ......
O    XXXXX     ......
OOOOOXXXXX           
     XXXXX           `,
							),
						},
					},
				},
			},
		},
		{
			name: "Test fixture 3",
			tests: []test{
				{
					handlerMaker: canvasCreate,
					req: internal.HttpTest{
						Request: internal.HttpTestRequest{
							Path:   "/",
							Method: http.MethodPatch,
							Body:   `{"name": "F3","fill": "","width":28,"height":12}`,
						},
						Want: internal.HttpTestWant{
							Status: http.StatusOK,
							Header: headerJSON,
							Body:   `{"id":"1","name":"F3","content":"` + str28x12 + `","width":28,"height":12}`,
						},
					},
				},
				{
					handlerMaker: canvasRectangle,
					req: internal.HttpTest{
						Request: internal.HttpTestRequest{
							Path:   "/",
							Method: http.MethodPatch,
							Body:   `{"top_left":{"x":15,"y":0},"width":7,"height":6,"fill":".","outline":""}`,
						},
						Want: internal.HttpTestWant{
							Status: http.StatusOK,
							Header: headerJSON,
							Body: internal.CanvasJsonFromText(
								"1",
								"F3",
								`
               .......      
               .......      
               .......      
               .......      
               .......      
               .......      
                            
                            
                            
                            
                            
                            `,
							),
						},
					},
				},
				{
					handlerMaker: canvasRectangle,
					req: internal.HttpTest{
						Request: internal.HttpTestRequest{
							Path:   "/",
							Method: http.MethodPatch,
							Body:   `{"top_left":{"x":0,"y":3},"width":8,"height":4,"fill":"","outline":"O"}`,
						},
						Want: internal.HttpTestWant{
							Status: http.StatusOK,
							Header: headerJSON,
							Body: internal.CanvasJsonFromText(
								"1",
								"F3",
								`
               .......      
               .......      
               .......      
OOOOOOOO       .......      
O      O       .......      
O      O       .......      
OOOOOOOO                    
                            
                            
                            
                            
                            `,
							),
						},
					},
				},
				{
					handlerMaker: canvasRectangle,
					req: internal.HttpTest{
						Request: internal.HttpTestRequest{
							Path:   "/",
							Method: http.MethodPatch,
							Body:   `{"top_left":{"x":5,"y":5},"width":5,"height":3,"fill":"X","outline":"X"}`,
						},
						Want: internal.HttpTestWant{
							Status: http.StatusOK,
							Header: headerJSON,
							Body: internal.CanvasJsonFromText(
								"1",
								"F3",
								`
               .......      
               .......      
               .......      
OOOOOOOO       .......      
O      O       .......      
O    XXXXX     .......      
OOOOOXXXXX                  
     XXXXX                  
                            
                            
                            
                            `,
							),
						},
					},
				},
				{
					handlerMaker: canvasFloodfill,
					req: internal.HttpTest{
						Request: internal.HttpTestRequest{
							Path:   "/",
							Method: http.MethodPatch,
							Body:   `{"start":{"x":0,"y":0},"fill":"-"}`,
						},
						Want: internal.HttpTestWant{
							Status: http.StatusOK,
							Header: headerJSON,
							Body: internal.CanvasJsonFromText(
								"1",
								"F3",
								`
---------------.......------
---------------.......------
---------------.......------
OOOOOOOO-------.......------
O      O-------.......------
O    XXXXX-----.......------
OOOOOXXXXX------------------
-----XXXXX------------------
----------------------------
----------------------------
----------------------------
----------------------------`,
							),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := makeDb()

			defer internal.Closed(db)

			for i := range tt.tests {
				if !tt.tests[i].req.Assert(t, tt.tests[i].handlerMaker(t, db)) {
					t.Errorf("%s [%d] failed", tt.name, i)
					return
				}
			}
		})
	}
}
