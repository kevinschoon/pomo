package tags

import (
	"fmt"
	"strings"
)

const deletePattern = "||DELETE||"

// Tags are an ordered collection of key or key-value
// pairs that can be associated with tasks
type Tags struct {
	keys []string
	kvs  map[string]string
}

// New creates a new Tags
func New() *Tags {
	return &Tags{
		keys: []string{},
		kvs:  map[string]string{},
	}
}

// FromMap loads Tags from a string map with order
func FromMap(order []string, kvs map[string]string) *Tags {
	tags := &Tags{
		keys: order,
		kvs:  kvs,
	}
	return tags
}

// FromKVs loads Tags from pairs of key=value strings
func FromKVs(kvs []string) (*Tags, error) {
	tags := New()
	for _, pair := range kvs {
		split := strings.Split(pair, "=")
		if len(split) == 1 {
			tags.Set(split[0], "")
		} else if len(split) == 2 {
			if split[1] == "" {
				split[1] = deletePattern
			}
			tags.Set(split[0], split[1])
		} else {
			return nil, fmt.Errorf("bad key pair: %s", pair)
		}
	}
	return tags, nil
}

// Keys returns ordered keys
func (t Tags) Keys() []string {
	return t.keys
}

// Len returns the length of the Tags
func (t Tags) Len() int {
	return len(t.kvs)
}

// HasTag returns true if there is a matching Tag
func (t Tags) HasTag(key string) bool {
	_, ok := t.kvs[key]
	return ok
}

// Get gets the value of a tag if it exists
func (t Tags) Get(key string) string {
	if value, ok := t.kvs[key]; ok {
		return value
	}
	return ""
}

// Set sets the key-value pair for a tag
func (t *Tags) Set(key, value string) {
	var hasKey bool
	for _, k := range t.keys {
		if k == key {
			hasKey = true
		}
	}
	if hasKey {
		t.kvs[key] = value
	} else {
		t.keys = append(t.keys, key)
		t.kvs[key] = value
	}
}

// Delete removes a key-value pair
func (t *Tags) Delete(key string) {
	var keys []string
	for _, k := range t.keys {
		if k != key {
			keys = append(keys, key)
		}
	}
	t.keys = keys
	delete(t.kvs, key)
}

// Contains returns true if this Tags instance
// contains the subset of other
func (t *Tags) Contains(other *Tags) bool {
	for _, key := range other.Keys() {
		if !(t.Get(key) == other.Get(key)) {
			return false
		}
	}
	return true
}

// Merge joins together two Tag pairs
func Merge(first, second *Tags) *Tags {
	for _, key := range second.Keys() {
		if second.Get(key) == deletePattern {
			first.Delete(key)
		} else {
			first.Set(key, second.Get(key))
		}
	}
	return first
}
