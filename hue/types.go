package hue

// Error represents an API error
type Error struct {
	Type        string `json:"type"`
	Address     string `json:"address"`
	Description string `json:"description"`
}

// ResourceIdentifier represents a reference to another resource
type ResourceIdentifier struct {
	RID   string `json:"rid"`
	RType string `json:"rtype"`
}

// Light represents a light resource
type Light struct {
	ID       string    `json:"id"`
	IDV1     string    `json:"id_v1"`
	Type     string    `json:"type"`
	Owner    ResourceIdentifier `json:"owner"`
	Metadata Metadata  `json:"metadata"`
	On       OnState   `json:"on"`
	Dimming  Dimming   `json:"dimming"`
	Color    *Color    `json:"color,omitempty"`
	ColorTemperature *ColorTemperature `json:"color_temperature,omitempty"`
	Dynamics *Dynamics `json:"dynamics,omitempty"`
	Effects  *Effects  `json:"effects,omitempty"`
	Alert    *Alert    `json:"alert,omitempty"`
	Mode     string    `json:"mode"`
}

// Group represents a grouped light resource
type Group struct {
	ID       string    `json:"id"`
	IDV1     string    `json:"id_v1"`
	Type     string    `json:"type"`
	Owner    *ResourceIdentifier `json:"owner,omitempty"`
	Metadata Metadata  `json:"metadata"`
	On       OnState   `json:"on"`
	Dimming  Dimming   `json:"dimming"`
	Color    *Color    `json:"color,omitempty"`
	ColorTemperature *ColorTemperature `json:"color_temperature,omitempty"`
	Dynamics *Dynamics `json:"dynamics,omitempty"`
	Effects  *Effects  `json:"effects,omitempty"`
	Alert    *Alert    `json:"alert,omitempty"`
}

// Scene represents a scene resource
type Scene struct {
	ID       string              `json:"id"`
	IDV1     string              `json:"id_v1"`
	Type     string              `json:"type"`
	Metadata Metadata            `json:"metadata"`
	Group    ResourceIdentifier  `json:"group"`
	Actions  []SceneAction       `json:"actions"`
	Palette  *ScenePalette       `json:"palette,omitempty"`
	Speed    float64             `json:"speed"`
	AutoDynamic bool             `json:"auto_dynamic"`
}

// Bridge represents bridge information
type Bridge struct {
	ID          string    `json:"id"`
	IDV1        string    `json:"id_v1"`
	Type        string    `json:"type"`
	Owner       ResourceIdentifier `json:"owner"`
	BridgeID    string    `json:"bridge_id"`
	TimeZone    TimeZone  `json:"time_zone"`
}

// Metadata contains name and archetype information
type Metadata struct {
	Name      string `json:"name"`
	Archetype string `json:"archetype"`
	Image     *ResourceIdentifier `json:"image,omitempty"`
}

// OnState represents the on/off state
type OnState struct {
	On bool `json:"on"`
}

// Dimming represents brightness control
type Dimming struct {
	Brightness    float64 `json:"brightness"`
	MinDimLevel   float64 `json:"min_dim_level,omitempty"`
}

// Color represents color settings
type Color struct {
	XY    XY      `json:"xy"`
	Gamut *Gamut  `json:"gamut,omitempty"`
	GamutType string `json:"gamut_type,omitempty"`
}

// XY represents CIE xy color coordinates
type XY struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Gamut represents the color gamut
type Gamut struct {
	Red   XY `json:"red"`
	Green XY `json:"green"`
	Blue  XY `json:"blue"`
}

// ColorTemperature represents color temperature settings
type ColorTemperature struct {
	Mirek       int     `json:"mirek"`
	MirekValid  bool    `json:"mirek_valid"`
	MirekSchema *MirekSchema `json:"mirek_schema,omitempty"`
}

// MirekSchema defines the valid range for color temperature
type MirekSchema struct {
	MirekMinimum int `json:"mirek_minimum"`
	MirekMaximum int `json:"mirek_maximum"`
}

// Dynamics represents transition dynamics
type Dynamics struct {
	Status   string `json:"status,omitempty"`
	Duration int    `json:"duration,omitempty"` // milliseconds
	Speed    float64 `json:"speed,omitempty"`
}

// Effects represents dynamic effects
type Effects struct {
	Effect       string   `json:"effect"`
	StatusValues []string `json:"status_values,omitempty"`
	Status       string   `json:"status,omitempty"`
	EffectValues []string `json:"effect_values,omitempty"`
}

// Alert represents alert effects
type Alert struct {
	ActionValues []string `json:"action_values,omitempty"`
	Action       string   `json:"action"`
}

// SceneAction represents an action in a scene
type SceneAction struct {
	Target ResourceIdentifier `json:"target"`
	Action LightUpdate        `json:"action"`
}

// ScenePalette represents scene color palette
type ScenePalette struct {
	Color        []PaletteColor `json:"color"`
	Dimming      []Dimming      `json:"dimming"`
	ColorTemperature []PaletteTemperature `json:"color_temperature"`
}

// PaletteColor represents a color in the palette
type PaletteColor struct {
	Color Color   `json:"color"`
	Dimming Dimming `json:"dimming"`
}

// PaletteTemperature represents a color temperature in the palette
type PaletteTemperature struct {
	ColorTemperature ColorTemperature `json:"color_temperature"`
	Dimming          Dimming          `json:"dimming"`
}

// TimeZone represents timezone information
type TimeZone struct {
	TimeZone string `json:"time_zone"`
}

// Update structures for API calls

// LightUpdate represents an update to a light
type LightUpdate struct {
	On               *OnState          `json:"on,omitempty"`
	Dimming          *Dimming          `json:"dimming,omitempty"`
	Color            *Color            `json:"color,omitempty"`
	ColorTemperature *ColorTemperature `json:"color_temperature,omitempty"`
	Dynamics         *Dynamics         `json:"dynamics,omitempty"`
	Effects          *Effects          `json:"effects,omitempty"`
	Alert            *Alert            `json:"alert,omitempty"`
}

// GroupUpdate represents an update to a group
type GroupUpdate struct {
	On               *OnState          `json:"on,omitempty"`
	Dimming          *Dimming          `json:"dimming,omitempty"`
	Color            *Color            `json:"color,omitempty"`
	ColorTemperature *ColorTemperature `json:"color_temperature,omitempty"`
	Dynamics         *Dynamics         `json:"dynamics,omitempty"`
	Effects          *Effects          `json:"effects,omitempty"`
	Alert            *Alert            `json:"alert,omitempty"`
}

// SceneCreate represents a scene creation request
type SceneCreate struct {
	Type     string             `json:"type"`
	Metadata Metadata           `json:"metadata"`
	Group    ResourceIdentifier `json:"group"`
	Actions  []SceneAction      `json:"actions"`
}