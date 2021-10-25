package sequel_test

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/fluxynet/ascanvas"
	"github.com/fluxynet/ascanvas/internal"
	"github.com/fluxynet/ascanvas/repo/sequel"
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

func TestRepository(t *testing.T) {
	var (
		db = makeDb()
		r  = RepositoryTest{DB: db}
	)

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	// order of the tests is important for value checking; must be: Create, Delete, Get, Update, List
	t.Run("repository_Create", r.Create)
	t.Run("repository_Delete", r.Delete)
	t.Run("repository_Get", r.Get)
	t.Run("repository_Update", r.Update)
	t.Run("repository_List", r.List)
}

type RepositoryTest struct {
	DB *sql.DB
}

func (r *RepositoryTest) Create(t *testing.T) {
	tests := []struct {
		name    string
		canvas  ascanvas.Canvas
		wantErr bool
	}{
		{
			name: "normal insert 1",
			canvas: ascanvas.Canvas{
				Id:      "1",
				Name:    "Canvas 1",
				Content: strings.Repeat(".", 255),
				Width:   10,
				Height:  20,
			},
			wantErr: false,
		},
		{
			name: "normal insert 2",
			canvas: ascanvas.Canvas{
				Id:      "2",
				Name:    "Canvas 2",
				Content: strings.Repeat(".", 15),
				Width:   10,
				Height:  20,
			},
			wantErr: false,
		},
		{
			name: "normal insert 3",
			canvas: ascanvas.Canvas{
				Id:      "3",
				Name:    "Canvas 3",
				Content: strings.Repeat(".", 65),
				Width:   10,
				Height:  20,
			},
			wantErr: false,
		},
		{
			name: "duplicate id insert",
			canvas: ascanvas.Canvas{
				Id:      "1",
				Name:    "Canvas 1B",
				Content: strings.Repeat(".", 10),
				Width:   1,
				Height:  2,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := sequel.Repository{
				DB: r.DB,
			}

			if err := repo.Create(context.Background(), tt.canvas); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %s, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func (r *RepositoryTest) Delete(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{
			name: "non-existent",
			id:   "100",
		},
		{
			name: "id = 3",
			id:   "3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := sequel.Repository{
				DB: r.DB,
			}

			if err := repo.Delete(context.Background(), tt.id); err != nil {
				t.Errorf("Delete() error = %v", err)
			}
		})
	}
}

func (r *RepositoryTest) Get(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		want    *ascanvas.Canvas
		wantErr error
	}{
		{
			name:    "non-existent",
			id:      "50",
			want:    nil,
			wantErr: ascanvas.ErrNotFound,
		},
		{
			name: "id = 1",
			id:   "1",
			want: &ascanvas.Canvas{
				Id:      "1",
				Name:    "Canvas 1",
				Content: strings.Repeat(".", 255),
				Width:   10,
				Height:  20,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := sequel.Repository{
				DB: r.DB,
			}

			got, err := repo.Get(context.Background(), tt.id)
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

func (r *RepositoryTest) Update(t *testing.T) {
	tests := []struct {
		name   string
		canvas ascanvas.Canvas
	}{
		{
			name: "non-existent",
			canvas: ascanvas.Canvas{
				Id:      "404",
				Name:    "Canvas 404",
				Content: strings.Repeat(".", 404),
				Width:   404,
				Height:  404,
			},
		},
		{
			name: "id = 2",
			canvas: ascanvas.Canvas{
				Id:      "2",
				Name:    "Canvas two",
				Content: strings.Repeat(".", 800*600),
				Width:   800,
				Height:  600,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := sequel.Repository{
				DB: r.DB,
			}

			if err := repo.Update(context.Background(), tt.canvas); err != nil {
				t.Errorf("Update() error = %v", err)
			}
		})
	}
}

func (r *RepositoryTest) List(t *testing.T) {
	var empty = makeDb()
	defer internal.Closed(empty)

	tests := []struct {
		name string
		db   *sql.DB
		want []ascanvas.Canvas
	}{
		{
			name: "empty db",
			db:   empty,
			want: []ascanvas.Canvas{},
		},
		{
			name: "non-empty db",
			db:   r.DB,
			want: []ascanvas.Canvas{
				{
					Id:      "1",
					Name:    "Canvas 1",
					Content: strings.Repeat(".", 255),
					Width:   10,
					Height:  20,
				},
				{
					Id:      "2",
					Name:    "Canvas two",
					Content: strings.Repeat(".", 800*600),
					Width:   800,
					Height:  600,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := sequel.Repository{
				DB: tt.db,
			}

			got, err := repo.List(context.Background())
			if err != nil {
				t.Errorf("List() error = %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}
