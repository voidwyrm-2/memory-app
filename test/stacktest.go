package main

import (
	memoryapp "github.com/voidwyrm-2/memory-app"
)

func stacktest() {
	var memory = memoryapp.NewMemory(64)

	// allocate a new area of memory to store some extra values we need
	memory.Allocate("misc", 9, memoryapp.NORMAL, true, true)

	// allocate a new area of memory as a stack, with a size of 10
	memory.Allocate("stack", 10, memoryapp.STACK, false, true)

	memory.PrintMemory(true)

	// push '21' to the stack
	memory.Pushi("stack", 21)

	memory.PrintMemory(true)

	// push '21' to the stack again
	memory.Pushi("stack", 21)

	memory.PrintMemory(true)

	// pop from the stack and store the popped value in the fifth memory cell
	memory.Pop(4, "stack")

	memory.PrintMemory(true)

	// pop from the stack and store the popped value in the seventh memory cell
	memory.Pop(6, "stack")

	memory.PrintMemory(true)

	// add the values of the fifth and seventh memory cells together and store the result in the ninth memory cell
	memory.Add(8, 6, 4)

	memory.PrintMemory(true)

	// say we have a bunch of values in an allocated area
	// in this case, we have a stack with '43', '123', '456', and '989'
	memory.Pushi("stack", 43)

	memory.PrintMemory(true)

	memory.Pushi("stack", 123)

	memory.PrintMemory(true)

	memory.Pushi("stack", 456)

	memory.PrintMemory(true)

	memory.Pushi("stack", 989)

	memory.PrintMemory(true)

	// but we aren't using those values anymore
	// and since that allocated area is there, we can't use it for other things
	// so we need to deallocate it
	memory.Deallocate("stack")

	memory.PrintMemory(true)
}
