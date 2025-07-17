package mcp

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/kungfusheep/hue/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// EventManager manages event streaming for MCP
type EventManager struct {
	client        *client.Client
	stream        *client.EventStream
	recentEvents  []client.Event
	eventsMutex   sync.RWMutex
	maxEvents     int
	streaming     bool
	streamingLock sync.Mutex
}

// Global event manager instance
var eventManager *EventManager

// InitEventManager initializes the global event manager
func InitEventManager(hueClient *client.Client) {
	eventManager = &EventManager{
		client:       hueClient,
		recentEvents: make([]client.Event, 0),
		maxEvents:    1000,
	}
}

// HandleStartEventStream starts the event stream
func HandleStartEventStream(hueClient *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if eventManager == nil {
			InitEventManager(hueClient)
		}

		eventManager.streamingLock.Lock()
		defer eventManager.streamingLock.Unlock()

		if eventManager.streaming {
			return mcp.NewToolResultText("Event stream is already running"), nil
		}

		// Get filter from arguments
		args := request.GetArguments()
		filterTypes := []string{}
		if filter, ok := args["filter"].(string); ok && filter != "" {
			filterTypes = strings.Split(filter, ",")
		}

		// Start the stream
		stream, err := hueClient.StreamEvents(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to start event stream: %v", err)), nil
		}

		eventManager.stream = stream
		eventManager.streaming = true

		// Start processing events in background
		go eventManager.processEvents(filterTypes)

		result := "Event stream started successfully"
		if len(filterTypes) > 0 {
			result += fmt.Sprintf(" with filter: %s", strings.Join(filterTypes, ", "))
		}
		
		return mcp.NewToolResultText(result), nil
	}
}

// HandleStopEventStream stops the event stream
func HandleStopEventStream(hueClient *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if eventManager == nil || !eventManager.streaming {
			return mcp.NewToolResultText("Event stream is not running"), nil
		}

		eventManager.streamingLock.Lock()
		defer eventManager.streamingLock.Unlock()

		if eventManager.stream != nil {
			eventManager.stream.Close()
			eventManager.stream = nil
		}
		eventManager.streaming = false

		return mcp.NewToolResultText("Event stream stopped"), nil
	}
}

// HandleGetRecentEvents returns recent events
func HandleGetRecentEvents(hueClient *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if eventManager == nil {
			return mcp.NewToolResultText("Event stream has not been started"), nil
		}

		args := request.GetArguments()
		limit := 50 // default
		if l, ok := args["limit"].(float64); ok {
			limit = int(l)
		}

		eventType := ""
		if t, ok := args["type"].(string); ok {
			eventType = t
		}

		eventManager.eventsMutex.RLock()
		defer eventManager.eventsMutex.RUnlock()

		var result strings.Builder
		result.WriteString(fmt.Sprintf("Recent events (total stored: %d):\n\n", len(eventManager.recentEvents)))

		count := 0
		// Show events in reverse order (newest first)
		for i := len(eventManager.recentEvents) - 1; i >= 0 && count < limit; i-- {
			event := eventManager.recentEvents[i]
			
			// Filter by type if specified
			if eventType != "" && event.Type != eventType {
				continue
			}

			result.WriteString(fmt.Sprintf("🔔 Event %s at %s\n", event.ID, event.CreationTime))
			result.WriteString(fmt.Sprintf("   Type: %s\n", event.Type))
			
			for _, data := range event.Data {
				result.WriteString(fmt.Sprintf("   • %s (%s)\n", data.Type, data.ID))
				
				// Show relevant details based on type
				switch data.Type {
				case "light":
					if data.On != nil {
						result.WriteString(fmt.Sprintf("     On: %v\n", data.On.On))
					}
					if data.Dimming != nil {
						result.WriteString(fmt.Sprintf("     Brightness: %.0f%%\n", data.Dimming.Brightness))
					}
					if data.Color != nil {
						result.WriteString(fmt.Sprintf("     Color: XY(%.3f, %.3f)\n", data.Color.XY.X, data.Color.XY.Y))
					}
				case "motion":
					if data.Motion != nil {
						result.WriteString(fmt.Sprintf("     Motion: %v\n", data.Motion.Motion))
					}
				case "button":
					if data.Button != nil && data.Button.ButtonReport != nil {
						result.WriteString(fmt.Sprintf("     Button: %s\n", data.Button.ButtonReport.Event))
					}
				case "temperature":
					if data.Temperature != nil {
						result.WriteString(fmt.Sprintf("     Temperature: %.1f°C\n", data.Temperature.Temperature))
					}
				case "scene":
					if data.Status != nil {
						result.WriteString(fmt.Sprintf("     Active: %s\n", data.Status.Active))
					}
				}
			}
			result.WriteString("\n")
			count++
		}

		if count == 0 {
			result.WriteString("No events found")
			if eventType != "" {
				result.WriteString(fmt.Sprintf(" of type '%s'", eventType))
			}
		}

		return mcp.NewToolResultText(result.String()), nil
	}
}

// HandleGetEventStreamStatus returns the current streaming status
func HandleGetEventStreamStatus(hueClient *client.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var result strings.Builder
		
		result.WriteString("Event Stream Status:\n")
		
		if eventManager == nil {
			result.WriteString("• Status: Not initialized\n")
		} else {
			eventManager.streamingLock.Lock()
			streaming := eventManager.streaming
			eventManager.streamingLock.Unlock()
			
			if streaming {
				result.WriteString("• Status: Running ✅\n")
			} else {
				result.WriteString("• Status: Stopped ❌\n")
			}
			
			eventManager.eventsMutex.RLock()
			eventCount := len(eventManager.recentEvents)
			eventManager.eventsMutex.RUnlock()
			
			result.WriteString(fmt.Sprintf("• Events buffered: %d\n", eventCount))
			result.WriteString(fmt.Sprintf("• Max buffer size: %d\n", eventManager.maxEvents))
		}
		
		return mcp.NewToolResultText(result.String()), nil
	}
}

// processEvents processes incoming events
func (em *EventManager) processEvents(filterTypes []string) {
	var events <-chan client.Event
	
	if len(filterTypes) > 0 {
		events = em.stream.FilterEvents(filterTypes...)
	} else {
		events = em.stream.Events()
	}
	
	for {
		select {
		case event, ok := <-events:
			if !ok {
				return
			}
			em.storeEvent(event)
			
		case err, ok := <-em.stream.Errors():
			if !ok {
				return
			}
			// Log error but continue
			fmt.Printf("Event stream error: %v\n", err)
		}
	}
}

// storeEvent stores an event in the recent events buffer
func (em *EventManager) storeEvent(event client.Event) {
	em.eventsMutex.Lock()
	defer em.eventsMutex.Unlock()
	
	em.recentEvents = append(em.recentEvents, event)
	
	// Trim buffer if too large
	if len(em.recentEvents) > em.maxEvents {
		// Keep the most recent events
		em.recentEvents = em.recentEvents[len(em.recentEvents)-em.maxEvents:]
	}
}

// Event type constants for filtering
const (
	EventTypeLight       = "light"
	EventTypeMotion      = "motion"
	EventTypeButton      = "button"
	EventTypeTemperature = "temperature"
	EventTypeLightLevel  = "light_level"
	EventTypeScene       = "scene"
	EventTypeGroup       = "grouped_light"
	EventTypeUpdate      = "update"
	EventTypeAdd         = "add"
	EventTypeDelete      = "delete"
)