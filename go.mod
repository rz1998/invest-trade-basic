module github.com/rz1998/invest-trade-basic

go 1.21

toolchain go1.23.0

replace (
	github.com/rz1998/invest-basic v0.0.2 => /home/rz1998/workspace/investBasic
	github.com/rz1998/invest-stra-signal v0.0.1 => /home/rz1998/workspace/invest-stra-signal
)

require (
	github.com/rz1998/invest-basic v0.0.2
	github.com/rz1998/invest-stra-signal v0.0.1
	github.com/zeromicro/go-zero v1.6.5
)

require (
	github.com/fatih/color v1.16.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	go.opentelemetry.io/otel v1.19.0 // indirect
	go.opentelemetry.io/otel/trace v1.19.0 // indirect
	go.uber.org/automaxprocs v1.5.3 // indirect
	golang.org/x/sys v0.23.0 // indirect
)
