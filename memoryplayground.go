package memoryplayground

import (
	"fmt"
	"strings"
)

/*
	c0: zero register, non-writable and always zero

	c1: system code cell,(0 is run as normal, -1 is error, 1 is exit, 2 is write to STDOUT, 3 is read from STDIN)

	c2: offset of bytes to read/write from/to when writing to STDOUT/reading from STDIN

	c3: amount of bytes to read/write from/to when writing to STDOUT/reading from STDIN

	c4-c12: the global stack

	otherwise the cells are freely allocatable
*/

/*
The struct this package revolves around
*/
type Memory struct {
	memSize     uint32
	cells       []int
	allocations map[string]AllocEntry
}

/*
func (m Memory) formatAllocations() string {
	out := ""
	revTick := len(m.allocations) - 1
	for allocName, alloc := range m.allocations {

	}
	return out
}
*/

/*
Prints the struct's memory cells. If print_allocations is true, it also prints the struct's allocations

Memory cells that are -2 are unallocated, memorycells that are -1 are allocated but haven't been assigned yet,
and otherwise they are allocated and assigned
*/
func (m Memory) Print(print_allocations bool) {
	if print_allocations {
		fmt.Printf("%v\n\n", m.cells)
		fmt.Printf("%v\n\n\n", m.allocations)
	} else {
		fmt.Printf("%v\n\n\n", m.cells)
	}
}

func (m Memory) validateAllocation(name string, expectedType AllocEntryType) {
	if _, ok := m.allocations[name]; !ok {
		printErr("allocation area '%s' does not exist", name)
	} else if m.allocations[name]._type != expectedType && expectedType != _ANY {
		printErr("expected type '%s' for allocation area '%s', but it's type is '%s' instead", expectedType, name, m.allocations[name]._type)
	}
}

/*
 */
func (m Memory) getTypedAlloc(name string, _type AllocEntryType) AllocEntry {
	m.validateAllocation(name, _type)
	return m.allocations[name]
}

func (m Memory) getAlloc(name string) AllocEntry {
	return m.getTypedAlloc(name, _ANY)
}

/*
Validates that the given cell exists and is writable

cell is the global index of the cell
*/
func (m Memory) validateGlobalCell(cellName string, cell uint32) {
	if cellName == "" {
		cellName = "cell"
	}

	if cell > m.memSize {
		printErr("%s '%d' is outside the memory bounds", cellName, cell)
	}
}

/*
Validates that the given cell exists and is writable

cell is relative to the allocation start index

Add '!!write' onto the end of cellName to tell the method the cell is being written to
*/
func (m Memory) validateCell(allocName, cellName string, cell uint32) {
	isWriting := false
	if strings.ToLower(cellName[len(cellName)-7:]) == "!!write" {
		cellName = cellName[:len(cellName)-7]
		isWriting = true
	} else if cellName == "dst cell" {
		isWriting = true
	}

	if cellName == "GLOBAL_STACK" {
		printErr("the global stack is not accessible in any way by the user")
	} else if cellName == "" {
		cellName = "cell"
	}

	if isWriting {
		if !m.getAlloc(allocName).isWritable {
			printErr("cannot not write to %s '%d' as that cell is in non-writable allocated area '%s'", cellName, cell, allocName)
		}
	}
	m.validateGlobalCell(cellName, cell)
	m.validateGlobalCell(cellName, cell+m.getAlloc(allocName).start)
}

func (m Memory) findLatestOpenCellsOfSize(size uint32) uint32 {
	if size == 0 {
		panic("findLatestOpenCellsOfSize paramater 'size' cannot be 0")
	}

	var c uint32 = 0
	for c < m.memSize-1 {
		isInside, allocName := m.cellIsInsideAnyAllocs(c)
		if isInside {
			// fmt.Println(c)
			c += m.getAlloc(allocName).Size() + 1
			// fmt.Println(c)
			// fmt.Print("\n")
		} else {
			var count uint32 = 0
			var orig uint32 = c
			//fmt.Printf("%d(%d), %d\n", c, m.cells[c], count)
			for m.forceRead(c) == -2 && c < m.memSize-1 {
				// fmt.Printf("c: %d\n", c)
				if count == 0 {
					count++
				} else {
					count++
				}
				if count >= size {
					return orig
				}
				c++
			}
			//fmt.Println("HUH???")
			count = 0
		}
	}

	return 0
}

func (m Memory) cellIsInsideAnyAllocs(cell uint32) (bool, string) {
	for allocName, alloc := range m.allocations {
		if alloc.IsCellInside(cell) {
			return true, allocName
		}
	}
	return false, ""
}

