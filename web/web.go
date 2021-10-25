package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/fluxynet/ascanvas/internal"
)

const (
	// ContentTypeJSON is the content type for JSON
	ContentTypeJSON = "application/json"

	// ContentTypeHTML is the content type for HTML
	ContentTypeHTML = "text/html"

	// ContentTypeEventStream used for SSE
	ContentTypeEventStream = "text/event-stream"
)

// IDGetter gets id from a request
type IDGetter func(r *http.Request) (string, error)

var (
	// ErrInvalidRequest means a request is either nil or not appropriate for the requested action
	ErrInvalidRequest = errors.New("request is invalid")

	// ErrPayloadUnverified payload could not be verified wrt signature
	ErrPayloadUnverified = errors.New("payload could not be verified")

	// ErrStreamingNotSupported means the browser does not support SSE streaming
	ErrStreamingNotSupported = errors.New("payload could not be verified")

	// ErrIDMissing from query string
	ErrIDMissing = errors.New("id missing from request")
)

// Print sends data to the browser
func Print(w http.ResponseWriter, status int, ctype string, content []byte) {
	if ctype != "" {
		w.Header().Set("Content-Type", ctype)
	}

	w.WriteHeader(status)
	w.Write(content)
}

// Json to the browser
func Json(w http.ResponseWriter, status int, r interface{}) {
	var b, err = json.Marshal(r)
	if err == nil {
		Print(w, status, ContentTypeJSON, b)
	} else {
		JsonError(w, http.StatusInternalServerError, err)
	}
}

func PrintStream(w io.Writer, f http.Flusher, name string, data []byte) error {
	var err error

	name = strings.ReplaceAll(name, "\n", "_")
	if name != "" {
		_, err = fmt.Fprintf(w, "event: %s\n", name)
	}

	if err == nil {
		_, err = fmt.Fprintf(w, "data: %s\n\n", bytes.ReplaceAll(data, []byte("\n"), []byte("")))
	}

	if err == nil {
		f.Flush()
	}

	return err
}

func PrintJSONStream(w io.Writer, f http.Flusher, name string, data interface{}) error {
	var b, err = json.Marshal(data)
	if err != nil {
		return err
	}

	return PrintStream(w, f, name, b)
}

// JsonError to the browser in json format
func JsonError(w http.ResponseWriter, status int, err error) {
	var m string

	if err != nil {
		m = strings.ReplaceAll(err.Error(), `"`, `\"`)
	}

	Print(w, status, ContentTypeJSON, []byte(`{"error":"`+m+`"}`))
}

// ReadBody from an http.Request
func ReadBody(r *http.Request) ([]byte, error) {
	if r == nil {
		return nil, ErrInvalidRequest
	}

	switch r.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		break
	default:
		return nil, ErrInvalidRequest
	}

	if r.Body == nil {
		return nil, nil
	}

	defer internal.Closed(r.Body)
	var b, err = io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// ReadJsonBodyInto reads json body into a target structure
func ReadJsonBodyInto(r *http.Request, target interface{}) error {
	var b, err = ReadBody(r)

	if err != nil {
		return err
	}

	return json.Unmarshal(b, target)
}

// Response is a generic reply
type Response struct {
	Message string `json:"message"`
}

// StaticIDGetter always returns the same id, err combination. useful for testing
func StaticIDGetter(id string, err error) IDGetter {
	return func(r *http.Request) (string, error) {
		return id, err
	}
}

// ChiIDGetter for chi mux library
func ChiIDGetter(r *http.Request) (string, error) {
	var id = chi.URLParam(r, "id")

	if id == "" {
		return "", ErrIDMissing
	}

	return id, nil
}
