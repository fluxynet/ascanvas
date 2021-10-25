package canvas

import (
	"errors"
	"net/http"

	"github.com/fluxynet/ascanvas"
	"github.com/fluxynet/ascanvas/web"
)

type WebCanvas struct {
	Service *ascanvas.CanvasService
	GetID   web.IDGetter
}

// Create http.HandleFunc compatible handler for listing ascanvas.Canvas
// @Summary "Create all canvas items"
// @Accept json
// @Produce json
// @Param CreateArgs body ascanvas.CreateArgs true "Canvas creation details"
// @Success 201 {array} ascanvas.Canvas
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
// @Router / [post]
func (s WebCanvas) Create(w http.ResponseWriter, r *http.Request) {
	var (
		args   ascanvas.CreateArgs
		canvas *ascanvas.Canvas
		err    error

		ctx = r.Context()
	)

	err = web.ReadJsonBodyInto(r, &args)
	if err != nil {
		web.JsonError(w, http.StatusBadRequest, err)
		return
	}

	canvas, err = s.Service.Create(ctx, args)
	if err == nil {
		web.Json(w, http.StatusOK, canvas)
		return
	} else if errors.Is(err, ascanvas.ErrInvalidInput) {
		web.JsonError(w, http.StatusBadRequest, err)
	} else {
		web.JsonError(w, http.StatusInternalServerError, err)
	}
}

// List http.HandleFunc compatible handler for listing ascanvas.Canvas
// @Summary "List all canvas items"
// @Accept json
// @Produce json
// @Success 200 {array} ascanvas.Canvas
// @Failure 500 {object} web.Response
// @Router / [get]
func (s WebCanvas) List(w http.ResponseWriter, r *http.Request) {
	var items, err = s.Service.List(r.Context())
	if err == nil {
		web.Json(w, http.StatusOK, items)
		return
	}

	web.JsonError(w, http.StatusInternalServerError, err)
}

// Get http.HandleFunc compatible handler for getting a specific ascanvas.Canvas
// @Summary Get a specific canvas by id
// @Accept json
// @Produce json
// @Param id path string true "Identifier of canvas to fetch"
// @Success 200 {object} ascanvas.Canvas
// @Failure 404 {object} web.Response
// @Failure 500 {object} web.Response
// @Router /{id} [get]
func (s WebCanvas) Get(w http.ResponseWriter, r *http.Request) {
	var (
		id, err = s.GetID(r)

		canvas *ascanvas.Canvas
	)

	if err != nil {
		web.JsonError(w, http.StatusBadRequest, err)
		return
	}

	canvas, err = s.Service.Get(r.Context(), id)

	if err == nil {
		web.Json(w, http.StatusOK, canvas)
		return
	}

	web.JsonError(w, http.StatusInternalServerError, err)
}

// Delete http.HandleFunc compatible handler for deleting a specific ascanvas.Canvas
// @Summary "Delete a specific canvas item by id"
// @Accept json
// @Param id path string true "Identifier of canvas to delete"
// @Success 204
// @Failure 404 {object} web.Response
// @Failure 500 {object} web.Response
// @Router /{id} [delete]
func (s WebCanvas) Delete(w http.ResponseWriter, r *http.Request) {
	var id, err = s.GetID(r)

	if err != nil {
		web.JsonError(w, http.StatusBadRequest, err)
		return
	}

	err = s.Service.Delete(r.Context(), id)

	if err == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	web.JsonError(w, http.StatusInternalServerError, err)
}

// @Summary "Obtain an SSE live stream of all canvas events"
// @Accept json
// @Produce text/event-stream
// @Produce json
// @Success 200
// @Failure 500
// @Router /events [get]

