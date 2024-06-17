package memoryplayground

import (
	"fmt"
	"math/rand"
	"os"
)

// AllocEntry types

/*
can be
_ANY,
NORMAL,
STACK,
QUEUE,
DLLIST
*/
type AllocEntryType string

const (
	_ANY   AllocEntryType = "ANY"    // Nonspecific AllocEntryType. if used in Alloc(), it's treated as NORMAL
	NORMAL AllocEntryType = "NORMAL" // Normal area of non-expandable writable/unwritable memory, and can be used as an array if needed.
	STACK  AllocEntryType = "STACK"  // Stack; First-In Last-Out; non-expandable.
	QUEUE  AllocEntryType = "QUEUE"  // Queue; Last-In First-Out; non-expandable.
	DLLIST AllocEntryType = "DLLIST" // Doubly Linked List; expandable.
)

var AllocEntryTypes = []AllocEntryType{
	NORMAL,
	STACK,
	QUEUE,
	DLLIST,
}

// functions

// alias for os.Exit(0)
func exit() {
	os.Exit(0)
}

func printErr(format string, a ...any) {
	fmt.Println(fmt.Errorf(format, a...).Error())
	exit()
}

func generateRandomName(length int) string {
	out := ""
	for range length {
		out += fmt.Sprint(rand.Intn(10))
	}
	return out
}
