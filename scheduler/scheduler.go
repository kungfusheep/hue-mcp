package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kungfusheep/hue-mcp/hue"
)

// Command represents a scheduled command
type Command struct {
	Type      string                 // "light", "group", "scene", etc.
	Action    string                 // "on", "off", "color", "brightness", etc.
	Target    string                 // ID of the target (light, group, etc.)
	Params    map[string]interface{} // Additional parameters
	Delay     time.Duration          // Delay before executing this command
}

// Sequence represents a sequence of commands
type Sequence struct {
	ID       string
	Name     string
	Commands []Command
	Loop     bool          // Whether to loop the sequence
	Running  bool
	stopChan chan struct{}
}

// Scheduler manages scheduled lighting operations
type Scheduler struct {
	client    *hue.Client
	sequences map[string]*Sequence
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewScheduler creates a new scheduler
func NewScheduler(client *hue.Client) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		client:    client,
		sequences: make(map[string]*Sequence),
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Stop stops all sequences and shuts down the scheduler
func (s *Scheduler) Stop() {
	s.cancel()
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	for _, seq := range s.sequences {
		if seq.Running && seq.stopChan != nil {
			close(seq.stopChan)
		}
	}
}

// ExecuteCommand executes a single command asynchronously
func (s *Scheduler) ExecuteCommand(cmd Command) error {
	go func() {
		// Apply delay if specified
		if cmd.Delay > 0 {
			select {
			case <-time.After(cmd.Delay):
			case <-s.ctx.Done():
				return
			}
		}
		
		// Execute the command
		ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
		defer cancel()
		
		s.executeCommandSync(ctx, cmd)
	}()
	
	return nil
}

// ExecuteSequence starts executing a sequence asynchronously
func (s *Scheduler) ExecuteSequence(seq *Sequence) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if seq.ID == "" {
		seq.ID = fmt.Sprintf("seq_%d", time.Now().UnixNano())
	}
	
	if _, exists := s.sequences[seq.ID]; exists && s.sequences[seq.ID].Running {
		return "", fmt.Errorf("sequence %s is already running", seq.ID)
	}
	
	seq.Running = true
	seq.stopChan = make(chan struct{})
	s.sequences[seq.ID] = seq
	
	// Start the sequence in a goroutine
	go s.runSequence(seq)
	
	return seq.ID, nil
}

// StopSequence stops a running sequence
func (s *Scheduler) StopSequence(sequenceID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	seq, exists := s.sequences[sequenceID]
	if !exists {
		return fmt.Errorf("sequence %s not found", sequenceID)
	}
	
	if seq.Running && seq.stopChan != nil {
		close(seq.stopChan)
		seq.Running = false
	}
	
	return nil
}

// GetSequences returns all sequences
func (s *Scheduler) GetSequences() map[string]*Sequence {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Return a copy to avoid concurrent modification
	result := make(map[string]*Sequence)
	for k, v := range s.sequences {
		result[k] = v
	}
	return result
}

// runSequence executes a sequence of commands
func (s *Scheduler) runSequence(seq *Sequence) {
	defer func() {
		s.mu.Lock()
		seq.Running = false
		s.mu.Unlock()
	}()
	
	for {
		for _, cmd := range seq.Commands {
			// Check if we should stop
			select {
			case <-seq.stopChan:
				return
			case <-s.ctx.Done():
				return
			default:
			}
			
			// Apply delay if specified
			if cmd.Delay > 0 {
				select {
				case <-time.After(cmd.Delay):
				case <-seq.stopChan:
					return
				case <-s.ctx.Done():
					return
				}
			}
			
			// Execute the command
			ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
			s.executeCommandSync(ctx, cmd)
			cancel()
		}
		
		// If not looping, we're done
		if !seq.Loop {
			break
		}
	}
}

// executeCommandSync executes a command synchronously
func (s *Scheduler) executeCommandSync(ctx context.Context, cmd Command) error {
	switch cmd.Type {
	case "light":
		return s.executeLightCommand(ctx, cmd)
	case "group":
		return s.executeGroupCommand(ctx, cmd)
	case "scene":
		return s.executeSceneCommand(ctx, cmd)
	default:
		return fmt.Errorf("unknown command type: %s", cmd.Type)
	}
}

// executeLightCommand executes a light command
func (s *Scheduler) executeLightCommand(ctx context.Context, cmd Command) error {
	switch cmd.Action {
	case "on":
		return s.client.TurnOnLight(ctx, cmd.Target)
	case "off":
		return s.client.TurnOffLight(ctx, cmd.Target)
	case "brightness":
		if brightness, ok := cmd.Params["brightness"].(float64); ok {
			return s.client.SetLightBrightness(ctx, cmd.Target, brightness)
		}
		return fmt.Errorf("brightness parameter required")
	case "color":
		if color, ok := cmd.Params["color"].(string); ok {
			return s.client.SetLightColor(ctx, cmd.Target, color)
		}
		return fmt.Errorf("color parameter required")
	default:
		return fmt.Errorf("unknown light action: %s", cmd.Action)
	}
}

// executeGroupCommand executes a group command
func (s *Scheduler) executeGroupCommand(ctx context.Context, cmd Command) error {
	switch cmd.Action {
	case "on":
		return s.client.TurnOnGroup(ctx, cmd.Target)
	case "off":
		return s.client.TurnOffGroup(ctx, cmd.Target)
	case "brightness":
		if brightness, ok := cmd.Params["brightness"].(float64); ok {
			return s.client.SetGroupBrightness(ctx, cmd.Target, brightness)
		}
		return fmt.Errorf("brightness parameter required")
	case "color":
		if color, ok := cmd.Params["color"].(string); ok {
			return s.client.SetGroupColor(ctx, cmd.Target, color)
		}
		return fmt.Errorf("color parameter required")
	default:
		return fmt.Errorf("unknown group action: %s", cmd.Action)
	}
}

// executeSceneCommand executes a scene command
func (s *Scheduler) executeSceneCommand(ctx context.Context, cmd Command) error {
	if cmd.Action == "recall" || cmd.Action == "activate" {
		return s.client.ActivateScene(ctx, cmd.Target)
	}
	return fmt.Errorf("unknown scene action: %s", cmd.Action)
}

