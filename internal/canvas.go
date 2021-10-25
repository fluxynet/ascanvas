package internal

import (
	"encoding/json"
	"strings"

	"github.com/fluxynet/ascanvas"
)

func CanvasFromText(id, name string, s string) *ascanvas.Canvas {
	var canvas = ascanvas.Canvas{
		Id:   id,
		Name: name,
	}

	s = strings.Trim(s, "\n")
	l := strings.Split(s, "\n")

	canvas.Height = len(l)
	if canvas.Height > 0 {
		canvas.Width = len(l[0])
	}

	canvas.Content = strings.Join(l, "")

	return &canvas
}

func CanvasJsonFromText(id string, name string, s string) string {
	var (
		b   strings.Builder
		enc = json.NewEncoder(&b)
	)
	c := CanvasFromText(id, name, s)

	_ = enc.Encode(c)

	return strings.Trim(b.String(), "\n")
}
