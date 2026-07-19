# MapleStory Private Server (FM)

A high-performance MapleStory private server implementation written in Go, featuring a custom thread-based architecture designed for scalability and efficiency.

## Architecture Overview

This server implementation uses a custom IO/Logic thread pool architecture instead of traditional actor models, providing:

- **Configurable Thread Pools**: Separate IO and Logic thread counts for optimal resource utilization
- **Dynamic Routing**: Intelligent task routing with rerouting capabilities
- **Thread-Per-Logic Lua States**: Isolated Lua script execution per logic thread
- **Type-Safe Packet Handling**: Generic-based packet system for compile-time safety
- **Efficient Resource Management**: Memory-conscious design with proper resource cleanup

### Key Components

- **Login Server**: Handles authentication, character selection, and server selection
- **Game Server**: Manages gameplay, character interactions, and world state

### MapActor messages

- Handlers are registered by message type (`services/game/actor/message_handlers.go`); `MapActor.Receive` dispatches through that registry after `EnsureDeliver` unwrap.
- Lua map builtins use the private `luaMapCall` intermediary. Synchronous work runs through the registered `MapCallHandler` immediately when the caller is already on the target map actor; cross-actor work is sent as `MapCall` and resumes through `MapCallAck`. Promise-based work uses `MapCallAsync` and always yields. Core systems do not use this intermediary and send actor messages directly.

## Project Structure

```
fm/
├── core/           # Core game logic and entity management
├── entity/         # Game entities (characters, mobs, npcs, items)
├── protocol/       # Network protocol definitions and packet structures
├── server/         # Server implementations (login, game)
├── common/         # Shared utilities and infrastructure
├── login/          # Login server implementation
├── game/           # Game server implementation
└── legacy/         # Original actor-based implementation
```

## Import Cycle Prevention Rules

**CRITICAL**: To prevent import cycles, the following rules must be strictly followed:

### Package Dependency Rules

1. **core package**:
   - ✅ Can import: `common`, `protocol`
   - ❌ **NEVER** import: `entity`, `server`

2. **entity package**:
   - ✅ Can import: `common`, `core`
   - ❌ **NEVER** import: `protocol`, `server`

3. **protocol package**:
   - ✅ Can import: `common`
   - ❌ **NEVER** import: `core`, `entity`, `server`

4. **server package**:
   - ✅ Can import: `common`, `core`, `entity`, `protocol`
   - ✅ Can implement interfaces defined in other packages

### Interface Design Pattern

- **Listener interfaces** are defined in `entity` package
- **Listener implementations** are implemented in `server` package
- This prevents circular dependencies while allowing proper separation of concerns

### Violation Consequences

Breaking these rules will result in:
- Import cycle compilation errors
- Circular dependency issues
- Build failures

**Always verify package imports before committing changes!**

## Configuration

The `config.yaml` file contains all server configuration:

```yaml
login:
  port: 7100
  io_threads: 2      # Number of IO threads
  logic_threads: 4   # Number of logic threads

game:
  port: 7111
  io_threads: 4      # Number of IO threads
  logic_threads: 8   # Number of logic threads
```

### Thread Configuration Guidelines

- **IO Threads**: Typically 1-2 threads per CPU core for network operations
- **Logic Threads**: Can be higher than CPU cores since game logic often waits on I/O
- **Login Server**: Lower thread counts (fewer concurrent operations)
- **Game Server**: Higher thread counts (more complex game logic)

## Server Architecture

### Threading Model

```
Client Connections
       ↓
   IO Thread Pool (Network I/O)
       ↓
   Task Dispatching (Routing)
       ↓
   Logic Thread Pool (Game Logic)
       ↓
   Lua Script Execution
```

### Routing Strategies

- **Login Server**: Socket File Descriptor (FD) based modulo operation for thread assignment
- **Game Server**: Character's current Map ID based modulo operation for thread assignment

### Dynamic Rerouting

The system supports dynamic task rerouting when routing keys change:
- Tasks are re-enqueued to the correct logic thread before execution
- Prevents infinite rerouting loops with attempt tracking
- Maintains system responsiveness during character map transitions

## Quick Start

### Prerequisites

- Go 1.19 or later
- Windows 10/11 (tested on Windows 10.0.26100)

### Building

```bash
# Build both servers
.\build.bat

# Or build individually
go build ./login
go build ./game
```

### Running

```bash
# Start login server
cd build
start-login.bat

# Start game server (in separate terminal)
start-game.bat
```

### Command Line Options

**Login Server:**
```bash
./login-server -host=0.0.0.0 -port=8484 -logic-threads=4 -stats
```

**Game Server:**
```bash
./game-server -host=0.0.0.0 -port=8485 -logic-threads=8 -world=Scania -stats
```

## Testing

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./common/...
go test ./core/...
go test ./entity/...
```

## Performance Monitoring

The server includes built-in metrics tracking:
- Task execution counts
- Rerouting statistics
- Error rates
- Active connection counts
- Lua script execution metrics

Enable statistics with the `-stats` flag when starting servers.

## Development

### Adding New Features

1. Follow the import cycle prevention rules
2. Add tests for new functionality
3. Update documentation
4. Verify build success

### Debugging

- Use `-stats` flag for real-time monitoring
- Check logs for error messages
- Verify thread pool configuration
- Monitor memory usage

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

1. Follow the import cycle prevention rules strictly
2. Add comprehensive tests
3. Update documentation
4. Ensure builds pass on Windows
5. Follow Go coding conventions 