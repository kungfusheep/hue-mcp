package client

import (
	"context"
	"fmt"
)

// Room represents a room resource
type Room struct {
	ID       string              `json:"id"`
	IDV1     string              `json:"id_v1"`
	Type     string              `json:"type"`
	Services []ResourceIdentifier `json:"services"`
	Metadata Metadata            `json:"metadata"`
	Children []ResourceIdentifier `json:"children"`
}

// Zone represents a zone resource  
type Zone struct {
	ID       string              `json:"id"`
	IDV1     string              `json:"id_v1"`
	Type     string              `json:"type"`
	Services []ResourceIdentifier `json:"services"`
	Metadata Metadata            `json:"metadata"`
	Children []ResourceIdentifier `json:"children"`
}

// GetRooms returns all rooms
func (c *Client) GetRooms(ctx context.Context) ([]Room, error) {
	var response struct {
		Errors []Error `json:"errors"`
		Data   []Room  `json:"data"`
	}
	
	err := c.getJSON(ctx, "/resource/room", &response)
	if err != nil {
		return nil, err
	}
	
	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s", response.Errors[0].Description)
	}
	
	return response.Data, nil
}

// GetRoom returns a specific room
func (c *Client) GetRoom(ctx context.Context, id string) (*Room, error) {
	var response struct {
		Errors []Error `json:"errors"`
		Data   []Room  `json:"data"`
	}
	
	err := c.getJSON(ctx, fmt.Sprintf("/resource/room/%s", id), &response)
	if err != nil {
		return nil, err
	}
	
	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s", response.Errors[0].Description)
	}
	
	if len(response.Data) == 0 {
		return nil, fmt.Errorf("room not found")
	}
	
	return &response.Data[0], nil
}

// GetZones returns all zones
func (c *Client) GetZones(ctx context.Context) ([]Zone, error) {
	var response struct {
		Errors []Error `json:"errors"`
		Data   []Zone  `json:"data"`
	}
	
	err := c.getJSON(ctx, "/resource/zone", &response)
	if err != nil {
		return nil, err
	}
	
	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s", response.Errors[0].Description)
	}
	
	return response.Data, nil
}

// GetZone returns a specific zone
func (c *Client) GetZone(ctx context.Context, id string) (*Zone, error) {
	var response struct {
		Errors []Error `json:"errors"`
		Data   []Zone  `json:"data"`
	}
	
	err := c.getJSON(ctx, fmt.Sprintf("/resource/zone/%s", id), &response)
	if err != nil {
		return nil, err
	}
	
	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s", response.Errors[0].Description)
	}
	
	if len(response.Data) == 0 {
		return nil, fmt.Errorf("zone not found")
	}
	
	return &response.Data[0], nil
}