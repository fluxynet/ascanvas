package sequel

import (
	"context"
	"database/sql"
	_ "embed"

	"github.com/fluxynet/ascanvas"
	"github.com/fluxynet/ascanvas/internal"
)

// SQLiteSchemaInit initializes the schema for sqlite
//go:embed sqlite/00-init.sql
var SQLiteSchemaInit string

type Repository struct {
	DB *sql.DB
}

func (r Repository) Create(ctx context.Context, canvas ascanvas.Canvas) error {
	var _, err = r.DB.ExecContext(
		ctx,
		`INSERT INTO "canvas" ("id", "name", "content", "width", "height") VALUES (?,?,?,?,?)`,
		canvas.Id,
		canvas.Name,
		canvas.Content,
		canvas.Width,
		canvas.Height,
	)

	return err
}

func (r Repository) Update(ctx context.Context, canvas ascanvas.Canvas) error {
	var _, err = r.DB.ExecContext(
		ctx,
		`UPDATE "canvas" SET "name" = ?, "content" = ?, "width" = ?, "height" = ? WHERE "id" = ?`,
		canvas.Name,
		canvas.Content,
		canvas.Width,
		canvas.Height,
		canvas.Id,
	)

	return err
}

func (r Repository) Get(ctx context.Context, id string) (*ascanvas.Canvas, error) {
	var (
		canvas ascanvas.Canvas

		rows, err = r.DB.QueryContext(
			ctx,
			`SELECT "id", "name", "content", "width", "height" FROM "canvas" WHERE "id" = ?`,
			id,
		)
	)

	if err != nil {
		return nil, err
	}

	defer internal.Closed(rows)

	if rows.Next() {
		err = rows.Scan(
			&canvas.Id,
			&canvas.Name,
			&canvas.Content,
			&canvas.Width,
			&canvas.Height,
		)
	} else {
		return nil, ascanvas.ErrNotFound
	}

	return &canvas, err
}

func (r Repository) List(ctx context.Context) ([]ascanvas.Canvas, error) {
	var (
		canvases []ascanvas.Canvas

		rows, err = r.DB.QueryContext(
			ctx,
			`SELECT "id", "name", "content", "width", "height" FROM "canvas"`,
		)
	)

	if err != nil {
		return nil, err
	}

	defer internal.Closed(rows)

	for rows.Next() {
		var canvas ascanvas.Canvas

		err = rows.Scan(
			&canvas.Id,
			&canvas.Name,
			&canvas.Content,
			&canvas.Width,
			&canvas.Height,
		)

		if err != nil {
			return nil, err
		}

		canvases = append(canvases, canvas)
	}

	if canvases == nil {
		return []ascanvas.Canvas{}, nil
	}

	return canvases, err
}

func (r Repository) Delete(ctx context.Context, id string) error {
	var _, err = r.DB.ExecContext(ctx, `DELETE FROM "canvas" WHERE "id" = ?`, id)
	return err
}
