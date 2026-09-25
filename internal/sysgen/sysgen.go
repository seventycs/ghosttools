package sysgen

import (
	"bufio"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"ghosttools/internal/ui"
)

func Base64Codec() {
	fmt.Println("  [1] encode")
	fmt.Println("  [2] decode")
	fmt.Print("choice: ")
	var c string
	fmt.Scanln(&c)
	fmt.Print("input: ")
	reader := bufio.NewReader(os.Stdin)
	in, _ := reader.ReadString('\n')
	in = strings.TrimSpace(in)
	if c == "1" {
		fmt.Println(base64.StdEncoding.EncodeToString([]byte(in)))
	} else {
		b, err := base64.StdEncoding.DecodeString(in)
		if err != nil {
			ui.Red(err.Error())
			return
		}
		fmt.Println(string(b))
	}
}

func SystemInfo() {
	host, _ := os.Hostname()
	wd, _ := os.Getwd()
	fmt.Printf("  hostname: %s\n", host)
	fmt.Printf("  os:       %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("  cpus:     %d\n", runtime.NumCPU())
	fmt.Printf("  cwd:      %s\n", wd)
	fmt.Printf("  user:     %s\n", os.Getenv("USERNAME"))
}

func IPPinger() {
	ui.Cyan("enter host:")
	var h string
	fmt.Scanln(&h)
	for i := 0; i < 4; i++ {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", h+":80", 3*time.Second)
		if err != nil {
			fmt.Printf("  fail: %v\n", err)
			continue
		}
		conn.Close()
		fmt.Printf("  %s: %v\n", h, time.Since(start))
	}
}

func Obfuscator() {
	ui.Cyan("enter python file path:")
	var p string
	fmt.Scanln(&p)
	p = strings.TrimSpace(p)
	b, err := os.ReadFile(p)
	if err != nil {
		ui.Red(err.Error())
		return
	}
	encoded := base64.StdEncoding.EncodeToString(b)
	out := "import base64;exec(base64.b64decode('" + encoded + "').decode())"
	outPath := strings.TrimSuffix(p, ".py") + "_obf.py"
	os.WriteFile(outPath, []byte(out), 0644)
	ui.Green("saved to " + outPath)
}

func MetadataScan() {
	ui.Cyan("enter file path:")
	var p string
	fmt.Scanln(&p)
	p = strings.TrimSpace(p)
	stat, err := os.Stat(p)
	if err != nil {
		ui.Red(err.Error())
		return
	}
	fmt.Printf("  size: %d bytes\n", stat.Size())
	fmt.Printf("  mode: %v\n", stat.Mode())
	fmt.Printf("  mod:  %v\n", stat.ModTime())
}

func AppInfo() {
	fmt.Println("  ghosttools v1.0.0")
	fmt.Println("  made by seventycs")
	fmt.Println("  go runtime:", runtime.Version())
}

func AppConfig() {
	fmt.Println("  config: ./config.json")
	fmt.Println("  themes: green/cyan (default)")
}

func DisableAV() {
	if runtime.GOOS != "windows" {
		ui.Red("windows only")
		return
	}
	wd, _ := os.Getwd()
	exec.Command("powershell", "-Command",
		"Add-MpPreference -ExclusionPath '"+wd+"'").Run()
	ui.Green("exclusion added (requires admin)")
}

func ProxyScraper() {
	ui.Cyan("scraping proxies...")
	sources := []string{
		"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/http.txt",
		"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/socks4.txt",
		"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/socks5.txt",
	}
	var all []string
	for _, s := range sources {
		resp, err := http.Get(s)
		if err != nil {
			continue
		}
		body := make([]byte, 8192)
		n, _ := resp.Body.Read(body)
		resp.Body.Close()
		all = append(all, strings.Split(string(body[:n]), "\n")...)
	}
	os.WriteFile("proxies.txt", []byte(strings.Join(all, "\n")), 0644)
	ui.Green(fmt.Sprintf("scraped %d proxies to proxies.txt", len(all)))
}

func ProxyChecker() {
	b, err := os.ReadFile("proxies.txt")
	if err != nil {
		ui.Red("run proxy scraper first")
		return
	}
	lines := strings.Split(string(b), "\n")
	ui.Cyan(fmt.Sprintf("checking %d proxies...", len(lines)))
	good := []string{}
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		conn, err := net.DialTimeout("tcp", l, 2*time.Second)
		if err == nil {
			conn.Close()
			good = append(good, l)
		}
	}
	os.WriteFile("proxies_good.txt", []byte(strings.Join(good, "\n")), 0644)
	ui.Green(fmt.Sprintf("%d working", len(good)))
}

func RegistryEditor() {
	ui.Yellow("use regedit.exe for manual registry edits")
}

func ServiceManager() {
	out, _ := exec.Command("sc", "query", "type=", "service", "state=", "all").CombinedOutput()
	fmt.Println(string(out))
}

func StartupManager() {
	out, _ := exec.Command("wmic", "startup", "list", "brief").CombinedOutput()
	fmt.Println(string(out))
}

func FileShredder() {
	ui.Cyan("enter file path:")
	var p string
	fmt.Scanln(&p)
	p = strings.TrimSpace(p)
	b, err := os.ReadFile(p)
	if err != nil {
		ui.Red(err.Error())
		return
	}
	for i := 0; i < 3; i++ {
		rand.Read(b)
		os.WriteFile(p, b, 0644)
	}
	os.Remove(p)
	ui.Green("shredded")
}

func PasswordGen() {
	ui.Cyan("length (default 24):")
	var n int
	fmt.Scanln(&n)
	if n == 0 {
		n = 24
	}
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
	b := make([]byte, n)
	rand.Read(b)
	for i := range b {
		b[i] = chars[int(b[i])%len(chars)]
	}
	fmt.Println(ui.Green(string(b)))
}

func HashCracker() {
	ui.Cyan("enter hash:")
	var h string
	fmt.Scanln(&h)
	ui.Cyan("enter wordlist path:")
	var w string
	fmt.Scanln(&w)
	f, err := os.Open(strings.TrimSpace(w))
	if err != nil {
		ui.Red(err.Error())
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		word := sc.Text()
		m := md5.Sum([]byte(word))
		if hex.EncodeToString(m[:]) == h {
			ui.Green("md5: " + word)
			return
		}
		s := sha1.Sum([]byte(word))
		if hex.EncodeToString(s[:]) == h {
			ui.Green("sha1: " + word)
			return
		}
		s2 := sha256.Sum256([]byte(word))
		if hex.EncodeToString(s2[:]) == h {
			ui.Green("sha256: " + word)
			return
		}
	}
	ui.Red("not found")
}

func PortForward() {
	ui.Yellow("use ssh -L or ncat for port forwarding")
}

func PacketSniffer() {
	ui.Yellow("requires npcap + admin — use wireshark")
}

func WiFiScanner() {
	out, _ := exec.Command("netsh", "wlan", "show", "networks", "mode=bssid").CombinedOutput()
	fmt.Println(string(out))
}

func Debloater() {
	Debloat()
}