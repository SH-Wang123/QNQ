## Init
need 64bit gcc
1. `go mod init window_handler`
2. `go get fyne.io/fyne/v2`
2. `go get github.com/shirou/gopsutil`
3. `go get github.com/yusufpapurcu/wmi`
4. `go get -u github.com/gin-gonic/gin`
5. `go get -u github.com/spf13/cobra@latest`
2. `go mod tidy`

## Build
go build -ldflags -H=windowsgui main.go
$Env:GOOS = "linux"
$Env:GOARCH = "arm64"
go build
go build -ldflags="-s -w " -a -v -o  qnq_v0.0.4.5_windows_amd64.exe main.go