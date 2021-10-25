package ascanvas

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UUIDGeneratorFunc denotes functions that can provide uuids
type UUIDGeneratorFunc func() (string, error)

// UUIDGenerator is a uuid v4 generator
func UUIDGenerator() (string, error) {
	id, err := uuid.NewRandom()

	if err == nil {
		return id.String(), nil
	}

	return "", err
}

// StaticUUIDGenerator generates predefined uuid and err; useful for testing
func StaticUUIDGenerator(id string, err error) UUIDGeneratorFunc {
	return func() (string, error) {
		return id, err
	}
}

// BroadcastFunc denotes functions used for wrapping broadcasting features
type BroadcastFunc func(ctx context.Context, b CanvasBroadcaster, l *zap.Logger, event CanvasEvent)

// AsyncBroadcast uses a go routine
func AsyncBroadcast(ctx context.Context, b CanvasBroadcaster, l *zap.Logger, event CanvasEvent) {
	go SyncBroadcast(ctx, b, l, event)
}

// SyncBroadcast is blocking - useful for testing
func SyncBroadcast(ctx context.Context, b CanvasBroadcaster, l *zap.Logger, event CanvasEvent) {
	err := b.Broadcast(ctx, event)

	if err != nil {
		f := event.Canvas.AsLogFields()
		l.Error("Create::Broadcast::Failed", append(f, zap.Error(err))...)
	}
}

// CanvasService provides applicative features
type CanvasService struct {
	Repo        CanvasRepository
	BroadCaster CanvasBroadcaster
	Logger      *zap.Logger
	GenerateID  UUIDGeneratorFunc
	Broadcast   BroadcastFunc
}

