package mist

import (
	"strings"
	"testing"
)

// BenchmarkAddRemoveSimple
func BenchmarkAddRemoveSimple(b *testing.B) {
	node := newNode()
	keys := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node.Add(keys)
		node.Remove(keys)
	}
}

// BenchmarkMatchSimple
func BenchmarkMatchSimple(b *testing.B) {
	node := newNode()
	keys := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	node.Add(keys)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node.Match(keys)
	}
}

// BenchmarkAddRemoveComplex benchmarks to see how fast mist can add/remove keys to
// a subscription
func BenchmarkAddRemoveComplex(b *testing.B) {
	node := newNode()

	// create a giant slice of random keys
	keys := [][]string{}
	for i := 0; i < b.N; i++ {
		keys = append(keys, []string{randKey(), randKey(), randKey(), randKey(), randKey(), randKey(), randKey(), randKey()})
	}

	b.ResetTimer()

	// add/remove keys
	for _, k := range keys {
		node.Add(k)
		node.Remove(k)
	}
}

// BenchmarkMatchComplex benchmarks to see how fast mist can match a set of keys on a
// subscription
func BenchmarkMatchComplex(b *testing.B) {
	node := newNode()

	// create a giant slice of random keys
	keys := [][]string{}
	for i := 0; i < b.N; i++ {
		keys = append(keys, []string{randKey(), randKey(), randKey(), randKey(), randKey(), randKey(), randKey(), randKey()})
	}

	b.ResetTimer()

	// add/match keys
	for _, k := range keys {
		node.Add(k)
		node.Match(k)
	}
}

// TestEmptySubscription
func TestEmptySubscription(t *testing.T) {
	node := newNode()
	if len(node.ToSlice()) != 0 {
		t.Fatalf("Unexpected tags in new subscription!")
	}
}

// TestAddRemoveSimple
func TestAddRemoveSimple(t *testing.T) {
	node := newNode()

	//
	node.Add([]string{"a"})
	if len(node.ToSlice()) != 1 {
		t.Fatalf("Failed to add node")
	}

	node.Remove([]string{"a"})
	if len(node.ToSlice()) != 0 {
		t.Fatalf("Failed to remove node")
	}
}

// TestAddRemoveComplex
func TestAddRemoveComplex(t *testing.T) {
	node := newNode()

	// add/remove unordered keys; should remove
	node.Add([]string{"a", "b", "c"})
	node.Remove([]string{"c", "b", "a"})
	if len(node.ToSlice()) != 0 {
		t.Fatalf("Failed to remove node")
	}

	// add/remove incomplete keys; should not remove
	node.Add([]string{"a", "b", "c"})
	node.Remove([]string{"a"})
	node.Remove([]string{"b"})
	node.Remove([]string{"c"})
	node.Remove([]string{"a", "b"})
	node.Remove([]string{"b", "c"})
	node.Remove([]string{"a", "c"})
	node.Remove([]string{"b", "c", "d"})
	node.Remove([]string{"a", "b", "c", "d"})
	if len(node.ToSlice()) != 1 {
		t.Fatalf("Node unexpectedly removed")
	}

	// add duplicate keys; should only add once
	node.Add([]string{"a", "b", "c"})
	node.Add([]string{"a", "b", "c"})
	if len(node.ToSlice()) != 1 {
		t.Fatalf("Duplicate nodes added")
	}
	node.Remove([]string{"a", "b", "c"})
	if len(node.ToSlice()) != 0 {
		t.Fatalf("Failed to remove nodes")
	}

	// remove duplicate keys; should only remove once
	node.Add([]string{"a", "b", "c"})
	node.Remove([]string{"a", "b", "c"})
	node.Remove([]string{"c", "b", "a"})
	if len(node.ToSlice()) != 0 {
		t.Fatalf("Failed to remove nodes")
	}

	// add duplicate remote one; should leave no nodes
	node.Add([]string{"a", "b", "c"})
	node.Add([]string{"a", "b", "c"})
	node.Add([]string{"a", "b", "c"})
	node.Remove([]string{"a", "b", "c"})
	if len(node.ToSlice()) != 0 {
		t.Fatalf("Failed to remove nodes")
	}
}

