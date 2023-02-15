## Init
need 64bit gcc
1. `go mod init window_handler`
2. `go get fyne.io/fyne/v2`
2. `go get github.com/shirou/gopsutil`
2. `go mod tidy`

## Build
go install github.com/fyne-io/fyne-cross@latest
go build -ldflags " -s -w -H=windowsgui" -o ./bin/QNQ_v0.0.3_linux_x86.exe main.go
$Env:GOOS = "linux"
$Env:GOARCH = "arm64"

优化任务：
1. 优化Worker，使得Worker可以装载任何任务
2. 优化协程池，使得协程池可以执行任何任务，并且有扩容和阻塞功能
