module core

go 1.25.0

replace github.com/3rd/core/core-lib => ../core-lib

replace github.com/3rd/syslang/go-syslang => ../syslang/go-syslang

replace github.com/3rd/go-futui => ../go-futui

require (
	bazil.org/fuse v0.0.0-20230120002735-62a210ff1fd5
	github.com/3rd/core/core-lib v0.0.0
	github.com/3rd/go-futui v0.0.0-20240720131722-26cf9e0a36db
	github.com/atotto/clipboard v0.1.4
	github.com/gdamore/tcell/v2 v2.13.8
	github.com/joho/godotenv v1.5.1
	github.com/radovskyb/watcher v1.0.7
	github.com/spf13/cobra v1.8.1
)

require (
	github.com/3rd/syslang/go-syslang v0.0.0 // indirect
	github.com/gdamore/encoding v1.0.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.4.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/smacker/go-tree-sitter v0.0.0-20240827094217-dd81d9e9be82 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/teacat/noire v1.1.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/term v0.41.0 // indirect
	golang.org/x/text v0.35.0 // indirect
)
