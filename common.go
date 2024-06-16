package memoryapp

// AllocEntry types
type AllocEntryType string

const (
	_ANY   AllocEntryType = "ANY"
	NORMAL AllocEntryType = "NORMAL"
	STACK  AllocEntryType = "STACK"  // FILO
	QUEUE  AllocEntryType = "QUEUE"  //
	DLLIST AllocEntryType = "DLLIST" // doubly linked list
)
