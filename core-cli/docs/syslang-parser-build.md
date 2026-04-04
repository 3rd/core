# Syslang Parser Build

`go-syslang` pulls the parser in through cgo:

```go
// #include "../../../tree-sitter-syslang/src/parser.c"
// #include "../../../tree-sitter-syslang/src/scanner.c"
import "C"
```

## Build Behavior

`core-cli/Makefile` handles this in two places:

- `ensure_syslang_parser` regenerates `../syslang/tree-sitter-syslang/src/parser.c` when `grammar.js` is newer.
- `build`, `test`, and `install` all use `-a` so Go fully rebuilds packages that depend on the external cgo sources.

Relevant commands:

```make
make -C ../syslang/tree-sitter-syslang generate
go build -a .
go test -a ./...
go install -a .
```

## Recovery

```bash
make -C ../syslang/tree-sitter-syslang generate
go clean -cache
make build
make install

make test
go vet -a ./...
```

