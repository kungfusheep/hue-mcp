package hue

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// EventStream represents a connection to the Hue event stream
type EventStream struct {
	client    *Client
	events    chan Event
	errors    chan error
	done      chan bool
	reconnect bool
}

// Event represents a Hue v2 event
type Event struct {
	CreationTime string      `json:"creationtime"`
	ID           string      `json:"id"`
	Type         string      `json:"type"`
	Data         []EventData `json:"data"`
}

// EventData represents the data portion of an event
type EventData struct {
	ID     string      `json:"id"`
	IDV1   string      `json:"id_v1,omitempty"`
	Type   string      `json:"type"`
	Owner  *ResourceIdentifier `json:"owner,omitempty"`
	
	// Light events
	On               *On               `json:"on,omitempty"`
	Dimming          *Dimming          `json:"dimming,omitempty"`
	Color            *Color            `json:"color,omitempty"`
	ColorTemperature *ColorTemperature `json:"color_temperature,omitempty"`
	Effects          *Effects          `json:"effects,omitempty"`
	
	// Motion sensor events
	Motion *MotionReport `json:"motion,omitempty"`
	
	// Temperature sensor events
	Temperature *TemperatureReport `json:"temperature,omitempty"`
	
	// Light level sensor events
	Light *LightLevelReport `json:"light,omitempty"`
	
	// Button events
	Button *ButtonReport `json:"button,omitempty"`
	
	// Scene events
	Status *struct {
		Active string `json:"active"`
	} `json:"status,omitempty"`
	
	// Group events
	Alert *Alert `json:"alert,omitempty"`
}

// StreamEvents creates a new event stream connection
func (c *Client) StreamEvents(ctx context.Context) (*EventStream, error) {
	stream := &EventStream{
		client:    c,
		events:    make(chan Event, 100),
		errors:    make(chan error, 10),
		done:      make(chan bool),
		reconnect: true,
	}
	
	go stream.connect(ctx)
	
	return stream, nil
}

// Events returns the event channel
func (es *EventStream) Events() <-chan Event {
	return es.events
}

// Errors returns the error channel
func (es *EventStream) Errors() <-chan error {
	return es.errors
}

// Close stops the event stream
func (es *EventStream) Close() {
	es.reconnect = false
	close(es.done)
}

// connect establishes and maintains the SSE connection
func (es *EventStream) connect(ctx context.Context) {
	defer close(es.events)
	defer close(es.errors)
	
	for es.reconnect {
		select {
		case <-ctx.Done():
			return
		case <-es.done:
			return
		default:
			err := es.streamEvents(ctx)
			if err != nil {
				es.errors <- fmt.Errorf("stream error: %w", err)
				if es.reconnect {
					// Wait before reconnecting
					select {
					case <-time.After(5 * time.Second):
						continue
					case <-ctx.Done():
						return
					case <-es.done:
						return
					}
				}
			}
		}
	}
}

// streamEvents handles the actual SSE connection
func (es *EventStream) streamEvents(ctx context.Context) error {
	url := fmt.Sprintf("https://%s/eventstream/clip/v2", es.client.bridgeIP)
	
	req, err := es.client.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	
	// SSE requires these headers
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	
	resp, err := es.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	
	scanner := bufio.NewScanner(resp.Body)
	var eventData strings.Builder
	
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-es.done:
			return nil
		default:
			line := scanner.Text()
			
			if line == "" {
				// Empty line signals end of event
				if eventData.Len() > 0 {
					es.processEvent(eventData.String())
					eventData.Reset()
				}
				continue
			}
			
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				eventData.WriteString(data)
			} else if strings.HasPrefix(line, ": hi") {
				// Keepalive message, ignore
				continue
			}
		}
	}
	
	if err := scanner.Err(); err != nil {
		return err
	}
	
	return nil
}

// processEvent parses and sends an event
func (es *EventStream) processEvent(data string) {
	var events []Event
	if err := json.Unmarshal([]byte(data), &events); err != nil {
		es.errors <- fmt.Errorf("failed to parse event: %w", err)
		return
	}
	
	for _, event := range events {
		select {
		case es.events <- event:
		default:
			// Channel full, drop oldest event
			select {
			case <-es.events:
				es.events <- event
			default:
			}
		}
	}
}

// FilterEvents creates a filtered event stream
func (es *EventStream) FilterEvents(types ...string) <-chan Event {
	filtered := make(chan Event, 100)
	typeMap := make(map[string]bool)
	for _, t := range types {
		typeMap[t] = true
	}
	
	go func() {
		defer close(filtered)
		for event := range es.events {
			if typeMap[event.Type] || len(types) == 0 {
				filtered <- event
			}
		}
	}()
	
	return filtered
}