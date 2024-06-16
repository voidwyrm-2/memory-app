package main

import memoryapp "github.com/voidwyrm-2/memory-app"

func errortest() {
	memory := memoryapp.NewMemory(64)

	memory.Addi(4, 0, 20)
}
