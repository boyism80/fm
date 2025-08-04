package core

import (
	"fmt"
	"log"
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// LogicTask represents a task to be executed by a logic thread
type LogicTask struct {
	Predicate  func() bool       // First function: Check if logic can be executed
	Logic      func() error      // Second function: Execute the actual logic
	Callback   func(bool, error) // Third function: Report execution result
	Object     ThreadAssignable  // Object that this task belongs to (for thread validation)
	MaxRetries int               // Maximum number of retries for thread reassignment
	retryCount int               // Current retry count (internal use)
}

type LogicThread struct {
	id           int
	server       *Server
	stopChan     chan struct{}
	taskChan     chan *LogicTask // Channel for receiving logic tasks
	timerManager *TimerManager
	luaState     *lua.LState        // Pre-created Lua state for this thread
	luaMutex     sync.Mutex         // Mutex for Lua state access
	initFunc     func(*LogicThread) // Initialization function injected by server
}

// GetID returns the logic thread ID
func (t *LogicThread) GetID() int {
	return t.id
}

// GetLuaState returns the Lua state for this logic thread
func (t *LogicThread) GetLuaState() *lua.LState {
	t.luaMutex.Lock()
	defer t.luaMutex.Unlock()
	return t.luaState
}

// run executes the logic thread main loop
func (t *LogicThread) run() {
	log.Printf("Logic thread %d started", t.id)
	defer log.Printf("Logic thread %d stopped", t.id)

	// Initialize Lua state
	if t.initFunc != nil {
		t.initFunc(t)
	}

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
func (t *LogicThread) processTask(task *LogicTask) {
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
}

// SubmitTask submits a logic task to this logic thread for execution
func (t *LogicThread) SubmitTask(task *LogicTask) error {
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

// Schedule schedules a one-shot timer that will execute the given logic function
// after the specified duration. The logic will be executed on this logic thread.
func (t *LogicThread) Schedule(duration time.Duration, logic func() error, callback func(bool, error)) *Timer {
	timer := t.timerManager.Schedule(duration, logic, callback)

	// Start timer in a separate goroutine
	go t.runTimer(timer, duration)

	return timer
}

// runTimer runs the timer logic after the specified duration
func (t *LogicThread) runTimer(timer *Timer, duration time.Duration) {
	select {
	case <-time.After(duration):
		// Timer expired - create LogicTask and submit to this thread
		task := &LogicTask{
			Logic:    timer.Logic,
			Callback: timer.Callback,
		}
		t.taskChan <- task

	case <-timer.CancelChan:
		// Timer was cancelled
		return
	}

	// Clean up timer after execution or cancellation
	t.timerManager.RemoveTimer(timer.ID)
}

// GetTimerCount returns the number of active timers on this logic thread
func (t *LogicThread) GetTimerCount() int {
	return t.timerManager.GetTimerCount()
}

// SetRepeatingTimer sets a repeating timer that will execute the given logic function
// at the specified interval. The logic will be executed on this logic thread.
func (t *LogicThread) SetRepeatingTimer(interval time.Duration, logic func() error, callback func(bool, error)) *RepeatingTimer {
	repeatingTimer := t.timerManager.SetRepeatingTimer(interval, logic, callback)

	// Start repeating timer in a separate goroutine
	go t.runRepeatingTimer(repeatingTimer)

	return repeatingTimer
}

// runRepeatingTimer runs the repeating timer logic at the specified interval
func (t *LogicThread) runRepeatingTimer(timer *RepeatingTimer) {
	ticker := time.NewTicker(timer.Interval)
	defer ticker.Stop()

	// Execute immediately on start
	t.executeTimerLogic(timer)

	for {
		select {
		case <-ticker.C:
			if !timer.IsRunning {
				return
			}
			t.executeTimerLogic(timer)

		case <-timer.CancelChan:
			timer.IsRunning = false
			// Clean up timer
			t.timerManager.RemoveRepeatingTimer(timer.ID)
			return
		}
	}
}

// executeTimerLogic executes the timer logic as a LogicTask
func (t *LogicThread) executeTimerLogic(timer *RepeatingTimer) {
	task := &LogicTask{
		Logic:    timer.Logic,
		Callback: timer.Callback,
	}
	t.taskChan <- task
}

// GetRepeatingTimerCount returns the number of active repeating timers on this logic thread
func (t *LogicThread) GetRepeatingTimerCount() int {
	return t.timerManager.GetRepeatingTimerCount()
}

// CancelRepeatingTimer cancels a repeating timer by ID
func (t *LogicThread) CancelRepeatingTimer(timerID uint64) error {
	return t.timerManager.CancelRepeatingTimer(timerID)
}
