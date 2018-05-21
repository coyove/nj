// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package potatolang

import (
	"fmt"
	"testing"
)

func TestRedBlackTreePut(t *testing.T) {
	tree := NewMap()
	tree.Put("5", NewStringValue("e"))
	tree.Put("6", NewStringValue("f"))
	tree.Put("7", NewStringValue("g"))
	tree.Put("3", NewStringValue("c"))
	tree.Put("4", NewStringValue("d"))
	tree.Put("1", NewStringValue("x"))
	tree.Put("2", NewStringValue("b"))
	tree.Put("1", NewStringValue("a")) //overwrite

	if actualValue := tree.Size(); actualValue != 7 {
		t.Errorf("Got %v expected %v", actualValue, 7)
	}

	tests1 := [][]interface{}{
		{"1", NewStringValue("a"), true},
		{"2", NewStringValue("b"), true},
		{"3", NewStringValue("c"), true},
		{"4", NewStringValue("d"), true},
		{"5", NewStringValue("e"), true},
		{"6", NewStringValue("f"), true},
		{"7", NewStringValue("g"), true},
		{"8", NewValue(), false},
	}

	for _, test := range tests1 {
		// retrievals
		actualValue, actualFound := tree.Get(test[0].(string))
		if !actualValue.Equal(test[1].(Value)) || actualFound != test[2] {
			t.Errorf("Got %v expected %v", actualValue, test[1])
		}
	}
}

func TestRedBlackTreeRemove(t *testing.T) {
	tree := NewMap()
	tree.Put("5", NewStringValue("e"))
	tree.Put("6", NewStringValue("f"))
	tree.Put("7", NewStringValue("g"))
	tree.Put("3", NewStringValue("c"))
	tree.Put("4", NewStringValue("d"))
	tree.Put("1", NewStringValue("x"))
	tree.Put("2", NewStringValue("b"))
	tree.Put("1", NewStringValue("a")) //overwrite

	tree.Remove("5")
	tree.Remove("6")
	tree.Remove("7")
	tree.Remove("8")
	tree.Remove("5")

	if actualValue := tree.Size(); actualValue != 4 {
		t.Errorf("Got %v expected %v", actualValue, 7)
	}

	tests2 := [][]interface{}{
		{"1", NewStringValue("a"), true},
		{"2", NewStringValue("b"), true},
		{"3", NewStringValue("c"), true},
		{"4", NewStringValue("d"), true},
		{"5", NewValue(), false},
		{"6", NewValue(), false},
		{"7", NewValue(), false},
		{"8", NewValue(), false},
	}

	for _, test := range tests2 {
		actualValue, actualFound := tree.Get(test[0].(string))
		if !actualValue.Equal(test[1].(Value)) || actualFound != test[2] {
			t.Errorf("Got %v expected %v", actualValue, test[1])
		}
	}

	tree.Remove("1")
	tree.Remove("4")
	tree.Remove("2")
	tree.Remove("3")
	tree.Remove("2")
	tree.Remove("2")

	if empty, size := tree.Empty(), tree.Size(); empty != true || size != -0 {
		t.Errorf("Got %v expected %v", empty, true)
	}

}

