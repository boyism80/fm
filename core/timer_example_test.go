package core

import (
	"fmt"
	"log"
	"net"
	"time"
)

// ExampleTimerUsage demonstrates how to use the timer system
func ExampleTimerUsage() {
	// Create server config
	config := &ServerConfig{
		LogicThreadCount: 4,
		ClientFactory: func(conn net.Conn, clientID int) (Client, error) {
			// Mock client factory for example
			return nil, nil
		},
	}

	// Create and start server
	server, err := NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start("localhost", 8080)
	if err != nil {
		log.Fatal(err)
	}

	// Create a mock ThreadAssignable object
	mockObject := &MockThreadAssignable{hash: 123}

	// Get the logic thread for this object
	logicThread, err := server.GetLogicThread(mockObject)
	if err != nil {
		log.Fatal(err)
	}

	// Schedule a one-shot timer
	timer := logicThread.Schedule(
		time.Second*5,
		func() error {
			fmt.Println("Timer executed after 5 seconds!")
			return nil
		},
		func(success bool, err error) {
			if err != nil {
				fmt.Printf("Timer callback error: %v\n", err)
			} else {
				fmt.Println("Timer executed successfully")
			}
		},
	)

	// You can cancel the timer if needed
	// timer.Cancel()

	// Wait a bit to see the timer execute
	time.Sleep(time.Second * 6)

	// Get server stats
	stats := server.GetStats()
	fmt.Printf("Server stats: %+v\n", stats)

	// Stop the server
	server.Stop()
}

// ExampleRepeatingTimerUsage demonstrates how to use the repeating timer system
func ExampleRepeatingTimerUsage() {
	// Create server config
	config := &ServerConfig{
		LogicThreadCount: 4,
		ClientFactory: func(conn net.Conn, clientID int) (Client, error) {
			// Mock client factory for example
			return nil, nil
		},
	}

	// Create and start server
	server, err := NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start("localhost", 8080)
	if err != nil {
		log.Fatal(err)
	}

	// Set a repeating timer on all logic threads
	timers := server.SetTimer(
		time.Second*2,
		func() error {
			fmt.Println("Repeating timer executed every 2 seconds!")
			return nil
		},
		func(success bool, err error) {
			if err != nil {
				fmt.Printf("Repeating timer callback error: %v\n", err)
			} else {
				fmt.Println("Repeating timer executed successfully")
			}
		},
	)

	fmt.Printf("Set %d repeating timers across all logic threads\n", len(timers))

	// Wait for a few executions
	time.Sleep(time.Second * 10)

	// Get server stats
	stats := server.GetStats()
	fmt.Printf("Server stats: %+v\n", stats)

	// Cancel all repeating timers
	for _, timer := range timers {
		timer.Cancel()
	}

	// Stop the server
	server.Stop()
}

// ExampleTimerManagerUsage demonstrates how to use TimerManager directly
func ExampleTimerManagerUsage() {
	// Create a timer manager
	timerManager := NewTimerManager()

	// Schedule a one-shot timer
	timer := timerManager.Schedule(
		time.Second*3,
		func() error {
			fmt.Println("Timer executed!")
			return nil
		},
		func(success bool, err error) {
			if err != nil {
				fmt.Printf("Timer callback error: %v\n", err)
			} else {
				fmt.Println("Timer executed successfully")
			}
		},
	)

	// Set a repeating timer
	repeatingTimer := timerManager.SetRepeatingTimer(
		time.Second*1,
		func() error {
			fmt.Println("Repeating timer executed!")
			return nil
		},
		func(success bool, err error) {
			if err != nil {
				fmt.Printf("Repeating timer callback error: %v\n", err)
			}
		},
	)

	// Wait for some executions
	time.Sleep(time.Second * 5)

	// Cancel timers
	timer.Cancel()
	repeatingTimer.Cancel()

	fmt.Printf("Timer count: %d\n", timerManager.GetTimerCount())
	fmt.Printf("Repeating timer count: %d\n", timerManager.GetRepeatingTimerCount())
}

// MockThreadAssignable is a mock implementation for testing
type MockThreadAssignable struct {
	hash int
}

func (m *MockThreadAssignable) GetThreadHash() int {
	return m.hash
}