func (m Memory) forceRead(cell uint32) int {
	return m.cells[cell]
}

func (m *Memory) forceWrite(cell uint32, value int) {
	m.cells[cell] = value
}

func (m Memory) forceWriteArea(start, end uint32, value int) {
	for c := start; c < end+1; c++ {
		m.forceWrite(c, value)
	}
}

func (m Memory) read(allocName, cellName string, cell uint32) int {
	m.validateCell(allocName, cellName, cell)
	/*
		switch m.getAlloc(allocName)._type {
		case STACK:
			return m.forceRead(m.getAlloc(allocName).start + m.getAlloc(allocName).stackPointer)
		default:*/
	return m.forceRead(m.getAlloc(allocName).start + cell)
	//}

}

func (m Memory) write(allocName, cellName string, cell uint32, value int) {
	if !strings.HasSuffix(cellName, "!!write") && cellName != "dst cell" {
		cellName += "!!write"
	}

	m.validateCell(allocName, cellName, cell)
	m.forceWrite(m.getAlloc(allocName).start+cell, value)
}

func (m *Memory) Allocate(name string, start, size uint32, _type AllocEntryType, startIsRelative, isWritable, isLoadable bool) {
	var preoffset uint32 = 0
	if startIsRelative {
		preoffset = m.findLatestOpenCellsOfSize(size)
	}

	if _type == _ANY {
		_type = NORMAL
	}

	if size > m.memSize {
		printErr("allocation area size %d is larger than the memory size %d", size, m.memSize)
	}

	m.allocations[name] = NewAllocEntry(preoffset, (preoffset+size)-1, _type, isWritable, isLoadable)
	m.forceWriteArea(m.allocations[name].start, m.allocations[name].end, -1)
}

func (m *Memory) Deallocate(name string) {
	if strings.HasPrefix(name, "SYSTEM_") {
		printErr("allocated area '%s' cannot be deallocated because it is a system allocated area", name)
	} else if name == "GLOBAL_STACK" {
		printErr("the global stack cannot be deallocated")
	}
	m.validateAllocation(name, _ANY)

	start := m.getAlloc(name).start
	end := m.getAlloc(name).end

	delete(m.allocations, name)

	m.forceWriteArea(start, end, -2)
}

func (m *Memory) Pushi(dstAllocName string, immediate int) {
	stack := m.getTypedAlloc(dstAllocName, STACK)

	m.forceWrite(stack.start+stack.stackPointer, immediate)
	stack.stackPointer++
	m.allocations[dstAllocName] = stack
}

func (m Memory) Push(dstAllocName, srcAllocName string, srcCell uint32) {
	m.Pushi(dstAllocName, m.read(srcAllocName, "src cell", srcCell))
}

func (m Memory) Mov(dstAllocName, srcAllocName string, dstCell, srcCell uint32) {
	m.write(dstAllocName, "dst cell", dstCell, m.read(srcAllocName, "src cell", srcCell))
}

const _SYSTEM_REGISTERS_NONWRITABLE_SIZE = 1
const _SYSTEM_REGISTERS_WRITABLE_SIZE = 4
const _GLOBAL_STACK_SIZE = 8

func NewMemory(memSize uint32) Memory {
	memSize += _SYSTEM_REGISTERS_NONWRITABLE_SIZE + _SYSTEM_REGISTERS_WRITABLE_SIZE + _GLOBAL_STACK_SIZE
	newMemory := Memory{}

	newMemory.memSize = memSize
	newMemory.cells = make([]int, memSize)
	for i := range newMemory.cells {
		// if i >= 0 && i <= _SYSTEM_REGISTERS_WRITABLE_SIZE {
		// 	continue
		// }
		newMemory.cells[i] = -2
	}
	newMemory.allocations = make(map[string]AllocEntry)

	newMemory.Allocate("SYSTEM_REGISTERS_NONWRITABLE", 0, _SYSTEM_REGISTERS_NONWRITABLE_SIZE, NORMAL, false, false, false)
	newMemory.Allocate("SYSTEM_REGISTERS_WRITABLE", 0, _SYSTEM_REGISTERS_WRITABLE_SIZE, NORMAL, true, true, false)
	newMemory.Allocate("GLOBAL_STACK", 0, _GLOBAL_STACK_SIZE, STACK, true, false, false)
	newMemory.forceWriteArea(0, (_SYSTEM_REGISTERS_NONWRITABLE_SIZE+_SYSTEM_REGISTERS_WRITABLE_SIZE+_GLOBAL_STACK_SIZE)-1, 0)

	return newMemory
}
