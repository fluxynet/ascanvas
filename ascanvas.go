package ascanvas

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

//go:generate swag init -g cmd/ascanvas/cmd_serve.go -o docs/ascanvas

//go:generate mockery --name=CanvasRepository --output=repo/mocks --filename=repo_mocks.go
//go:generate mockery --name=CanvasBroadcaster --output=broadcaster/mocks --filename=broadcaster_mocks.go

// ObserveALL is a special keyword to observe all events
const ObserveALL = "all"

//go:embed assets/index.html
var IndexHTML []byte

var (
	// ErrInvalidInput is for validation errors
	ErrInvalidInput = errors.New("invalid input")

	// ErrNotFound is when we cannot find something
	ErrNotFound = errors.New("item not found")

	// ErrOutOfBounds is when a cooridnate is out of bounds
	ErrOutOfBounds = errors.New("out of bounds")
)

// Canvas is an ascii art drawing
type Canvas struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
}

func (c Canvas) String() string {
	var p string

	for t := strings.ReplaceAll(c.Content, " ", "_"); t != ""; t = t[c.Width:] {
		p += t[:c.Width] + "\n"
	}

	return fmt.Sprintf(
		"id = %s, name = %s, w = %d, h = %d\n%s\n",
		c.Id,
		c.Name,
		c.Width,
		c.Height,
		p,
	)
}

func (c Canvas) AsGrid() [][]string {
	var (
		grid = make([][]string, c.Height)
		p    = 0
	)

	for y := 0; y < c.Height; y++ {
		grid[y] = make([]string, c.Width)
		for x := 0; x < c.Width; x++ {
			grid[y][x] = string(c.Content[p])
			p += 1
		}
	}

	return grid
}

func (c *Canvas) FromGrid(grid [][]string) {
	var b strings.Builder

	for x := range grid {
		for y := range grid[x] {
			b.WriteString(grid[x][y])
		}
	}

	c.Content = b.String()
}

// AsLogFields is a helper for logging
func (c Canvas) AsLogFields() []zap.Field {
	return []zap.Field{
		zap.String("Id", c.Id),
		zap.String("Name", c.Name),
		zap.String("Content", c.Content),
		zap.Int("Width", c.Width),
		zap.Int("Height", c.Height),
	}
}

func (c Canvas) Contains(coords Coordinates) bool {
	return coords.Y >= 0 && coords.Y < c.Height && coords.X >= 0 && coords.X < c.Width
}

// Coordinates on a cartesian plane
type Coordinates struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (c Coordinates) String() string {
	return strconv.Itoa(c.X) + "," + strconv.Itoa(c.Y)
}

// CanvasRepository is for canvas persistence
type CanvasRepository interface {
	Create(ctx context.Context, canvas Canvas) error
	Update(ctx context.Context, canvas Canvas) error
	Get(ctx context.Context, id string) (*Canvas, error)
	List(ctx context.Context) ([]Canvas, error)
	Delete(ctx context.Context, id string) error
}

// CanvasEventName is the type of event emitted by CanvasEvent
type CanvasEventName string

const (
	// CanvasEventCreated is emitted when a new Canvas has been created
	CanvasEventCreated CanvasEventName = "CREATED"

	// CanvasEventUpdated is emitted when an existing Canvas has been updated
	CanvasEventUpdated CanvasEventName = "UPDATED"

	// CanvasEventDeleted is emitted when an existing Canvas has been deleted
	CanvasEventDeleted CanvasEventName = "DELETED"
)

// CanvasEvent emitted by a CanvasBroadcaster
type CanvasEvent struct {
	Name   CanvasEventName
	Canvas Canvas
}

// CanvasBroadcaster allows changes to be published or observed on a specific canvas
type CanvasBroadcaster interface {
	Observe(ctx context.Context, id string) (StopObserveFunc, <-chan CanvasEvent, error)
	Broadcast(ctx context.Context, event CanvasEvent) error
}

// StopObserveFunc can be used to stop listening after calling Observe
type StopObserveFunc func()
