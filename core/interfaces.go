package core

// ThreadAssignable represents objects that can be assigned to specific logic threads
// Flow: Object creation -> Thread assignment -> Logic processing
// Purpose: Abstract objects that need to run on specific logic threads
// Error Conditions: Invalid thread assignment, object state errors
type ThreadAssignable interface {
	// GetThreadHash returns a hash value used to determine which logic thread to use
	// Flow: Object identification -> Hash generation -> Thread assignment
	// Thread Assignment: Uses hash % logicThreadCount to determine target thread
	// Error Handling: Returns consistent hash for same object state
	GetThreadHash() int
}
