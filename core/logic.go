package core

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/common/types"
)

// LogicTask represents a task to be executed by a logic thread
// Flow: Predicate check -> Logic execution -> Result callback
// Components: Predicate function, logic function, result callback function
// Error Handling: Predicate failure, logic execution errors, callback failures
type LogicTask struct {
	Predicate  func() bool       // First function: Check if logic can be executed
	Logic      func() error      // Second function: Execute the actual logic
	Callback   func(bool, error) // Third function: Report execution result
	Object     ThreadAssignable  // Object that this task belongs to (for thread validation)
	MaxRetries int               // Maximum number of retries for thread reassignment
	retryCount int               // Current retry count (internal use)
}

// Thread Safety: Processes game state updates atomically
// Error Handling: Game logic errors, state inconsistencies
type LogicThread[T any] struct {
	id       int
	server   *Server[T]
	stopChan chan struct{}
	taskChan chan *LogicTask // Channel for receiving logic tasks
}

// run executes the logic thread main loop
// Flow: Receives tasks from channel -> Processes game logic -> Sends responses
// Game Logic: Character movement, combat, inventory, chat, etc.
// Thread Safety: Processes game state updates atomically
func (t *LogicThread[T]) run() {
	log.Printf("Logic thread %d started", t.id)
	defer log.Printf("Logic thread %d stopped", t.id)

	for {
		select {
		case <-t.stopChan:
			return
		case task := <-t.taskChan:
			t.processTask(task)
		}
	}
}

// processTask executes a logic task with the three-function pattern
// Flow: Thread validation -> Predicate check -> Logic execution -> Result callback
// Error Handling: Predicate failure, logic execution errors, callback failures, thread reassignment
// Thread Safety: Executes tasks atomically within the logic thread
func (t *LogicThread[T]) processTask(task *LogicTask) {
	if task == nil {
		log.Printf("Logic thread %d received nil task", t.id)
		return
	}

	// Validate thread assignment before execution
	if task.Object != nil {
		expectedThread, err := t.server.GetThreadForObject(task.Object)
		if err != nil {
			log.Printf("Logic thread %d: Failed to get thread for object: %v", t.id, err)
			if task.Callback != nil {
				task.Callback(false, err)
			}
			return
		}

		// Check if this task should be executed on a different thread
		if expectedThread != t.id {
			log.Printf("Logic thread %d: Task should be executed on thread %d, reassigning...", t.id, expectedThread)

			// Check retry limit
			if task.retryCount >= task.MaxRetries {
				log.Printf("Logic thread %d: Max retries (%d) exceeded for task reassignment", t.id, task.MaxRetries)
				if task.Callback != nil {
					task.Callback(false, fmt.Errorf("max retries exceeded for thread reassignment"))
				}
				return
			}

			// Increment retry count and reassign to correct thread
			task.retryCount++
			if err := t.server.SubmitLogicTaskToThread(expectedThread, task); err != nil {
				log.Printf("Logic thread %d: Failed to reassign task to thread %d: %v", t.id, expectedThread, err)
				if task.Callback != nil {
					task.Callback(false, fmt.Errorf("failed to reassign task: %w", err))
				}
			}
			return
		}
	}

	// First function: Check if logic can be executed
	canExecute := false
	if task.Predicate != nil {
		canExecute = task.Predicate()
	} else {
		canExecute = true // Default to true if no predicate provided
	}

	// Second function: Execute the logic if predicate allows
	var logicError error
	if canExecute && task.Logic != nil {
		logicError = task.Logic()
	} else if !canExecute {
		logicError = fmt.Errorf("predicate check failed")
	}

	// Third function: Report the result
	if task.Callback != nil {
		task.Callback(canExecute, logicError)
	}

	log.Printf("Logic thread %d processed task: predicate=%v, error=%v", t.id, canExecute, logicError)
}

// SubmitTask submits a logic task to this logic thread for execution
// Flow: Creates task -> Sends to task channel -> Logic thread processes
// Thread Safety: Safe to call from any thread, non-blocking
func (t *LogicThread[T]) SubmitTask(task *LogicTask) error {
	if task == nil {
		return fmt.Errorf("cannot submit nil task")
	}

	select {
	case t.taskChan <- task:
		return nil
	default:
		return fmt.Errorf("task channel is full")
	}
}

// processPacket processes a packet for a specific client
// Flow: Packet -> Client processing -> Response
// Thread Safety: Executes within logic thread context
func (t *LogicThread[T]) processPacket(client *Client[T], packet types.Packet) {
	// This is a placeholder for packet processing logic
	log.Printf("Logic thread %d processing packet for client %d", t.id, client.clientID)
}

// processPacketForClient processes a packet for a specific client with additional context
// Flow: Packet -> Client processing -> Response
// Thread Safety: Executes within logic thread context
func (t *LogicThread[T]) processPacketForClient(client *Client[T], packet types.Packet) {
	// This is a placeholder for packet processing logic with additional context
	log.Printf("Logic thread %d processing packet for client %d with additional context", t.id, client.clientID)
}
