package data_structures

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrefixTreeInsertAndContains(t *testing.T) {
	tree := NewPrefixTree()

	require.NoError(t, tree.Insert(".foo.bar"))
	require.NoError(t, tree.Insert(".foo.baz"))

	exists, err := tree.Contains(".foo.bar")
	require.NoError(t, err)
	require.True(t, exists)

	exists, err = tree.Contains(".foo.baz")
	require.NoError(t, err)
	require.True(t, exists)

	exists, err = tree.Contains(".foo")
	require.NoError(t, err)
	require.False(t, exists)

	exists, err = tree.Contains(".foo.qux")
	require.NoError(t, err)
	require.False(t, exists)
}

func TestPrefixTreeInsertRejectsInvalidKeys(t *testing.T) {
	tree := NewPrefixTree()

	require.Error(t, tree.Insert(""))
	require.Error(t, tree.Insert("foo.bar"))
	require.Error(t, tree.Insert("."))
	require.Error(t, tree.Insert(".foo..bar"))
	require.Error(t, tree.Insert(".foo."))
}

func TestPrefixTreeContainsRejectsInvalidKeys(t *testing.T) {
	tree := NewPrefixTree()

	_, err := tree.Contains("")
	require.Error(t, err)

	_, err = tree.Contains("foo.bar")
	require.Error(t, err)

	_, err = tree.Contains(".")
	require.Error(t, err)

	_, err = tree.Contains(".foo..bar")
	require.Error(t, err)
}

func TestPrefixTreeInsertAll(t *testing.T) {
	tree := NewPrefixTree()

	require.NoError(t, tree.InsertAll([]string{".foo.bar", ".foo.baz", ".qux"}))

	exists, err := tree.Contains(".foo.bar")
	require.NoError(t, err)
	require.True(t, exists)

	exists, err = tree.Contains(".foo.baz")
	require.NoError(t, err)
	require.True(t, exists)

	exists, err = tree.Contains(".qux")
	require.NoError(t, err)
	require.True(t, exists)
}

func TestPrefixTreeInsertAllRejectsInvalidKeysWithoutPartialInsert(t *testing.T) {
	tree := NewPrefixTree()

	err := tree.InsertAll([]string{".foo.bar", "foo.baz"})
	require.Error(t, err)

	exists, containsErr := tree.Contains(".foo.bar")
	require.NoError(t, containsErr)
	require.False(t, exists)
}

func TestPrefixTreeHasPrefix(t *testing.T) {
	tree := NewPrefixTree()
	require.NoError(t, tree.InsertAll([]string{".foo", ".alpha.beta"}))

	hasPrefix, err := tree.HasPrefix(".foo.bar.baz")
	require.NoError(t, err)
	require.True(t, hasPrefix)

	hasPrefix, err = tree.HasPrefix(".alpha.beta.gamma")
	require.NoError(t, err)
	require.True(t, hasPrefix)

	hasPrefix, err = tree.HasPrefix(".alpha")
	require.NoError(t, err)
	require.False(t, hasPrefix)

	hasPrefix, err = tree.HasPrefix(".unknown.path")
	require.NoError(t, err)
	require.False(t, hasPrefix)
}

func TestPrefixTreeHasPrefixRejectsInvalidKeys(t *testing.T) {
	tree := NewPrefixTree()

	_, err := tree.HasPrefix("")
	require.Error(t, err)

	_, err = tree.HasPrefix("foo.bar")
	require.Error(t, err)

	_, err = tree.HasPrefix(".")
	require.Error(t, err)
}
