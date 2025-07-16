# Phase 1 Quick Start Guide

## 🎯 Goal: Implement Critical Features for 70%+ Coverage

### Task 1: Event Streaming (Highest Priority)

**Why**: This is THE killer feature that makes v2 superior to v1. It enables:
- Real-time state updates
- Motion sensor triggers
- Button press events
- Light state changes
- No more polling!

**Implementation Plan**:

1. **Create SSE Client** (`hue/events.go`):
```go
type EventStream struct {
    client     *Client
    events     chan Event
    errors     chan error
    done       chan bool
    reconnect  bool
}

type Event struct {
    Type      string
    Data      []EventData
    Timestamp time.Time
}

func (c *Client) StreamEvents(ctx context.Context) (*EventStream, error) {
    // Connect to /eventstream/clip/v2
    // Parse SSE format
    // Dispatch events
}
```

2. **Add MCP Event Tools** (`mcp/events.go`):
- `start_event_stream` - Begin streaming
- `stop_event_stream` - Stop streaming
- `get_recent_events` - Get buffered events

3. **Event Types to Support**:
- Light state changes
- Motion detection
- Button presses
- Scene activation
- Group updates

### Task 2: Full CRUD Operations

**Quick Wins** (Do these first):

1. **Scene Creation with State Capture**:
```go
func (c *Client) CreateSceneFromCurrentState(ctx context.Context, name string, groupID string) (*Scene, error) {
    // Get all lights in group
    // Capture their current states
    // Create scene with actions
}
```

2. **Group Membership Management**:
```go
func (c *Client) AddLightToGroup(ctx context.Context, groupID, lightID string) error
func (c *Client) RemoveLightFromGroup(ctx context.Context, groupID, lightID string) error
```

3. **Delete Operations**:
```go
func (c *Client) DeleteScene(ctx context.Context, id string) error
func (c *Client) DeleteRule(ctx context.Context, id string) error
func (c *Client) DeleteSchedule(ctx context.Context, id string) error
```

### Task 3: Entertainment Streaming

**Note**: This is complex but incredibly valuable. Consider starting with a simple proof-of-concept.

**Steps**:
1. Research DTLS handshake requirements
2. Implement basic streaming client
3. Create entertainment area activation
4. Add high-frequency update support

### Quick Implementation Order

**Day 1-2**: 
- Set up SSE client structure
- Basic event parsing
- MCP tool for starting stream

**Day 3-4**:
- Scene creation with state capture
- Group membership management
- Basic delete operations

**Day 5-7**:
- Complete event type handling
- Reconnection logic
- Event filtering and buffering

**Week 2**:
- Entertainment mode research
- Basic DTLS implementation
- Simple streaming demo

## Code Structure Suggestions

```
hue-mcp/
├── hue/
│   ├── events.go       # New: SSE client
│   ├── streaming.go    # New: Entertainment streaming
│   └── crud.go         # New: Additional CRUD operations
├── mcp/
│   ├── events.go       # New: Event stream handlers
│   └── streaming.go    # New: Entertainment handlers
└── examples/
    ├── event_monitor/  # Example: Monitor all events
    └── disco_mode/     # Example: Entertainment demo
```

## Testing Plan

1. **Event Streaming Tests**:
   - Connect and receive events
   - Handle disconnection/reconnection
   - Filter event types
   - Performance under load

2. **CRUD Tests**:
   - Create scene from state
   - Modify group membership
   - Delete resources
   - Verify cascading deletes

3. **Integration Tests**:
   - Full workflow tests
   - Error handling
   - Edge cases

## Success Criteria

- [ ] Can receive real-time events from bridge
- [ ] Can create scenes capturing current state
- [ ] Can manage group membership
- [ ] Can delete resources
- [ ] Event stream stays connected for 1+ hours
- [ ] All existing tests still pass

## Next Steps

1. Start with `hue/events.go` - create the SSE client
2. Add basic event parsing
3. Create simple MCP tool to start streaming
4. Test with button presses and motion sensors
5. Iterate and improve

Remember: Event streaming alone will make this MCP incredibly powerful for automation!