// Observe http.HandleFunc compatible handler for server-sent events of ascanvas.Canvas
// @Summary "Obtain an SSE live stream of canvas events for a specific canvas id"
// @Accept json
// @Produce text/event-stream
// @Produce json
// @Success 200
// @Failure 500
// @Router /{id}/events [get]
// @Param id path string true "Identifier of canvas to observe"
func (s WebCanvas) Observe(w http.ResponseWriter, r *http.Request) {
	var (
		f, ok = w.(http.Flusher)
		ctx   = r.Context()

		id     string
		err    error
		stop   ascanvas.StopObserveFunc
		events <-chan ascanvas.CanvasEvent
	)

	if !ok {
		web.JsonError(w, http.StatusPreconditionFailed, web.ErrStreamingNotSupported)
	}

	id, err = s.GetID(r)
	if err != nil {
		id = ascanvas.ObserveALL
	}

	stop, events, err = s.Service.Observe(r.Context(), id)
	if err != nil {
		web.JsonError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-Type", web.ContentTypeEventStream)
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	go func(stop ascanvas.StopObserveFunc) {
		<-ctx.Done()
		stop()
	}(stop)

	for event := range events {
		_ = web.PrintJSONStream(w, f, string(event.Name), event.Canvas)
	}
}

// Rectangle http.HandleFunc compatible handler for performing ascanvas.Canvas ascanvas.TransformRectangle
// @Summary "Draw a rectangle on a specific canvas"
// @Accept json
// @Produce json
// @Param id path string true "Identifier of canvas to modify"
// @Param Transformation body ascanvas.TransformRectangleArgs true "Rectangle transformation details"
// @Success 200 {object} web.Response
// @Failure 400
// @Failure 500
// @Router /{id}/rectangle [patch]
func (s WebCanvas) Rectangle(w http.ResponseWriter, r *http.Request) {
	var (
		transformation ascanvas.TransformRectangleArgs
		canvas         *ascanvas.Canvas
		id             string
		err            error

		ctx = r.Context()
	)

	id, err = s.GetID(r)
	if err != nil {
		web.JsonError(w, http.StatusBadRequest, err)
		return
	}

	err = web.ReadJsonBodyInto(r, &transformation)
	if err != nil {
		web.JsonError(w, http.StatusBadRequest, err)
		return
	}

	canvas, err = s.Service.ApplyRectangle(ctx, id, transformation)
	if err == nil {
		web.Json(w, http.StatusOK, canvas)
		return
	} else if errors.Is(err, ascanvas.ErrInvalidInput) {
		web.JsonError(w, http.StatusBadRequest, err)
	} else {
		web.JsonError(w, http.StatusInternalServerError, err)
	}
}

// Floodfill http.HandleFunc compatible handler for performing ascanvas.Canvas ascanvas.TransformFloodfill
// @Summary "Apply flood fill on a specific canvas"
// @Accept json
// @Produce json
// @Param id path string true "Identifier of canvas to modify"
// @Param Transformation body ascanvas.TransformFloodfillArgs true "Flood fill transformation details"
// @Success 200 {object} web.Response
// @Failure 400
// @Failure 500
// @Router /{id}/floodfill [patch]
func (s WebCanvas) Floodfill(w http.ResponseWriter, r *http.Request) {
	var (
		transformation ascanvas.TransformFloodfillArgs
		canvas         *ascanvas.Canvas
		id             string
		err            error

		ctx = r.Context()
	)

	id, err = s.GetID(r)
	if err != nil {
		web.JsonError(w, http.StatusBadRequest, err)
		return
	}

	err = web.ReadJsonBodyInto(r, &transformation)
	if err != nil {
		web.JsonError(w, http.StatusBadRequest, err)
		return
	}

	canvas, err = s.Service.ApplyFloodfill(ctx, id, transformation)
	if err == nil {
		web.Json(w, http.StatusOK, canvas)
		return
	} else if errors.Is(err, ascanvas.ErrInvalidInput) {
		web.JsonError(w, http.StatusBadRequest, err)
	} else {
		web.JsonError(w, http.StatusInternalServerError, err)
	}
}