// TestList
func TestList(t *testing.T) {
	node := newNode()

	// test simple list; length should be 1 and value should be "a"
	node.Add([]string{"a"})
	list := node.ToSlice()
	if len(list) != 1 {
		t.Fatalf("Wrong number of keys - Expecting 1 got %v", len(list))
	}
	if len(list[0]) != 1 {
		t.Fatalf("Wrong number of keys - Expecing 2 got %v", len(list[0]))
	}
	if strings.Join(list[0], ",") != "a" {
		t.Fatalf("Wrong tags - Expecing 'a' got %v", list[0])
	}

	node.Add([]string{"a", "b"})
	list = node.ToSlice()
	if len(list) != 2 {
		t.Fatalf("Wrong number of keys - Expecting 2 got %v", len(list))
	}
	if len(list[1]) != 2 {
		t.Fatalf("Wrong number of keys - Expecing 2 got %v", len(list[1]))
	}
	if strings.Join(list[1], ",") != "a,b" {
		t.Fatalf("Wrong tags - Expecing 'a,b' got %v", list[1])
	}

	node.Add([]string{"a", "b", "c"})
	list = node.ToSlice()
	if len(list) != 3 {
		t.Fatalf("wrong length of list. Expecting 3 got %v", len(list))
	}
	if len(list[2]) != 3 {
		t.Fatalf("Wrong number of keys - Expecing 3 got %v", len(list[2]))
	}
	if strings.Join(list[2], ",") != "a,b,c" {
		t.Fatalf("Wrong tags - Expecing 'a,b,c' got %v", list[1])
	}
}

// TestMatchSimple
func TestMatchSimple(t *testing.T) {
	node := newNode()

	// simple match
	node.Add([]string{"a"})
	if !node.Match([]string{"a"}) {
		t.Fatalf("Expected match!")
	}

	//
	node.Add([]string{"a", "b"})
	if !node.Match([]string{"a", "b"}) {
		t.Fatalf("Expected match!")
	}

	//
	node.Add([]string{"a", "b", "c"})
	if !node.Match([]string{"a", "b", "c"}) {
		t.Fatalf("Expected match!")
	}
}

// TestMatchComplex
func TestMatchComplex(t *testing.T) {
	node := newNode()

	// match unordered keys; should match
	node.Add([]string{"a", "b", "c"})
	if !node.Match([]string{"c", "b", "a"}) {
		t.Fatalf("Expected match!")
	}
	node.Remove([]string{"a", "b", "c"})

	// match multiple subs with single match; should match
	node.Add([]string{"a", "b", "e"})
	node.Add([]string{"c"})
	if !node.Match([]string{"a", "b", "c", "d"}) {
		t.Fatalf("Expected match!")
	}
	node.Remove([]string{"a", "b", "e"})
	node.Remove([]string{"c"})

	// match incomplete keys; should not match
	node.Add([]string{"a", "b", "c"})
	if node.Match([]string{"a"}) {
		t.Fatalf("Unexpected match!")
	}
	if node.Match([]string{"a", "b"}) {
		t.Fatalf("Unexpected match!")
	}
	if node.Match([]string{"b", "c"}) {
		t.Fatalf("Unexpected match!")
	}
	if node.Match([]string{"a", "c"}) {
		t.Fatalf("Unexpected match!")
	}
	node.Remove([]string{"a", "b", "c"})

	// match incomplete keys; should not match
	node.Add([]string{"a", "b"})
	node.Add([]string{"c", "d"})
	if node.Match([]string{"b", "c"}) {
		t.Fatalf("Unexpected match!")
	}
	node.Remove([]string{"a", "b"})
	node.Remove([]string{"c", "d"})
}
