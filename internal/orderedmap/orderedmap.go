package orderedmap

import (
	"container/list"
)

// Element can be used to get next
type Element struct {
	*keyValue
	elem *list.Element // needed to able delete with O(1)
}

// Next returns next element to current one
func (el *Element) Next() *Element {
	// This is wrapper against linked list to avoid type casting in client code

	e := el.elem.Next()
	if e == nil {
		return nil
	}

	return &Element{
		elem:     e,
		keyValue: e.Value.(*keyValue),
	}
}

type keyValue struct {
	Key, Value interface{}
}

// OrderedMap is ordered map structure
// It is using power of linked list and hash table
type OrderedMap struct {
	ll   *list.List
	data map[interface{}]*list.Element
}

// New creates ordered map
func New() *OrderedMap {
	return &OrderedMap{
		ll:   list.New(),
		data: make(map[interface{}]*list.Element),
	}
}

// Load returns the value stored in the map for a key, or nil if no value is present. The ok result indicates whether value was found in the map.
func (m *OrderedMap) Load(key interface{}) (value interface{}, ok bool) {
	elem, ok := m.data[key]
	if !ok {
		return nil, false
	}

	return elem.Value.(*keyValue).Value, true
}

// Store sets the value for a key.
func (m *OrderedMap) Store(key, value interface{}) {
	_, ok := m.data[key]
	if !ok {
		el := m.ll.PushBack(&keyValue{Key: key, Value: value})
		m.data[key] = el
	} else {
		m.data[key].Value.(*keyValue).Value = value
	}
}

// Delete deletes the value for a key.
func (m *OrderedMap) Delete(key interface{}) {
	if elem, ok := m.data[key]; ok {
		delete(m.data, key)
		m.ll.Remove(elem) // O(1)
	}
}

// Front returns the first element of ordered map or nil if the ordered map is empty.
func (m *OrderedMap) Front() *Element {
	front := m.ll.Front()
	if front == nil {
		return nil
	}

	return &Element{
		elem:     front,
		keyValue: front.Value.(*keyValue),
	}
}

// Back returns the last element of ordered map l or nil if the ordered map is empty.
func (m *OrderedMap) Back() *Element {
	back := m.ll.Back()
	if back == nil {
		return nil
	}

	return &Element{
		elem:     back,
		keyValue: back.Value.(*keyValue),
	}
}
