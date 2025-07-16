package hue

import (
	"context"
	"fmt"
)

// Entertainment represents an entertainment configuration
type Entertainment struct {
	ID               string                    `json:"id"`
	IDV1             string                    `json:"id_v1"`
	Type             string                    `json:"type"`
	Metadata         Metadata                  `json:"metadata"`
	ConfigurationType string                   `json:"configuration_type"`
	Status           string                    `json:"status"`
	ActiveStreamer   *ResourceIdentifier       `json:"active_streamer,omitempty"`
	StreamProxy      StreamProxy               `json:"stream_proxy"`
	Channels         []EntertainmentChannel    `json:"channels"`
	Locations        *EntertainmentLocations   `json:"locations,omitempty"`
	LightServices    []ResourceIdentifier      `json:"light_services"`
}

// StreamProxy contains streaming proxy information
type StreamProxy struct {
	Mode string `json:"mode"`
	Node ResourceIdentifier `json:"node"`
}

// EntertainmentChannel represents a channel configuration
type EntertainmentChannel struct {
	ChannelID     int                    `json:"channel_id"`
	Position      EntertainmentPosition  `json:"position"`
	Members       []ChannelMember        `json:"members"`
}

// ChannelMember represents a light in an entertainment channel
type ChannelMember struct {
	Service ResourceIdentifier `json:"service"`
	Index   int                `json:"index"`
}

// EntertainmentPosition represents a 3D position
type EntertainmentPosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// EntertainmentLocations contains location bounds
type EntertainmentLocations struct {
	ServiceLocations []ServiceLocation `json:"service_locations"`
}

// ServiceLocation represents a service's physical location
type ServiceLocation struct {
	Service   ResourceIdentifier    `json:"service"`
	Position  EntertainmentPosition `json:"position"`
	Positions []EntertainmentPosition `json:"positions,omitempty"`
}

// GetEntertainmentConfigurations returns all entertainment configurations
func (c *Client) GetEntertainmentConfigurations(ctx context.Context) ([]Entertainment, error) {
	var response struct {
		Errors []Error         `json:"errors"`
		Data   []Entertainment `json:"data"`
	}
	
	err := c.getJSON(ctx, "/resource/entertainment_configuration", &response)
	if err != nil {
		return nil, err
	}
	
	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s", response.Errors[0].Description)
	}
	
	return response.Data, nil
}

// GetEntertainmentConfiguration returns a specific entertainment configuration
func (c *Client) GetEntertainmentConfiguration(ctx context.Context, id string) (*Entertainment, error) {
	var response struct {
		Errors []Error         `json:"errors"`
		Data   []Entertainment `json:"data"`
	}
	
	err := c.getJSON(ctx, fmt.Sprintf("/resource/entertainment_configuration/%s", id), &response)
	if err != nil {
		return nil, err
	}
	
	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s", response.Errors[0].Description)
	}
	
	if len(response.Data) == 0 {
		return nil, fmt.Errorf("entertainment configuration not found")
	}
	
	return &response.Data[0], nil
}

// StartEntertainment starts entertainment mode
func (c *Client) StartEntertainment(ctx context.Context, id string) error {
	update := map[string]interface{}{
		"action": "start",
	}
	_, err := c.put(ctx, fmt.Sprintf("/resource/entertainment_configuration/%s", id), update)
	return err
}

// StopEntertainment stops entertainment mode
func (c *Client) StopEntertainment(ctx context.Context, id string) error {
	update := map[string]interface{}{
		"action": "stop",
	}
	_, err := c.put(ctx, fmt.Sprintf("/resource/entertainment_configuration/%s", id), update)
	return err
}