package local

import (
	"path/filepath"
	"testing"

	"github.com/3rd/core/core-lib/wiki"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalNode(t *testing.T) {
	t.Run("Resolve basic properties", func(t *testing.T) {
		path, err := filepath.Abs("../../test-data/wiki/default/root-1")
		require.NoError(t, err)

		node, err := NewLocalNode(path)
		assert.NoError(t, err)
		assert.Equal(t, wiki.NodeID("root-1"), node.GetID())
		assert.Equal(t, "root-1", node.GetName())
		assert.Equal(t, path, node.GetPath())
		assert.False(t, node.IsParsed())

		content, err := node.GetContent()
		assert.NoError(t, err)
		assert.Equal(t, "This is the root-1 node.\n", content)
	})

	t.Run("Parse node and resolve custom title", func(t *testing.T) {
		path, err := filepath.Abs("../../test-data/wiki/default/root-2")
		require.NoError(t, err)

		node, err := NewLocalNode(path)
		require.NoError(t, err)

		err = node.Parse()
		assert.NoError(t, err)
		assert.True(t, node.IsParsed())
		assert.Equal(t, "Custom title", node.GetName())
	})
}
