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
		assert.Equal(t, "root-1", node.GetID())
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

		err = node.Parse("full")
		assert.NoError(t, err)
		assert.True(t, node.IsParsed())
		assert.Equal(t, "Custom title", node.GetName())
	})

	t.Run("Parse only node meta and resolve custom title", func(t *testing.T) {
		path, err := filepath.Abs("../../test-data/wiki/default/root-2")
		require.NoError(t, err)

		node, err := NewLocalNode(path)
		require.NoError(t, err)

		err = node.Parse("meta")
		assert.NoError(t, err)
		assert.True(t, node.IsParsed())
		assert.Equal(t, "Custom title", node.GetName())
	})

	t.Run("Get tasks", func(t *testing.T) {
		path, err := filepath.Abs("../../test-data/wiki/tasks/sample")
		require.NoError(t, err)

		node, err := NewLocalNode(path)
		require.NoError(t, err)

		err = node.Parse("full")
		assert.NoError(t, err)
		assert.True(t, node.IsParsed())

		tasks := node.GetTasks()
		assert.Len(t, tasks, 4)

		// [ ] task 1
		assert.Equal(t, "task 1", tasks[0].Text)
		assert.Equal(t, wiki.TASK_STATUS_DEFAULT, tasks[0].Status)
		assert.Equal(t, uint32(2), tasks[0].LineNumber)
		assert.Equal(t, uint32(0), tasks[0].Priority)

		// [-] task 2
		assert.Equal(t, "task 2", tasks[1].Text)
		assert.Equal(t, wiki.TASK_STATUS_ACTIVE, tasks[1].Status)
		assert.Equal(t, uint32(3), tasks[1].LineNumber)
		assert.Equal(t, uint32(0), tasks[1].Priority)

		// [x] task 2-2 Session: 2024.01.01 01:00-02:00
		assert.Equal(t, "task 2-2", tasks[2].Text)
		assert.Equal(t, wiki.TASK_STATUS_DONE, tasks[2].Status)
		assert.Equal(t, uint32(4), tasks[2].LineNumber)
		assert.Equal(t, uint32(0), tasks[2].Priority)
		assert.Equal(t, 1, len(tasks[2].Sessions))
		assert.NotNil(t, tasks[2].GetLastWorkSession())

		// [ ] task 3 Schedule: 2024.01.01 10:00
		assert.Equal(t, "task 3", tasks[3].Text)
		assert.Equal(t, wiki.TASK_STATUS_DEFAULT, tasks[3].Status)
		assert.Equal(t, uint32(6), tasks[3].LineNumber)
		assert.Equal(t, uint32(0), tasks[3].Priority)
		assert.NotNil(t, tasks[3].Schedule)
	})
}
