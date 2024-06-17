package main

import mp "github.com/voidwyrm-2/memory-playground"

func firsttest() {
	memory := mp.NewMemory(64)

	memory.Print(true)

	memory.Allocate("teststack", 0, 5, mp.STACK, true, false, false)

	memory.Print(true)

	memory.Pushi("teststack", 42)

	memory.Print(true)

	memory.Pushi("teststack", 60)

	memory.Print(true)

	memory.Pushi("teststack", 32)

	memory.Print(true)

	memory.Deallocate("teststack")

	memory.Print(true)
}
