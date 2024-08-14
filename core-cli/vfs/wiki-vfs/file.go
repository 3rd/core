package wikivfs

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"bazil.org/fuse"
	"bazil.org/fuse/fuseutil"
	"github.com/3rd/core/core-lib/wiki/local"
)

type WikiVFSFile struct {
	wiki     *local.LocalWiki
	path     string
	markdown *string
}

var (
	markdownCache = make(map[string]*string)
	cacheMutex    sync.RWMutex
)

func (w WikiVFSFile) Convert() error {
	cacheMutex.RLock()
	cachedMarkdown, exists := markdownCache[w.path]
	cacheMutex.RUnlock()

	if exists {
		w.markdown = cachedMarkdown
		return nil
	}

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	path := strings.TrimSuffix(w.path, ".md")
	wikiNodes, err := w.wiki.GetNodes()
	if err != nil {
		return fmt.Errorf("failed to get wiki nodes: %w", err)
	}

	for _, wikiNode := range wikiNodes {
		if path == wikiNode.GetPath() {
			wikiNode.Parse("full")
			markdown := wikiNode.ToMarkdown()
			markdownCache[w.path] = &markdown
			w.markdown = &markdown
			return nil
		}
	}

	return fmt.Errorf("no matching wiki node found for path: %s", path)
}

func (w WikiVFSFile) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 2
	a.Mode = 0o700
	a.Uid = uint32(os.Getuid())
	a.Gid = uint32(os.Getgid())

	if err := w.Convert(); err != nil {
		return fmt.Errorf("failed to convert: %w", err)
	}

	cacheMutex.RLock()
	markdown := markdownCache[w.path]
	cacheMutex.RUnlock()

	if markdown != nil {
		a.Size = uint64(len(*markdown))
	} else {
		stat, err := os.Stat(w.path)
		if err != nil {
			return fmt.Errorf("failed to stat file: %w", err)
		}
		a.Size = uint64(stat.Size())
	}

	return nil
}

func (w WikiVFSFile) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	if err := w.Convert(); err != nil {
		return fmt.Errorf("failed to convert: %w", err)
	}

	cacheMutex.RLock()
	markdown := markdownCache[w.path]
	cacheMutex.RUnlock()

	if markdown != nil {
		fuseutil.HandleRead(req, resp, []byte(*markdown))
		return nil
	}

	// fallback to real content
	content, err := os.ReadFile(w.path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	fuseutil.HandleRead(req, resp, content)
	return nil
}
