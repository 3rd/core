package wikivfs

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/3rd/core/core-lib/wiki/local"
)

type WikiVFSDir struct {
	wiki *local.LocalWiki
	path string
}

func (WikiVFSDir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Mode = os.ModeDir | 0o700
	a.Uid = uint32(os.Getuid())
	a.Gid = uint32(os.Getgid())
	return nil
}

func (w WikiVFSDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	files, err := os.ReadDir(w.path)
	if err != nil {
		return nil, err
	}

	var dirents []fuse.Dirent

	wikiNodes, err := w.wiki.GetNodes()
	if err != nil {
		return nil, err
	}

	for _, node := range files {
		nodeType := fuse.DT_File
		nodePath := filepath.Join(w.path, node.Name())
		nodeName := node.Name()

		if node.IsDir() {
			nodeType = fuse.DT_Dir
		} else {
			// add '.md" if it's a wiki node
			for _, wikiNode := range wikiNodes {
				if nodePath == wikiNode.GetPath() {
					nodeName = nodeName + ".md"
					break
				}
			}
		}

		dirents = append(dirents, fuse.Dirent{
			Name: nodeName,
			Type: nodeType,
		})
	}

	return dirents, nil
}

func (w WikiVFSDir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	// strip '.md' suffix
	name = strings.TrimSuffix(name, ".md")

	stat, err := os.Stat(filepath.Join(w.path, name))
	if err != nil {
		return nil, syscall.ENOENT
	}
	if stat.IsDir() {
		return WikiVFSDir{
			wiki: w.wiki,
			path: filepath.Join(w.path, name),
		}, nil
	}
	return WikiVFSFile{
		wiki: w.wiki,
		path: filepath.Join(w.path, name),
	}, nil
}
