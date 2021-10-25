package memory_test

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/fluxynet/ascanvas"
	"github.com/fluxynet/ascanvas/broadcaster/memory"
)

type receiver struct {
	ObserveId string
	Want      []ascanvas.CanvasEvent
	Events    []ascanvas.CanvasEvent
	Ctx       context.Context
	Closed    bool
}

func (r *receiver) receive(events <-chan ascanvas.CanvasEvent) {
	for e := range events {
		r.Events = append(r.Events, e)
	}
	r.Closed = true
}

func TestMemory(t *testing.T) {
	tests := []struct {
		name      string
		events    []ascanvas.CanvasEvent
		receivers []receiver
	}{
		{
			name: "1 listener id = 1, 1 event id 1",
			events: []ascanvas.CanvasEvent{
				{
					Name: ascanvas.CanvasEventCreated,
					Canvas: ascanvas.Canvas{
						Id:      "1",
						Name:    "C1",
						Content: "XY",
						Width:   1,
						Height:  2,
					},
				},
			},
			receivers: []receiver{
				{
					ObserveId: "1",
					Ctx:       context.Background(),
					Want: []ascanvas.CanvasEvent{
						{
							Name: ascanvas.CanvasEventCreated,
							Canvas: ascanvas.Canvas{
								Id:      "1",
								Name:    "C1",
								Content: "XY",
								Width:   1,
								Height:  2,
							},
						},
					},
				},
			},
		},
		{
			name: "1 listener id = 1, 1 event id 2",
			events: []ascanvas.CanvasEvent{
				{
					Name: ascanvas.CanvasEventUpdated,
					Canvas: ascanvas.Canvas{
						Id:      "2",
						Name:    "C2",
						Content: "QWERTY",
						Width:   3,
						Height:  2,
					},
				},
			},
			receivers: []receiver{
				{
					ObserveId: "",
					Ctx:       context.Background(),
					Want:      []ascanvas.CanvasEvent{},
				},
			},
		},
		{
			name: "1 listener id = 2, 2 events id 1, id 2",
			events: []ascanvas.CanvasEvent{
				{
					Name: ascanvas.CanvasEventDeleted,
					Canvas: ascanvas.Canvas{
						Id:      "1",
						Name:    "C1",
						Content: "XY",
						Width:   1,
						Height:  2,
					},
				},
				{
					Name: ascanvas.CanvasEventCreated,
					Canvas: ascanvas.Canvas{
						Id:      "2",
						Name:    "C2",
						Content: "ABCD",
						Width:   2,
						Height:  2,
					},
				},
			},
			receivers: []receiver{
				{
					ObserveId: "2",
					Ctx:       context.Background(),
					Want: []ascanvas.CanvasEvent{
						{
							Name: ascanvas.CanvasEventCreated,
							Canvas: ascanvas.Canvas{
								Id:      "2",
								Name:    "C2",
								Content: "ABCD",
								Width:   2,
								Height:  2,
							},
						},
					},
				},
			},
		},
		{
			name: "3 listeners: all, id = 1, id = 2; 2 events id 1, id 2",
			events: []ascanvas.CanvasEvent{
				{
					Name: ascanvas.CanvasEventDeleted,
					Canvas: ascanvas.Canvas{
						Id:      "1",
						Name:    "C1",
						Content: "XY",
						Width:   1,
						Height:  2,
					},
				},
				{
					Name: ascanvas.CanvasEventCreated,
					Canvas: ascanvas.Canvas{
						Id:      "2",
						Name:    "C2",
						Content: "ABCD",
						Width:   2,
						Height:  2,
					},
				},
			},
			receivers: []receiver{
				{
					ObserveId: ascanvas.ObserveALL,
					Ctx:       context.Background(),
					Want: []ascanvas.CanvasEvent{
						{
							Name: ascanvas.CanvasEventDeleted,
							Canvas: ascanvas.Canvas{
								Id:      "1",
								Name:    "C1",
								Content: "XY",
								Width:   1,
								Height:  2,
							},
						},
						{
							Name: ascanvas.CanvasEventCreated,
							Canvas: ascanvas.Canvas{
								Id:      "2",
								Name:    "C2",
								Content: "ABCD",
								Width:   2,
								Height:  2,
							},
						},
					},
				},
				{
					ObserveId: "1",
					Ctx:       context.Background(),
					Want: []ascanvas.CanvasEvent{
						{
							Name: ascanvas.CanvasEventDeleted,
							Canvas: ascanvas.Canvas{
								Id:      "1",
								Name:    "C1",
								Content: "XY",
								Width:   1,
								Height:  2,
							},
						},
					},
				},
				{
					ObserveId: "2",
					Ctx:       context.Background(),
					Want: []ascanvas.CanvasEvent{
						{
							Name: ascanvas.CanvasEventCreated,
							Canvas: ascanvas.Canvas{
								Id:      "2",
								Name:    "C2",
								Content: "ABCD",
								Width:   2,
								Height:  2,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			mem := memory.New()

			wg.Add(len(tt.receivers))

			for i := range tt.receivers {
				_, c, err := mem.Observe(tt.receivers[i].Ctx, tt.receivers[i].ObserveId)

				if err != nil {
					t.Errorf("Observe() error = %v", err)
					return
				}

				go func(i int) {
					defer wg.Done()
					tt.receivers[i].receive(c)
				}(i)
			}

			for i := range tt.events {
				err := mem.Broadcast(context.Background(), tt.events[i])

				if err != nil {
					t.Errorf("Broadcast() error = %v", err)
					return
				}
			}

			err := mem.Close()
			if err != nil {
				t.Errorf("Close() error = %v", err)
				return
			}

			wg.Wait()

			for i := range tt.receivers {
				if !tt.receivers[i].Closed {
					t.Errorf("receiver not closed #%d", i)
					return
				}

				if len(tt.receivers[i].Events) == 0 && len(tt.receivers[i].Want) == 0 {
					// okay
				} else if !reflect.DeepEqual(tt.receivers[i].Events, tt.receivers[i].Want) {
					t.Errorf("receivers got = %v, want %v", tt.receivers[i].Events, tt.receivers[i].Want)
					return
				}
			}
		})
	}
}
