package storage

import (
	"unsafe"

	"github.com/yohamta/donburi/internal/component"
)

// SimpleStorage is a structure that stores the pointer to data of each component.
// It stores the pointers in the two dimensional slice.
// First dimension is the archetype index.
// Second dimension is the component index.
// The component index is used to access the component data in the archetype.
type SimpleStorage struct {
	storages [][]unsafe.Pointer
}

// NewSimpleStorage creates a new empty structure that stores the pointer to data of each component.
func NewSimpleStorage() *SimpleStorage {
	return &SimpleStorage{
		storages: make([][]unsafe.Pointer, 256),
	}
}

// PushComponent stores the new data of the component in the archetype.
func (cs *SimpleStorage) PushComponent(component *component.ComponentType, archetypeIndex ArchetypeIndex) {
	if v := cs.storages[archetypeIndex]; v == nil {
		cs.storages[archetypeIndex] = []unsafe.Pointer{}
	}
	// TODO: optimize to avoid allocation
	componentValue := component.New()
	cs.storages[archetypeIndex] = append(cs.storages[archetypeIndex], componentValue)
}

// Component returns the pointer to data of the component in the archetype.
func (cs *SimpleStorage) Component(archetypeIndex ArchetypeIndex, componentIndex ComponentIndex) unsafe.Pointer {
	if int(archetypeIndex) >= len(cs.storages) {
		panic("archetypeIndex out of range")
	}
	if int(componentIndex) >= len(cs.storages[archetypeIndex]) {
		panic("componentIndex out of range")
	}
	return cs.storages[archetypeIndex][componentIndex]
}

// SetComponent sets the pointer to data of the component in the archetype.
func (cs *SimpleStorage) SetComponent(archetypeIndex ArchetypeIndex, componentIndex ComponentIndex, component unsafe.Pointer) {
	cs.storages[archetypeIndex][componentIndex] = component
}

// MoveComponent moves the pointer to data of the component in the archetype.
func (cs *SimpleStorage) MoveComponent(source ArchetypeIndex, index ComponentIndex, dst ArchetypeIndex) {
	src_slice := cs.storages[source]
	dst_slice := cs.storages[dst]

	value := src_slice[index]
	src_slice[index] = src_slice[len(src_slice)-1]
	src_slice = src_slice[:len(src_slice)-1]
	cs.storages[source] = src_slice

	dst_slice = append(dst_slice, value)
	cs.storages[dst] = dst_slice
}

// SwapRemove removes the pointer to data of the component in the archetype.
func (cs *SimpleStorage) SwapRemove(archetypeIndex ArchetypeIndex, componentIndex ComponentIndex) unsafe.Pointer {
	componentValue := cs.storages[archetypeIndex][componentIndex]
	cs.storages[archetypeIndex][componentIndex] = cs.storages[archetypeIndex][len(cs.storages[archetypeIndex])-1]
	cs.storages[archetypeIndex] = cs.storages[archetypeIndex][:len(cs.storages[archetypeIndex])-1]
	return componentValue
}