func TestRedBlackTreeLeftAndRight(t *testing.T) {
	tree := NewMap()

	if actualValue := tree.Left(); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	if actualValue := tree.Right(); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	tree.Put("1", NewStringValue("a"))
	tree.Put("5", NewStringValue("e"))
	tree.Put("6", NewStringValue("f"))
	tree.Put("7", NewStringValue("g"))
	tree.Put("3", NewStringValue("c"))
	tree.Put("4", NewStringValue("d"))
	tree.Put("1", NewStringValue("x")) // overwrite
	tree.Put("2", NewStringValue("b"))

	if actualValue, expectedValue := fmt.Sprintf("%s", tree.Left().Key), "1"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := fmt.Sprintf("%v", tree.Left().Value), `"x"`; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := fmt.Sprintf("%s", tree.Right().Key), "7"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := fmt.Sprintf("%v", tree.Right().Value), `"g"`; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

// func TestRedBlackTreeCeilingAndFloor(t *testing.T) {
// 	tree := NewWithIntComparator()

// 	if node, found := tree.Floor(0); node != nil || found {
// 		t.Errorf("Got %v expected %v", node, "<nil>")
// 	}
// 	if node, found := tree.Ceiling(0); node != nil || found {
// 		t.Errorf("Got %v expected %v", node, "<nil>")
// 	}

// 	tree.Put(5, "e")
// 	tree.Put(6, "f")
// 	tree.Put(7, "g")
// 	tree.Put(3, "c")
// 	tree.Put(4, "d")
// 	tree.Put(1, "x")
// 	tree.Put(2, "b")

// 	if node, found := tree.Floor(4); node.Key != 4 || !found {
// 		t.Errorf("Got %v expected %v", node.Key, 4)
// 	}
// 	if node, found := tree.Floor(0); node != nil || found {
// 		t.Errorf("Got %v expected %v", node, "<nil>")
// 	}

// 	if node, found := tree.Ceiling(4); node.Key != 4 || !found {
// 		t.Errorf("Got %v expected %v", node.Key, 4)
// 	}
// 	if node, found := tree.Ceiling(8); node != nil || found {
// 		t.Errorf("Got %v expected %v", node, "<nil>")
// 	}
// }

// func TestRedBlackTreeIteratorNextOnEmpty(t *testing.T) {
// 	tree := NewWithIntComparator()
// 	it := tree.Iterator()
// 	for it.Next() {
// 		t.Errorf("Shouldn't iterate on empty tree")
// 	}
// }

// func TestRedBlackTreeIteratorPrevOnEmpty(t *testing.T) {
// 	tree := NewWithIntComparator()
// 	it := tree.Iterator()
// 	for it.Prev() {
// 		t.Errorf("Shouldn't iterate on empty tree")
// 	}
// }

// func TestRedBlackTreeIterator1Next(t *testing.T) {
// 	tree := NewWithIntComparator()
// 	tree.Put(5, "e")
// 	tree.Put(6, "f")
// 	tree.Put(7, "g")
// 	tree.Put(3, "c")
// 	tree.Put(4, "d")
// 	tree.Put(1, "x")
// 	tree.Put(2, "b")
// 	tree.Put(1, "a") //overwrite
// 	// │   ┌── 7
// 	// └── 6
// 	//     │   ┌── 5
// 	//     └── 4
// 	//         │   ┌── 3
// 	//         └── 2
// 	//             └── 1
// 	it := tree.Iterator()
// 	count := 0
// 	for it.Next() {
// 		count++
// 		key := it.Key()
// 		switch key {
// 		case count:
// 			if actualValue, expectedValue := key, count; actualValue != expectedValue {
// 				t.Errorf("Got %v expected %v", actualValue, expectedValue)
// 			}
// 		default:
// 			if actualValue, expectedValue := key, count; actualValue != expectedValue {
// 				t.Errorf("Got %v expected %v", actualValue, expectedValue)
// 			}
// 		}
// 	}
// 	if actualValue, expectedValue := count, tree.Size(); actualValue != expectedValue {
// 		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
// 	}
// }

// func TestRedBlackTreeIterator1Prev(t *testing.T) {
// 	tree := NewWithIntComparator()
// 	tree.Put(5, "e")
// 	tree.Put(6, "f")
// 	tree.Put(7, "g")
// 	tree.Put(3, "c")
// 	tree.Put(4, "d")
// 	tree.Put(1, "x")
// 	tree.Put(2, "b")
// 	tree.Put(1, "a") //overwrite
// 	// │   ┌── 7
// 	// └── 6
// 	//     │   ┌── 5
// 	//     └── 4
// 	//         │   ┌── 3
// 	//         └── 2
// 	//             └── 1
// 	it := tree.Iterator()
// 	for it.Next() {
// 	}
// 	countDown := tree.size
// 	for it.Prev() {
// 		key := it.Key()
// 		switch key {
// 		case countDown:
// 			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
// 				t.Errorf("Got %v expected %v", actualValue, expectedValue)
// 			}
// 		default:
// 			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
// 				t.Errorf("Got %v expected %v", actualValue, expectedValue)
// 			}
// 		}
// 		countDown--
// 	}
// 	if actualValue, expectedValue := countDown, 0; actualValue != expectedValue {
// 		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
// 	}
// }

// func TestRedBlackTreeIterator2Next(t *testing.T) {
// 	tree := NewWithIntComparator()
// 	tree.Put(3, "c")
// 	tree.Put(1, "a")
// 	tree.Put(2, "b")
// 	it := tree.Iterator()
// 	count := 0
// 	for it.Next() {
// 		count++
// 		key := it.Key()
// 		switch key {
// 		case count:
// 			if actualValue, expectedValue := key, count; actualValue != expectedValue {
// 				t.Errorf("Got %v expected %v", actualValue, expectedValue)
// 			}
// 		default:
// 			if actualValue, expectedValue := key, count; actualValue != expectedValue {
// 				t.Errorf("Got %v expected %v", actualValue, expectedValue)
// 			}
// 		}
// 	}
// 	if actualValue, expectedValue := count, tree.Size(); actualValue != expectedValue {
// 		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
// 	}
// }

// func TestRedBlackTreeIterator2Prev(t *testing.T) {
// 	tree := NewWithIntComparator()
// 	tree.Put(3, "c")
// 	tree.Put(1, "a")
// 	tree.Put(2, "b")
// 	it := tree.Iterator()
// 	for it.Next() {
// 	}
// 	countDown := tree.size
// 	for it.Prev() {
// 		key := it.Key()
// 		switch key {
// 		case countDown:
// 			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
// 				t.Errorf("Got %v expected %v", actualValue, expectedValue)
// 			}
// 		default:
// 			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
// 				t.Errorf("Got %v expected %v", actualValue, expectedValue)
// 			}
// 		}
// 		countDown--
// 	}
// 	if actualValue, expectedValue := countDown, 0; actualValue != expectedValue {
// 		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
// 	}
// }

// func TestRedBlackTreeIterator3Next(t *testing.T) {
// 	tree := NewWithIntComparator()
// 	tree.Put(1, "a")
// 	it := tree.Iterator()
// 	count := 0
// 	for it.Next() {
// 		count++
// 		key := it.Key()
// 		switch key {
// 		case count:
// 			if actualValue, expectedValue := key, count; actualValue != expectedValue {
// 				t.Errorf("Got %v expected %v", actualValue, expectedValue)
// 			}
// 		default:
// 			if actualValue, expectedValue := key, count; actualValue != expectedValue {
// 				t.Errorf("Got %v expected %v", actualValue, expectedValue)
// 			}
// 		}
// 	}
// 	if actualValue, expectedValue := count, tree.Size(); actualValue != expectedValue {
// 		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
// 	}
// }

// func TestRedBlackTreeIterator3Prev(t *testing.T) {
// 	tree := NewWithIntComparator()
// 	tree.Put(1, "a")
// 	it := tree.Iterator()
// 	for it.Next() {
// 	}
// 	countDown := tree.size
// 	for it.Prev() {
// 		key := it.Key()
// 		switch key {
// 		case countDown:
// 			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
// 				t.Errorf("Got %v expected %v", actualValue, expectedValue)
// 			}
// 		default:
// 			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
// 				t.Errorf("Got %v expected %v", actualValue, expectedValue)
// 			}
// 		}
// 		countDown--
// 	}
// 	if actualValue, expectedValue := countDown, 0; actualValue != expectedValue {
// 		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
// 	}
// }

// func TestRedBlackTreeIterator4Next(t *testing.T) {
// 	tree := NewWithIntComparator()
// 	tree.Put(13, 5)
// 	tree.Put(8, 3)
// 	tree.Put(17, 7)
// 	tree.Put(1, 1)
// 	tree.Put(11, 4)
// 	tree.Put(15, 6)
// 	tree.Put(25, 9)
// 	tree.Put(6, 2)
// 	tree.Put(22, 8)
// 	tree.Put(27, 10)
// 	// │           ┌── 27
// 	// │       ┌── 25
// 	// │       │   └── 22
// 	// │   ┌── 17
// 	// │   │   └── 15
// 	// └── 13
// 	//     │   ┌── 11
// 	//     └── 8
// 	//         │   ┌── 6
// 	//         └── 1
// 	it := tree.Iterator()
// 	count := 0
// 	for it.Next() {
// 		count++
// 		value := it.Value()
// 		switch value {
// 		case count:
// 			if actualValue, expectedValue := value, count; actualValue != expectedValue {
// 				t.Errorf("Got %v expected %v", actualValue, expectedValue)
// 			}
// 		default:
// 			if actualValue, expectedValue := value, count; actualValue != expectedValue {
// 				t.Errorf("Got %v expected %v", actualValue, expectedValue)
// 			}
// 		}
// 	}
// 	if actualValue, expectedValue := count, tree.Size(); actualValue != expectedValue {
// 		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
// 	}
// }

// func TestRedBlackTreeIterator4Prev(t *testing.T) {
// 	tree := NewWithIntComparator()
// 	tree.Put(13, 5)
// 	tree.Put(8, 3)
// 	tree.Put(17, 7)
// 	tree.Put(1, 1)
// 	tree.Put(11, 4)
// 	tree.Put(15, 6)
// 	tree.Put(25, 9)
// 	tree.Put(6, 2)
// 	tree.Put(22, 8)
// 	tree.Put(27, 10)
// 	// │           ┌── 27
// 	// │       ┌── 25
// 	// │       │   └── 22
// 	// │   ┌── 17
// 	// │   │   └── 15
// 	// └── 13
// 	//     │   ┌── 11
// 	//     └── 8
// 	//         │   ┌── 6
// 	//         └── 1
// 	it := tree.Iterator()
// 	count := tree.Size()
// 	for it.Next() {
// 	}
// 	for it.Prev() {
// 		value := it.Value()
// 		switch value {
// 		case count:
// 			if actualValue, expectedValue := value, count; actualValue != expectedValue {
// 				t.Errorf("Got %v expected %v", actualValue, expectedValue)
// 			}
// 		default:
// 			if actualValue, expectedValue := value, count; actualValue != expectedValue {
// 				t.Errorf("Got %v expected %v", actualValue, expectedValue)
// 			}
// 		}
// 		count--
// 	}
// 	if actualValue, expectedValue := count, 0; actualValue != expectedValue {
// 		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
// 	}
// }

// func TestRedBlackTreeIteratorBegin(t *testing.T) {
// 	tree := NewWithIntComparator()
// 	tree.Put(3, "c")
// 	tree.Put(1, "a")
// 	tree.Put(2, "b")
// 	it := tree.Iterator()

// 	if it.node != nil {
// 		t.Errorf("Got %v expected %v", it.node, nil)
// 	}

// 	it.Begin()

// 	if it.node != nil {
// 		t.Errorf("Got %v expected %v", it.node, nil)
// 	}

// 	for it.Next() {
// 	}

// 	it.Begin()

// 	if it.node != nil {
// 		t.Errorf("Got %v expected %v", it.node, nil)
// 	}

// 	it.Next()
// 	if key, value := it.Key(), it.Value(); key != 1 || value != "a" {
// 		t.Errorf("Got %v,%v expected %v,%v", key, value, 1, "a")
// 	}
// }

// func TestRedBlackTreeIteratorEnd(t *testing.T) {
// 	tree := NewWithIntComparator()
// 	it := tree.Iterator()

// 	if it.node != nil {
// 		t.Errorf("Got %v expected %v", it.node, nil)
// 	}

// 	it.End()
// 	if it.node != nil {
// 		t.Errorf("Got %v expected %v", it.node, nil)
// 	}

// 	tree.Put(3, "c")
// 	tree.Put(1, "a")
// 	tree.Put(2, "b")
// 	it.End()
// 	if it.node != nil {
// 		t.Errorf("Got %v expected %v", it.node, nil)
// 	}

// 	it.Prev()
// 	if key, value := it.Key(), it.Value(); key != 3 || value != "c" {
// 		t.Errorf("Got %v,%v expected %v,%v", key, value, 3, "c")
// 	}
// }

// func TestRedBlackTreeIteratorFirst(t *testing.T) {
// 	tree := NewWithIntComparator()
// 	tree.Put(3, "c")
// 	tree.Put(1, "a")
// 	tree.Put(2, "b")
// 	it := tree.Iterator()
// 	if actualValue, expectedValue := it.First(), true; actualValue != expectedValue {
// 		t.Errorf("Got %v expected %v", actualValue, expectedValue)
// 	}
// 	if key, value := it.Key(), it.Value(); key != 1 || value != "a" {
// 		t.Errorf("Got %v,%v expected %v,%v", key, value, 1, "a")
// 	}
// }

// func TestRedBlackTreeIteratorLast(t *testing.T) {
// 	tree := NewWithIntComparator()
// 	tree.Put(3, "c")
// 	tree.Put(1, "a")
// 	tree.Put(2, "b")
// 	it := tree.Iterator()
// 	if actualValue, expectedValue := it.Last(), true; actualValue != expectedValue {
// 		t.Errorf("Got %v expected %v", actualValue, expectedValue)
// 	}
// 	if key, value := it.Key(), it.Value(); key != 3 || value != "c" {
// 		t.Errorf("Got %v,%v expected %v,%v", key, value, 3, "c")
// 	}
// }
