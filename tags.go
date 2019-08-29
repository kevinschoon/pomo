package main

import (
	"fmt"
	"strings"
)

type Tags struct {
	keys []string
	kvs  map[string]string
}

func NewTags() *Tags {
	return &Tags{
		keys: []string{},
		kvs:  map[string]string{},
	}
}

func NewTagsFromKVs(kvs []string) (*Tags, error) {
	tags := NewTags()
	for _, pair := range kvs {
		split := strings.Split(pair, "=")
		if len(split) != 2 {
			return nil, fmt.Errorf("bad key pair: %s", pair)
		}
		tags.keys = append(tags.keys, split[0])
		tags.kvs[split[0]] = split[1]
	}
	return tags, nil
}

func (t Tags) Keys() []string {
	return t.keys
}

func (t Tags) Get(key string) string {
	if value, ok := t.kvs[key]; ok {
		return value
	}
	return ""
}

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

func MergeTags(first, second *Tags) *Tags {
	for _, key := range second.Keys() {
		if second.Get(key) == "" {
			first.Delete(key)
		} else {
			first.Set(key, second.Get(key))
		}
	}
	return first
}
