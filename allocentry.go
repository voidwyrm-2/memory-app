package memoryapp

type AllocEntry struct {
	start uint32
	end   uint32
	_type AllocEntryType

	isWritable   bool
	isLoadable   bool
	stackPointer uint32
}

func (ae AllocEntry) Size() uint32 {
	return ae.end - ae.start
}

func (ae AllocEntry) IsCellInside(cell uint32) bool {
	if cell >= ae.start && cell <= ae.end {
		return true
	}
	return false
}

func (ae AllocEntry) Copy() AllocEntry {
	newAllocEntry := AllocEntry{}

	newAllocEntry.start = ae.start
	newAllocEntry.end = ae.end
	newAllocEntry._type = ae._type
	newAllocEntry.isWritable = ae.isWritable
	newAllocEntry.isLoadable = ae.isLoadable
	newAllocEntry.stackPointer = ae.stackPointer

	return newAllocEntry
}

func NewAllocEntry(start, end uint32, _type AllocEntryType, isWritable, isLoadable bool) AllocEntry {
	newAllocEntry := AllocEntry{}

	newAllocEntry.start = start
	newAllocEntry.end = end
	newAllocEntry._type = _type
	newAllocEntry.isWritable = isWritable
	newAllocEntry.isLoadable = isLoadable
	newAllocEntry.stackPointer = 0

	return newAllocEntry
}