type CreateArgs struct {
	Name   string `json:"name"`
	Fill   string `json:"fill"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// AsLogFields is a helper for logging
func (a CreateArgs) AsLogFields() []zap.Field {
	return []zap.Field{
		zap.String("Name", a.Name),
		zap.String("Fill", a.Fill),
		zap.Int("Width", a.Width),
		zap.Int("Height", a.Height),
	}
}

func (a CreateArgs) Validate() error {
	var errs []string

	if a.Name == "" {
		errs = append(errs, "name cannot be empty")
	}

	if len(a.Fill) != 1 {
		errs = append(errs, "fill must be exactly one character")
	}

	if a.Width < 1 {
		errs = append(errs, "width cannot be less than 1")
	}

	if a.Height < 1 {
		errs = append(errs, "height cannot be less than 1")
	}

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("%w: %s", ErrInvalidInput, strings.Join(errs, ", "))
}

// Create a new Canvas
func (s CanvasService) Create(ctx context.Context, args CreateArgs) (*Canvas, error) {
	var (
		id  string
		err error
	)

	if args.Fill == "" {
		args.Fill = " "
	}

	s.Logger.Debug("Create::Validating", args.AsLogFields()...)

	if err := args.Validate(); err != nil {
		s.Logger.Debug("Create::Validate::Failed", zap.Error(err))
		return nil, err
	}

	id, err = s.GenerateID()
	if err != nil {
		return nil, err
	}

	var canvas = Canvas{
		Id:      id,
		Name:    args.Name,
		Content: strings.Repeat(args.Fill, args.Width*args.Height),
		Width:   args.Width,
		Height:  args.Height,
	}

	s.Logger.Debug("Create::BeforeCreate", canvas.AsLogFields()...)

	err = s.Repo.Create(ctx, canvas)
	if err != nil {
		return nil, err
	}

	s.Logger.Debug("Created::Created", canvas.AsLogFields()...)
	s.Broadcast(ctx, s.BroadCaster, s.Logger, CanvasEvent{
		Name:   CanvasEventCreated,
		Canvas: canvas,
	})

	return &canvas, nil
}

// List all available Canvas items
func (s CanvasService) List(ctx context.Context) ([]Canvas, error) {
	s.Logger.Debug("List::Fetching")

	var items, err = s.Repo.List(ctx)
	if err == nil {
		s.Logger.Debug("List::Fetched", zap.Int("count", len(items)))
	} else {
		s.Logger.Error("List::Failed", zap.Error(err))
	}

	return items, err
}

// Get one specific Canvas item
func (s CanvasService) Get(ctx context.Context, id string) (*Canvas, error) {
	s.Logger.Debug("Get::Fetching")

	var canvas, err = s.Repo.Get(ctx, id)
	if err == nil {
		s.Logger.Debug("Get::Fetched", canvas.AsLogFields()...)
	} else if err == ErrNotFound {
		s.Logger.Debug("Get:NotFound", zap.String("id", id))
	} else {
		s.Logger.Error("Get::Failed", zap.Error(err))
	}

	return canvas, err
}

// Delete a specific Canvas item
func (s CanvasService) Delete(ctx context.Context, id string) error {
	s.Logger.Debug("Delete::Fetching")

	var canvas, err = s.Repo.Get(ctx, id)
	if err == nil {
		s.Logger.Debug("Delete::Fetched", canvas.AsLogFields()...)
	} else if err == ErrNotFound {
		s.Logger.Debug("Delete:NotFound", zap.String("id", id))
		return err
	} else {
		s.Logger.Error("Delete::Fetch::Failed", zap.Error(err))
		return err
	}

	err = s.Repo.Delete(ctx, id)
	if err != nil {
		s.Logger.Error("Delete::Failed", zap.Error(err))
		return err
	}

	s.Logger.Debug("Delete::Deleted", canvas.AsLogFields()...)
	s.Broadcast(ctx, s.BroadCaster, s.Logger, CanvasEvent{
		Name:   CanvasEventDeleted,
		Canvas: *canvas,
	})

	return nil
}

func (s CanvasService) Observe(ctx context.Context, id string) (StopObserveFunc, <-chan CanvasEvent, error) {
	s.Logger.Debug("Observe::Acquiring")

	if id == "" {
		id = ObserveALL
	}

	var stop, c, err = s.BroadCaster.Observe(ctx, id)
	if err == nil {
		s.Logger.Debug("Observe::Acquired")
	} else {
		s.Logger.Error("Observe::Failed", zap.Error(err))
	}

	return stop, c, err
}

type TransformRectangleArgs struct {
	TopLeft Coordinates `json:"top_left"`
	Width   int         `json:"width"`
	Height  int         `json:"height"`
	Fill    string      `json:"fill"`
	Outline string      `json:"outline"`
}

func (a TransformRectangleArgs) Validate() error {
	var errs []string

	if a.TopLeft.X < 0 {
		errs = append(errs, "TopLeft.X must not be negative")
	}

	if a.TopLeft.Y < 0 {
		errs = append(errs, "TopLeft.Y must not be negative")
	}

	if a.Width < 0 {
		errs = append(errs, "Width cannot be negative")
	}

	if a.Height < 0 {
		errs = append(errs, "Height cannot be negative")
	}

	if a.Width == 0 && a.Height == 0 {
		errs = append(errs, "At least one of Width and Height must be greater than zero")
	}

	if a.Fill == "" && a.Outline == "" {
		errs = append(errs, "Atleast one of Fill and Outline is required")
	}

	if a.Fill != "" && len(a.Fill) != 1 {
		errs = append(errs, "Fill must be not be longer than 1 character")
	}

	if a.Outline != "" && len(a.Outline) != 1 {
		errs = append(errs, "Outline must be not be longer than 1 character")
	}

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("%w: %s", ErrInvalidInput, strings.Join(errs, ", "))
}

// ApplyRectangle loads a Canvas and uses TransformRectangle on it
func (s CanvasService) ApplyRectangle(ctx context.Context, id string, args TransformRectangleArgs) (*Canvas, error) {
	s.Logger.Debug("ApplyRectangle::Fetching")

	var canvas, err = s.Repo.Get(ctx, id)
	if err == nil {
		s.Logger.Debug("ApplyRectangle::Fetched", canvas.AsLogFields()...)
	} else if err == ErrNotFound {
		s.Logger.Debug("ApplyRectangle:NotFound", zap.String("id", id))
		return nil, err
	} else {
		s.Logger.Error("ApplyRectangle::Fetch::Failed", zap.Error(err))
		return nil, err
	}

	s.Logger.Debug("ApplyRectangle::Transform")
	err = TransformRectangle(canvas, args)

	if err == nil {
		s.Logger.Debug("ApplyRectangle::Transformed", canvas.AsLogFields()...)
	} else {
		s.Logger.Error("ApplyRectangle::Transform::Failed", zap.Error(err))
		return nil, err
	}

	s.Logger.Debug("ApplyRectangle::Updating")
	err = s.Repo.Update(ctx, *canvas)

	if err != nil {
		s.Logger.Error("ApplyRectangle::Update::Failed", zap.Error(err))
		return nil, err
	}

	s.Logger.Debug("ApplyRectangle::Updated", zap.String("id", canvas.Id))
	s.Broadcast(ctx, s.BroadCaster, s.Logger, CanvasEvent{
		Name:   CanvasEventUpdated,
		Canvas: *canvas,
	})

	return canvas, nil
}

type TransformFloodfillArgs struct {
	Start Coordinates `json:"start"`
	Fill  string      `json:"fill"`
}

func (a TransformFloodfillArgs) Validate() error {
	var errs []string

	if a.Start.X < 0 {
		errs = append(errs, "Start.X must not be negative")
	}

	if a.Start.Y < 0 {
		errs = append(errs, "Start.Y must not be negative")
	}

	if a.Fill == "" || len(a.Fill) != 1 {
		errs = append(errs, "Fill must contain exactly 1 character")
	}

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("%w: %s", ErrInvalidInput, strings.Join(errs, ", "))
}

func (s CanvasService) ApplyFloodfill(ctx context.Context, id string, args TransformFloodfillArgs) (*Canvas, error) {
	s.Logger.Debug("ApplyFloodfill::Fetching")

	var canvas, err = s.Repo.Get(ctx, id)
	if err == nil {
		s.Logger.Debug("ApplyFloodfill::Fetched", canvas.AsLogFields()...)
	} else {
		s.Logger.Error("ApplyFloodfill::Fetch::Failed", zap.Error(err))
		return nil, err
	}

	s.Logger.Debug("ApplyFloodfill::Transform")
	err = TransformFloodfill(canvas, args)

	if err == nil {
		s.Logger.Debug("ApplyFloodfill::Transformed", canvas.AsLogFields()...)
	} else if err == ErrNotFound {
		s.Logger.Debug("ApplyFloodfill:NotFound", zap.String("id", id))
	} else {
		s.Logger.Error("ApplyFloodfill::Transform::Failed", zap.Error(err))
		return nil, err
	}

	s.Logger.Debug("ApplyFloodfill::Updating")
	err = s.Repo.Update(ctx, *canvas)

	if err != nil {
		s.Logger.Error("ApplyFloodfill::Update::Failed", zap.Error(err))
		return nil, err
	}

	s.Logger.Debug("ApplyFloodfill::Updated", zap.String("id", canvas.Id))
	s.Broadcast(ctx, s.BroadCaster, s.Logger, CanvasEvent{
		Name:   CanvasEventUpdated,
		Canvas: *canvas,
	})

	return canvas, err
}
