package memoryapp

import (
	"fmt"
	"os"
	"strings"
)

/*
   c0: zero register, non-writable and always zero
   c1: system code cell,(0 is run as normal, -1 is error, 1 is exit, 2 is write to STDOUT, 3 is read from STDIN)
   c2: writing scan/reading collection beginning,
   c3: writing scan/reading collection ending,
   otherwise the cells are freely allocatable
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

func swap(first *AllocEntry, second *AllocEntry) {
	temp := *first
	*first = *second
	*second = temp
}

func allocEntryBubbleSort(arr []AllocEntry) []AllocEntry {
	for i := 0; i < len(arr)-1; i++ {
		swapped := false
		for j := 0; j < len(arr)-1-i; j++ {
			if arr[j].start > arr[j+1].start {
				swap(&arr[j], &arr[j+1])
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}

	return arr
}

func (m Memory) PrintMemory(printAllocations bool) {
	fmt.Println(m.cells)
	fmt.Println(m.allocations)
	if printAllocations {
		var allocs []AllocEntry
		for _, a := range m.allocations {
			allocs = append(allocs, a)
		}
		allocs = allocEntryBubbleSort(allocs)

		var acc = []string{}
		for _, alloc := range allocs {
			var acc2 = []string{}
			for range alloc.Size() {
				acc2 = append(acc2, "*")
			}
			acc = append(acc, "["+strings.Join(acc2, " ")+"]")
		}
		fmt.Println(strings.Join(acc, "")[2:] + "\n\n")
	}
}

func (m Memory) validateCell(cell uint32, name string, isWriting bool) error {
	if cell > m.memSize {
		return fmt.Errorf("%s '%d' is outside the memory bounds", name, cell)
	} else if isWriting || name == "dst cell" {
		for allocName, alloc := range m.allocations {
			if alloc.IsCellInside(cell) {
				if alloc.isWritable {
					return nil
				}
				return fmt.Errorf("cannot not write to cell '%d' as that cell is in non-writable allocated area '%s'", cell, allocName)
			}
		}
		return fmt.Errorf("cannot not write to cell '%d' as that cell is not in a writable allocated area", cell)
	}

	return nil
}

func (m Memory) validateAllocation(name string, expectedType AllocEntryType) error {
	if _, ok := m.allocations[name]; !ok {
		return fmt.Errorf("allocation area '%s' does not exist", name)
	} else if m.allocations[name]._type != expectedType && expectedType != _ANY {
		return fmt.Errorf("expected type '%s' for allocation area '%s', but it's type is '%s' instead", expectedType, name, m.allocations[name]._type)
	}

	return nil
}

func (m Memory) getCell(cell uint32) int {
	return m.cells[cell]
}

func (m Memory) getLatestOpenCell() uint32 {
	var acc uint32 = 0
	for _, alloc := range m.allocations {
		acc += alloc.Size() + 1
	}
	return acc
}

func (m *Memory) forceWrite(cell uint32, value int) {
	m.cells[cell] = value
}

// allocates a range of memory
func (m *Memory) Allocate(name string, size uint32, _type AllocEntryType, isWritable, isLoadable bool) {
	var preoffset = m.getLatestOpenCell()
	if size > m.memSize {
		fmt.Println(fmt.Errorf("allocation area size %d is larger than the memory size %d", size, m.memSize))
		os.Exit(0)
	}
	if _, ok := m.allocations[name]; ok {
		fmt.Println(fmt.Errorf("allocation area %s already exists", name).Error())
		os.Exit(0)
	}
	m.allocations[name] = NewAllocEntry(preoffset, (preoffset+size)-1, _type, isWritable, isLoadable)
}

func (m *Memory) Deallocate(name string) error {
	if err := m.validateAllocation(name, _ANY); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	allocStart := m.allocations[name].start
	allocEnd := m.allocations[name].end
	delete(m.allocations, name)
	for c := allocStart; c < allocEnd+1; c++ {
		m.forceWrite(c, 0)
	}
	return nil
}

func (m Memory) Addi(dstCell, srcCell uint32, immediate int) error {
	if err := m.validateCell(dstCell, "dst cell", true); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	if err := m.validateCell(srcCell, "src cell 1", false); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	m.forceWrite(dstCell, m.getCell(srcCell)+immediate)
	return nil
}

func (m Memory) Pushi(dstAlloc string, immediate int) {
	if err := m.validateAllocation(dstAlloc, STACK); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	m.forceWrite(m.allocations[dstAlloc].start+m.allocations[dstAlloc].stackPointer, immediate)
	__temp := m.allocations[dstAlloc].Copy()
	__temp.stackPointer++
	m.allocations[dstAlloc] = __temp
}

func (m Memory) Pop(dstCell uint32, srcAlloc string) {
	if err := m.validateAllocation(srcAlloc, STACK); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	__temp := m.allocations[srcAlloc].Copy()
	__temp.stackPointer--
	m.allocations[srcAlloc] = __temp

	m.Mov(dstCell, m.allocations[srcAlloc].start+m.allocations[srcAlloc].stackPointer)

	m.forceWrite(m.allocations[srcAlloc].start+m.allocations[srcAlloc].stackPointer, 0)
}

func (m Memory) Add(dstCell, srcCell1, srcCell2 uint32) {
	if err := m.validateCell(srcCell2, "src cell 2", false); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	m.Addi(dstCell, srcCell1, m.getCell(srcCell2))
}

/*
Copies a value from cell `srcCell` to cell `dstCell`

Technically just an alias for `Memory.Addi([dstCell], [srcCell], 0)`
*/
func (m Memory) Mov(dstCell, srcCell uint32) {
	m.Addi(dstCell, srcCell, 0)
}

func NewMemory(memSize uint32) Memory {
	newMemory := Memory{}

	newMemory.memSize = memSize
	newMemory.cells = make([]int, memSize)
	newMemory.allocations = make(map[string]AllocEntry)

	newMemory.Allocate("SYSTEM_REGISTERS_NONWRITABLE", 1, NORMAL, false, false)
	newMemory.Allocate("SYSTEM_REGISTERS_WRITABLE", 3, NORMAL, true, false)

	return newMemory
}
