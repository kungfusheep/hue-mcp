package hue

import (
	"context"
	"fmt"
)

// Motion sensor types

// Motion represents a motion sensor resource
type Motion struct {
	ID       string             `json:"id"`
	IDV1     string             `json:"id_v1"`
	Type     string             `json:"type"`
	Owner    ResourceIdentifier `json:"owner"`
	Enabled  bool               `json:"enabled"`
	Motion   MotionReport       `json:"motion"`
}

// MotionReport contains motion detection data
type MotionReport struct {
	Motion      bool   `json:"motion"`
	MotionValid bool   `json:"motion_valid"`
	MotionReport *struct {
		Changed string `json:"changed"`
		Motion  bool   `json:"motion"`
	} `json:"motion_report,omitempty"`
}

// Temperature sensor types

// Temperature represents a temperature sensor resource
type Temperature struct {
	ID          string             `json:"id"`
	IDV1        string             `json:"id_v1"`
	Type        string             `json:"type"`
	Owner       ResourceIdentifier `json:"owner"`
	Enabled     bool               `json:"enabled"`
	Temperature TemperatureReport  `json:"temperature"`
}

// TemperatureReport contains temperature data
type TemperatureReport struct {
	Temperature      float64 `json:"temperature"`
	TemperatureValid bool    `json:"temperature_valid"`
	TemperatureReport *struct {
		Changed     string  `json:"changed"`
		Temperature float64 `json:"temperature"`
	} `json:"temperature_report,omitempty"`
}

// Light level sensor types

// LightLevel represents a light level sensor resource
type LightLevel struct {
	ID         string             `json:"id"`
	IDV1       string             `json:"id_v1"`
	Type       string             `json:"type"`
	Owner      ResourceIdentifier `json:"owner"`
	Enabled    bool               `json:"enabled"`
	LightLevel LightLevelReport   `json:"light"`
}

// LightLevelReport contains light level data
type LightLevelReport struct {
	LightLevel      int  `json:"light_level"`
	LightLevelValid bool `json:"light_level_valid"`
	LightLevelReport *struct {
		Changed    string `json:"changed"`
		LightLevel int    `json:"light_level"`
	} `json:"light_level_report,omitempty"`
}

// Button types

// Button represents a button resource (like dimmer switches)
type Button struct {
	ID         string              `json:"id"`
	IDV1       string              `json:"id_v1"`
	Type       string              `json:"type"`
	Owner      ResourceIdentifier  `json:"owner"`
	Metadata   Metadata            `json:"metadata"`
	Button     ButtonReport        `json:"button"`
}

// ButtonReport contains button state
type ButtonReport struct {
	ButtonReport *struct {
		Updated string `json:"updated"`
		Event   string `json:"event"`
	} `json:"button_report,omitempty"`
	RepeatInterval int    `json:"repeat_interval"`
	EventValues    []string `json:"event_values"`
}

// GetMotionSensors returns all motion sensors
func (c *Client) GetMotionSensors(ctx context.Context) ([]Motion, error) {
	var response struct {
		Errors []Error  `json:"errors"`
		Data   []Motion `json:"data"`
	}
	
	err := c.getJSON(ctx, "/resource/motion", &response)
	if err != nil {
		return nil, err
	}
	
	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s", response.Errors[0].Description)
	}
	
	return response.Data, nil
}

// GetTemperatureSensors returns all temperature sensors
func (c *Client) GetTemperatureSensors(ctx context.Context) ([]Temperature, error) {
	var response struct {
		Errors []Error       `json:"errors"`
		Data   []Temperature `json:"data"`
	}
	
	err := c.getJSON(ctx, "/resource/temperature", &response)
	if err != nil {
		return nil, err
	}
	
	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s", response.Errors[0].Description)
	}
	
	return response.Data, nil
}

// GetLightLevelSensors returns all light level sensors
func (c *Client) GetLightLevelSensors(ctx context.Context) ([]LightLevel, error) {
	var response struct {
		Errors []Error      `json:"errors"`
		Data   []LightLevel `json:"data"`
	}
	
	err := c.getJSON(ctx, "/resource/light_level", &response)
	if err != nil {
		return nil, err
	}
	
	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s", response.Errors[0].Description)
	}
	
	return response.Data, nil
}

// GetButtons returns all buttons (dimmer switches, etc)
func (c *Client) GetButtons(ctx context.Context) ([]Button, error) {
	var response struct {
		Errors []Error  `json:"errors"`
		Data   []Button `json:"data"`
	}
	
	err := c.getJSON(ctx, "/resource/button", &response)
	if err != nil {
		return nil, err
	}
	
	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("API error: %s", response.Errors[0].Description)
	}
	
	return response.Data, nil
}