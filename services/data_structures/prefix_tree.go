package data_structures

import (
	"fmt"
	"strings"
)

type PrefixTree struct {
	root *prefixTreeNode
}

type prefixTreeNode struct {
	children map[string]*prefixTreeNode
	terminal bool
}

func NewPrefixTree() *PrefixTree {
	return &PrefixTree{
		root: &prefixTreeNode{
			children: map[string]*prefixTreeNode{},
		},
	}
}

func (t *PrefixTree) Insert(key string) error {
	parts, err := parseDotPrefixedKey(key)
	if err != nil {
		return err
	}

	t.insertParts(parts)
	return nil
}

func (t *PrefixTree) InsertAll(keys []string) error {
	parsedKeys := make([][]string, len(keys))
	for i, key := range keys {
		parts, err := parseDotPrefixedKey(key)
		if err != nil {
			return fmt.Errorf("invalid key at index %d: %w", i, err)
		}
		parsedKeys[i] = parts
	}

	for _, parts := range parsedKeys {
		t.insertParts(parts)
	}
	return nil
}

func (t *PrefixTree) Contains(key string) (bool, error) {
	parts, err := parseDotPrefixedKey(key)
	if err != nil {
		return false, err
	}

	current := t.root
	for _, part := range parts {
		next, exists := current.children[part]
		if !exists {
			return false, nil
		}
		current = next
	}
	return current.terminal, nil
}

func (t *PrefixTree) HasPrefix(key string) (bool, error) {
	parts, err := parseDotPrefixedKey(key)
	if err != nil {
		return false, err
	}

	current := t.root
	for _, part := range parts {
		next, exists := current.children[part]
		if !exists {
			return false, nil
		}
		current = next
		if current.terminal {
			return true, nil
		}
	}
	return false, nil
}

func (t *PrefixTree) insertParts(parts []string) {
	current := t.root
	for _, part := range parts {
		next, exists := current.children[part]
		if !exists {
			next = &prefixTreeNode{
				children: map[string]*prefixTreeNode{},
			}
			current.children[part] = next
		}
		current = next
	}
	current.terminal = true
}

func parseDotPrefixedKey(key string) ([]string, error) {
	if key == "" {
		return nil, fmt.Errorf("key cannot be empty")
	}
	if !strings.HasPrefix(key, ".") {
		return nil, fmt.Errorf("key must start with dot")
	}
	if key == "." {
		return nil, fmt.Errorf("key must include at least one segment")
	}

	parts := strings.Split(key[1:], ".")
	for _, part := range parts {
		if part == "" {
			return nil, fmt.Errorf("key cannot contain empty segments")
		}
	}
	return parts, nil
}
