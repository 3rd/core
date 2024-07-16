module core

go 1.22.5

replace github.com/3rd/core/core-lib => ../core-lib

replace github.com/3rd/syslang/go-syslang => ../syslang/go-syslang

replace github.com/3rd/go-futui => ../go-futui

require (
	github.com/3rd/core/core-lib v0.0.0
	github.com/joho/godotenv v1.5.1
	github.com/spf13/cobra v1.8.1
)

require (
	github.com/3rd/syslang/go-syslang v0.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/k0kubun/pp/v3 v3.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/smacker/go-tree-sitter v0.0.0-20240625050157-a31a98a7c0f6 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
)
