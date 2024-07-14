package local

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/3rd/core/core-lib/wiki"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalWiki(t *testing.T) {
	testDataPath, err := filepath.Abs("../../test-data/wiki/default")
	require.NoError(t, err)

	config := LocalWikiConfig{
		Root:  testDataPath,
		Parse: true,
	}

	t.Run("NewLocalWiki", func(t *testing.T) {
		localWiki, err := NewLocalWiki(config)
		assert.NoError(t, err)
		assert.NotNil(t, localWiki)
	})

	t.Run("GetNodes", func(t *testing.T) {
		localWiki, err := NewLocalWiki(config)
		require.NoError(t, err)

		nodes, err := localWiki.GetNodes()
		assert.NoError(t, err)
		assert.Len(t, nodes, 4)
	})

	t.Run("FindNode", func(t *testing.T) {
		localWiki, err := NewLocalWiki(config)
		require.NoError(t, err)

		t.Run("Find existing node", func(t *testing.T) {
			node, err := localWiki.FindNode(func(n wiki.Node) bool {
				return n.GetName() == "Custom title"
			})
			assert.NoError(t, err)
			assert.NotNil(t, node)
			assert.Equal(t, "Custom title", node.GetName())
		})

		t.Run("Find non-existing node", func(t *testing.T) {
			node, err := localWiki.FindNode(func(n wiki.Node) bool {
				return n.GetName() == "Non-existing node"
			})
			assert.NoError(t, err)
			assert.Nil(t, node)
		})
	})

	t.Run("GetNode", func(t *testing.T) {
		localWiki, err := NewLocalWiki(config)
		require.NoError(t, err)

		t.Run("Get existing node", func(t *testing.T) {
			node := localWiki.GetNode("root-1")
			assert.NotNil(t, node)
			assert.Equal(t, "root-1", node.GetName())
		})

		t.Run("Get non-existing node", func(t *testing.T) {
			node := localWiki.GetNode("non-existing")
			assert.Nil(t, node)
		})
	})

	t.Run("Refresh", func(t *testing.T) {
		localWiki, err := NewLocalWiki(config)
		require.NoError(t, err)
		nodes, err := localWiki.GetNodes()
		assert.NoError(t, err)
		assert.Len(t, nodes, 4)

		tmpPath := filepath.Join(testDataPath, "tmp")
		err = os.WriteFile(tmpPath, []byte(""), 0644)
		require.NoError(t, err)

		err = localWiki.Refresh()
		assert.NoError(t, err)
		nodes, err = localWiki.GetNodes()
		assert.NoError(t, err)
		assert.Len(t, nodes, 5)

		err = os.Remove(tmpPath)
		require.NoError(t, err)
	})
}
