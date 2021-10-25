package ascanvas_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"go.uber.org/zap/zaptest"

	"github.com/fluxynet/ascanvas"
	mb "github.com/fluxynet/ascanvas/broadcaster/mocks"
	mr "github.com/fluxynet/ascanvas/repo/mocks"
)

func TestCanvasService_Create(t *testing.T) {
	tests := []struct {
		name     string
		mustCall bool
		args     ascanvas.CreateArgs
		want     *ascanvas.Canvas
		wantErr  bool
	}{
		{
			name:     "name empty",
			mustCall: false,
			args: ascanvas.CreateArgs{
				Name:   "",
				Fill:   ".",
				Width:  10,
				Height: 20,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:     "fill empty",
			mustCall: true,
			args: ascanvas.CreateArgs{
				Name:   "Foo",
				Fill:   "",
				Width:  2,
				Height: 1,
			},
			want: &ascanvas.Canvas{
				Id:      "1",
				Name:    "Foo",
				Content: "  ",
				Width:   2,
				Height:  1,
			},
		},
		{
			name:     "width negative",
			mustCall: false,
			args: ascanvas.CreateArgs{
				Name:   "Foo",
				Fill:   ".",
				Width:  -1,
				Height: 30,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:     "height negative",
			mustCall: false,
			args: ascanvas.CreateArgs{
				Name:   "Foo",
				Fill:   ".",
				Width:  40,
				Height: -1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:     "normal create",
			mustCall: true,
			args: ascanvas.CreateArgs{
				Name:   "Test",
				Fill:   ".",
				Width:  10,
				Height: 15,
			},
			want: &ascanvas.Canvas{
				Id:      "1",
				Name:    "Test",
				Content: strings.Repeat(".", 150),
				Width:   10,
				Height:  15,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			repo := &mr.CanvasRepository{}
			brd := &mb.CanvasBroadcaster{}
			s := ascanvas.CanvasService{
				Repo:        repo,
				BroadCaster: brd,
				Logger:      zaptest.NewLogger(t),
				GenerateID:  ascanvas.StaticUUIDGenerator("1", nil),
				Broadcast:   ascanvas.SyncBroadcast,
			}

			var event ascanvas.CanvasEvent

			if tt.mustCall {
				event = ascanvas.CanvasEvent{
					Name:   ascanvas.CanvasEventCreated,
					Canvas: *tt.want,
				}

				repo.On("Create", ctx, *tt.want).Return(nil)
				brd.On("Broadcast", ctx, event).Return(nil)
			}

			got, err := s.Create(ctx, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			if tt.mustCall && !repo.AssertNumberOfCalls(t, "Create", 1) &&
				!repo.AssertCalled(t, "Create", *tt.want) {
				t.Errorf("Call Error: Create")
				return
			} else if !tt.mustCall && !repo.AssertNotCalled(t, "Create") {
				t.Errorf("Call Error: Create")
				return
			}

			if tt.mustCall && !brd.AssertNumberOfCalls(t, "Broadcast", 1) ||
				!brd.AssertCalled(t, "Broadcast", ctx, event) {
				t.Errorf("Call Error: Broadcast")
				return
			} else if !tt.mustCall && !repo.AssertNotCalled(t, "Broadcast") {
				t.Errorf("Call Error: Broadcast")
				return
			}

			if got.Id == "" {
				t.Errorf("Create() Id is empty")
				return
			}

			if got.Name != tt.want.Name {
				t.Errorf("Create() Name = %v, Want = %v", got.Name, tt.want.Name)
				return
			}

			if got.Content != tt.want.Content {
				t.Errorf("Create() Content = %v, Want = %v", got.Content, tt.want.Content)
				return
			}

			if got.Width != tt.want.Width {
				t.Errorf("Create() Width = %v, Want = %v", got.Width, tt.want.Width)
				return
			}

			if got.Height != tt.want.Height {
				t.Errorf("Create() Height = %v, Want = %v", got.Height, tt.want.Height)
				return
			}

		})
	}
}

func TestCanvasService_Delete(t *testing.T) {
	type repoGet struct {
		count        int
		id           string
		returnCanvas *ascanvas.Canvas
		returnErr    error
	}

	type repoDelete struct {
		count     int
		id        string
		returnErr error
	}

	tests := []struct {
		name       string
		id         string
		repoGet    repoGet
		repoDelete repoDelete
		wantErr    error
	}{
		{
			name: "delete not found",
			id:   "2",
			repoGet: repoGet{
				count:        1,
				id:           "2",
				returnCanvas: nil,
				returnErr:    ascanvas.ErrNotFound,
			},
			repoDelete: repoDelete{
				count: 0,
			},
			wantErr: ascanvas.ErrNotFound,
		},
		{
			name: "delete ok",
			id:   "2",
			repoGet: repoGet{
				id:    "2",
				count: 1,
				returnCanvas: &ascanvas.Canvas{
					Id:      "2",
					Name:    "Foo",
					Content: "ABCDEF",
					Width:   2,
					Height:  3,
				},
				returnErr: nil,
			},
			repoDelete: repoDelete{
				count:     1,
				id:        "2",
				returnErr: nil,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			repo := &mr.CanvasRepository{}
			brd := &mb.CanvasBroadcaster{}
			s := ascanvas.CanvasService{
				Repo:        repo,
				BroadCaster: brd,
				Logger:      zaptest.NewLogger(t),
				Broadcast:   ascanvas.SyncBroadcast,
			}

			var event ascanvas.CanvasEvent

			if tt.repoGet.returnCanvas != nil {
				event = ascanvas.CanvasEvent{
					Name:   ascanvas.CanvasEventDeleted,
					Canvas: *tt.repoGet.returnCanvas,
				}
			}

			repo.On("Get", ctx, tt.repoGet.id).Return(tt.repoGet.returnCanvas, tt.repoGet.returnErr)
			repo.On("Delete", ctx, tt.repoDelete.id).Return(tt.repoDelete.returnErr)
			brd.On("Broadcast", ctx, event).Return(nil)

			if err := s.Delete(ctx, tt.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("Delete() error = %v", err)
				return
			}

			if tt.repoGet.count == 0 && !repo.AssertNotCalled(t, "Get") {
				return
			} else if !repo.AssertCalled(t, "Get", ctx, tt.id) ||
				repo.AssertNumberOfCalls(t, "Get", tt.repoGet.count) {
				return
			}

			if tt.repoDelete.count == 0 && !repo.AssertNotCalled(t, "Delete") {
				return
			} else if !repo.AssertCalled(t, "Delete", ctx, tt.id) ||
				repo.AssertNumberOfCalls(t, "Delete", tt.repoDelete.count) {
				return
			}

			if !brd.AssertNumberOfCalls(t, "Broadcast", 1) ||
				!brd.AssertCalled(t, "Broadcast", event) {
				return
			}
		})
	}
}

func TestCanvasService_Get(t *testing.T) {
	type RepoGet struct {
		Id           string
		ReturnCanvas *ascanvas.Canvas
		ReturnErr    error
	}

	tests := []struct {
		name    string
		repoGet RepoGet
		id      string
		want    *ascanvas.Canvas
		wantErr error
	}{
		{
			name: "existing id",
			id:   "2",
			repoGet: RepoGet{
				Id: "2",
				ReturnCanvas: &ascanvas.Canvas{
					Id:      "2",
					Name:    "Foo",
					Content: "XXYY",
					Width:   2,
					Height:  2,
				},
				ReturnErr: nil,
			},
			want: &ascanvas.Canvas{
				Id:      "2",
				Name:    "Foo",
				Content: "XXYY",
				Width:   2,
				Height:  2,
			},
			wantErr: nil,
		},
		{
			name: "non-existent",
			id:   "1",
			repoGet: RepoGet{
				Id:           "1",
				ReturnCanvas: nil,
				ReturnErr:    ascanvas.ErrNotFound,
			},
			want:    nil,
			wantErr: ascanvas.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			repo := &mr.CanvasRepository{}
			repo.On("Get", ctx, tt.repoGet.Id).Return(tt.repoGet.ReturnCanvas, tt.repoGet.ReturnErr)

			brd := &mb.CanvasBroadcaster{}

			s := ascanvas.CanvasService{
				Repo:        repo,
				BroadCaster: brd,
				Logger:      zaptest.NewLogger(t),
			}

			if !brd.AssertNotCalled(t, "Broadcast") {
				return
			}

			got, err := s.Get(ctx, tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCanvasService_List(t *testing.T) {
	err := errors.New("foo")

	tests := []struct {
		name    string
		want    []ascanvas.Canvas
		wantErr error
	}{
		{
			name:    "error",
			want:    []ascanvas.Canvas{},
			wantErr: err,
		},
		{
			name:    "empty",
			want:    nil,
			wantErr: nil,
		},
		{
			name: "non-empty",
			want: []ascanvas.Canvas{
				{
					Id:      "1",
					Name:    "Foo",
					Content: "AB",
					Width:   2,
					Height:  1,
				},
				{
					Id:      "100",
					Name:    "Bar",
					Content: "ABCD",
					Width:   1,
					Height:  4,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			brd := &mb.CanvasBroadcaster{}

			repo := &mr.CanvasRepository{}
			repo.On("List", ctx).Return(tt.want, tt.wantErr)

			s := ascanvas.CanvasService{
				Repo:        repo,
				BroadCaster: brd,
				Logger:      zaptest.NewLogger(t),
			}

			got, err := s.List(ctx)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("List() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}

			if !repo.AssertNumberOfCalls(t, "List", 1) {
				return
			}

			if !brd.AssertNotCalled(t, "Observe") {
				return
			}
		})
	}
}

func TestCanvasService_Observe(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		callID  string
		wantErr error
	}{
		{
			name:    "specific id",
			id:      "5",
			callID:  "5",
			wantErr: nil,
		},
		{
			name:    "all",
			id:      ascanvas.ObserveALL,
			callID:  ascanvas.ObserveALL,
			wantErr: nil,
		},
		{
			name:    "empty",
			id:      "",
			callID:  ascanvas.ObserveALL,
			wantErr: nil,
		},
		{
			name:    "error",
			id:      "123",
			callID:  "123",
			wantErr: errors.New("foo"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			repo := &mr.CanvasRepository{}

			ch := make(<-chan ascanvas.CanvasEvent)
			brd := &mb.CanvasBroadcaster{}
			brd.On("Observe", ctx, tt.callID).Return(nil, ch, tt.wantErr)

			s := ascanvas.CanvasService{
				Repo:        repo,
				BroadCaster: brd,
				Logger:      zaptest.NewLogger(t),
			}

			_, got, err := s.Observe(ctx, tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Observe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, ch) {
				t.Errorf("Observe() got = %v, want %v", got, ch)
			}

			if !brd.AssertNumberOfCalls(t, "Observe", 1) ||
				!brd.AssertCalled(t, "Observe", ctx, tt.callID) {
				return
			}
		})
	}
}

func TestCanvasService_ApplyFloodfill(t *testing.T) {
	errFoo := errors.New("foo")

	type repoGet struct {
		ReturnCanvas *ascanvas.Canvas
		ReturnErr    error
	}

	type repoUpdate struct {
		Canvas    ascanvas.Canvas
		ReturnErr error
	}

	tests := []struct {
		name       string
		id         string
		args       ascanvas.TransformFloodfillArgs
		repoGet    repoGet
		repoUpdate repoUpdate
		want       *ascanvas.Canvas
		wantErr    error
	}{
		{
			name: "not found",
			id:   "1",
			args: ascanvas.TransformFloodfillArgs{
				Start: ascanvas.Coordinates{
					X: 0,
					Y: 0,
				},
				Fill: "x",
			},
			repoGet: repoGet{
				ReturnCanvas: nil,
				ReturnErr:    ascanvas.ErrNotFound,
			},
			want:    nil,
			wantErr: ascanvas.ErrNotFound,
		},
		{
			name: "updated err",
			id:   "1",
			args: ascanvas.TransformFloodfillArgs{
				Start: ascanvas.Coordinates{
					X: 0,
					Y: 0,
				},
				Fill: "x",
			},
			repoGet: repoGet{
				ReturnCanvas: &ascanvas.Canvas{
					Id:      "1",
					Name:    "Foo",
					Content: "....",
					Width:   2,
					Height:  2,
				},
				ReturnErr: nil,
			},
			repoUpdate: repoUpdate{
				Canvas: ascanvas.Canvas{
					Id:      "1",
					Name:    "Foo",
					Content: "xxxx",
					Width:   2,
					Height:  2,
				},
				ReturnErr: errFoo,
			},
			want:    nil,
			wantErr: errFoo,
		},
		{
			name: "updated ok",
			id:   "1",
			args: ascanvas.TransformFloodfillArgs{
				Start: ascanvas.Coordinates{
					X: 0,
					Y: 0,
				},
				Fill: "x",
			},
			repoGet: repoGet{
				ReturnCanvas: &ascanvas.Canvas{
					Id:      "1",
					Name:    "Foo",
					Content: "....",
					Width:   2,
					Height:  2,
				},
				ReturnErr: nil,
			},
			repoUpdate: repoUpdate{
				Canvas: ascanvas.Canvas{
					Id:      "1",
					Name:    "Foo",
					Content: "xxxx",
					Width:   2,
					Height:  2,
				},
				ReturnErr: nil,
			},
			want: &ascanvas.Canvas{
				Id:      "1",
				Name:    "Foo",
				Content: "xxxx",
				Width:   2,
				Height:  2,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			repo := &mr.CanvasRepository{}
			brd := &mb.CanvasBroadcaster{}

			s := &ascanvas.CanvasService{
				Repo:        repo,
				BroadCaster: brd,
				Logger:      zaptest.NewLogger(t),
				Broadcast:   ascanvas.SyncBroadcast,
			}

			var event ascanvas.CanvasEvent

			if tt.want != nil {
				event = ascanvas.CanvasEvent{
					Name:   ascanvas.CanvasEventUpdated,
					Canvas: *tt.want,
				}
			}

			repo.On("Get", ctx, tt.id).Return(tt.repoGet.ReturnCanvas, tt.repoGet.ReturnErr)
			repo.On("Update", ctx, tt.repoUpdate.Canvas).Return(tt.repoUpdate.ReturnErr)

			brd.On("Broadcast", ctx, event).Return(nil)

			got, err := s.ApplyFloodfill(ctx, tt.id, tt.args)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ApplyFloodfill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApplyFloodfill() got = %v, want %v", got, tt.want)
			}

			if tt.wantErr != nil && !brd.AssertNotCalled(t, "Broadcast") {
				t.Errorf("Call Error = Broadcast")
				return
			} else if tt.wantErr == nil && !brd.AssertCalled(t, "Broadcast", ctx, event) && !brd.AssertNumberOfCalls(t, "Broadcast", 1) {
				t.Errorf("Call Error = Broadcast")
				return
			}
		})
	}
}

func TestCanvasService_ApplyRectangle(t *testing.T) {
	errFoo := errors.New("foo")

	type repoGet struct {
		ReturnCanvas *ascanvas.Canvas
		ReturnErr    error
	}

	type repoUpdate struct {
		Canvas    ascanvas.Canvas
		ReturnErr error
	}

	tests := []struct {
		name       string
		id         string
		args       ascanvas.TransformRectangleArgs
		repoGet    repoGet
		repoUpdate repoUpdate
		want       *ascanvas.Canvas
		wantErr    error
	}{
		{
			name: "not found",
			id:   "1",
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 0,
					Y: 0,
				},
				Width:  2,
				Height: 2,
				Fill:   "x",
			},
			repoGet: repoGet{
				ReturnCanvas: nil,
				ReturnErr:    ascanvas.ErrNotFound,
			},
			want:    nil,
			wantErr: ascanvas.ErrNotFound,
		},
		{
			name: "updated err",
			id:   "1",
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 0,
					Y: 0,
				},
				Width:  2,
				Height: 2,
				Fill:   "x",
			},
			repoGet: repoGet{
				ReturnCanvas: &ascanvas.Canvas{
					Id:      "1",
					Name:    "Foo",
					Content: "....",
					Width:   2,
					Height:  2,
				},
				ReturnErr: nil,
			},
			repoUpdate: repoUpdate{
				Canvas: ascanvas.Canvas{
					Id:      "1",
					Name:    "Foo",
					Content: "xxxx",
					Width:   2,
					Height:  2,
				},
				ReturnErr: errFoo,
			},
			want:    nil,
			wantErr: errFoo,
		},
		{
			name: "updated ok",
			id:   "1",
			args: ascanvas.TransformRectangleArgs{
				TopLeft: ascanvas.Coordinates{
					X: 0,
					Y: 0,
				},
				Width:  2,
				Height: 2,
				Fill:   "x",
			},
			repoGet: repoGet{
				ReturnCanvas: &ascanvas.Canvas{
					Id:      "1",
					Name:    "Foo",
					Content: "....",
					Width:   2,
					Height:  2,
				},
				ReturnErr: nil,
			},
			repoUpdate: repoUpdate{
				Canvas: ascanvas.Canvas{
					Id:      "1",
					Name:    "Foo",
					Content: "xxxx",
					Width:   2,
					Height:  2,
				},
				ReturnErr: nil,
			},
			want: &ascanvas.Canvas{
				Id:      "1",
				Name:    "Foo",
				Content: "xxxx",
				Width:   2,
				Height:  2,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			repo := &mr.CanvasRepository{}
			brd := &mb.CanvasBroadcaster{}

			s := &ascanvas.CanvasService{
				Repo:        repo,
				BroadCaster: brd,
				Logger:      zaptest.NewLogger(t),
				Broadcast:   ascanvas.SyncBroadcast,
			}

			var event ascanvas.CanvasEvent

			if tt.want != nil {
				event = ascanvas.CanvasEvent{
					Name:   ascanvas.CanvasEventUpdated,
					Canvas: *tt.want,
				}
			}

			repo.On("Get", ctx, tt.id).Return(tt.repoGet.ReturnCanvas, tt.repoGet.ReturnErr)
			repo.On("Update", ctx, tt.repoUpdate.Canvas).Return(tt.repoUpdate.ReturnErr)

			brd.On("Broadcast", ctx, event).Return(nil)

			got, err := s.ApplyRectangle(ctx, tt.id, tt.args)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ApplyRectangle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApplyRectangle() got = %v, want %v", got, tt.want)
			}

			if tt.wantErr != nil && !brd.AssertNotCalled(t, "Broadcast") {
				t.Errorf("Call Error = Broadcast")
				return
			} else if tt.wantErr == nil && !brd.AssertCalled(t, "Broadcast", ctx, event) && !brd.AssertNumberOfCalls(t, "Broadcast", 1) {
				t.Errorf("Call Error = Broadcast")
				return
			}
		})
	}
}
