package wikivfs

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	"github.com/3rd/core/core-lib/wiki/local"
)

type WikiFS struct {
	rootPath string
	wiki     *local.LocalWiki
}

func (w WikiFS) Root() (fs.Node, error) {
	return WikiVFSDir{
		path: w.rootPath,
		wiki: w.wiki,
	}, nil
}

type WikiVFS struct {
	rootPath   string
	wiki       *local.LocalWiki
	mountPoint string
	conn       *fuse.Conn
}

func NewWikiVFS(wiki *local.LocalWiki, rootPath string, mountPoint string) (*WikiVFS, error) {
	vfs := WikiVFS{
		rootPath:   rootPath,
		wiki:       wiki,
		mountPoint: mountPoint,
	}
	conn, err := fuse.Mount(mountPoint)
	if err != nil {
		return nil, err
	}
	vfs.conn = conn
	return &vfs, nil
}

func (vfs *WikiVFS) Mount() error {
	err := fs.Serve(vfs.conn, WikiFS{
		rootPath: vfs.rootPath,
		wiki:     vfs.wiki,
	})
	return err
}

func (vfs *WikiVFS) Close() {
	vfs.conn.Close()
}
