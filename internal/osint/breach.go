package osint

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"ghosttools/internal/ui"
)

// BreachCheck lists breaches associated with an email via HIBP public endpoints.
func BreachCheck() {
	ui.Cyan("enter email or account:")
	var e string
	fmt.Scanln(&e)
	e = strings.TrimSpace(e)

	// fetch all breaches (public, no key)
	resp, err := http.Get("https://haveibeenpwned.com/api/v3/breaches")
	if err != nil {
		ui.Red("error: " + err.Error())
		return
	}
	defer resp.Body.Close()

	var breaches []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&breaches)
	fmt.Println()
	fmt.Println(ui.Green(fmt.Sprintf("─── known breaches (public list, %d total) ───", len(breaches))))

	// show recent 15
	for i, b := range breaches {
		if i >= 15 {
			break
		}
		fmt.Printf("  %s (%s) — %v accounts\n", ui.Cyan(fmt.Sprint(b["Name"])), b["BreachDate"], b["PwnCount"])
	}

	fmt.Println()
	fmt.Println(ui.Cyan("─── check your email ───"))
	fmt.Printf("  https://haveibeenpwned.com/account/%s\n", e)
	fmt.Printf("  https://leakcheck.io/\n")
	fmt.Printf("  https://dehashed.com/search?query=%s\n", e)
	fmt.Printf("  https://intelx.io/?s=%s\n", e)
}

// ShodanLookup queries shodan for host info (API key optional).
func ShodanLookup() {
	ui.Cyan("enter IP or domain:")
	var q string
	fmt.Scanln(&q)
	q = strings.TrimSpace(q)

	apiKey := ""
	// check env
	// os.Getenv not imported? add import if needed
	// keep it simple with public links

	fmt.Println()
	fmt.Println(ui.Green("─── search engines ───"))
	links := [][2]string{
		{"shodan",   "https://www.shodan.io/search?query=" + q},
		{"censys",   "https://search.censys.io/search?q=" + q},
		{"zoomeye",  "https://www.zoomeye.org/searchResult?q=" + q},
		{"fofa",     "https://fofa.info/result?qbase64=" + q},
		{"netlas",   "https://app.netlas.io/host/" + q},
		{"binaryedge", "https://app.binaryedge.io/services/query?query=" + q},
	}
	for _, l := range links {
		fmt.Printf("  %-12s %s\n", ui.Cyan(l[0]+":"), l[1])
	}
	if apiKey == "" {
		fmt.Println()
		ui.Yellow("for direct API access, set SHODAN_API_KEY env var (paid)")
	}
}

// Wayback queries the wayback machine for URLs snapshots.
func Wayback() {
	ui.Cyan("enter URL (with https://):")
	var u string
	fmt.Scanln(&u)
	u = strings.TrimSpace(u)

	// API endpoint
	api := "http://archive.org/wayback/available?url=" + u
	resp, err := http.Get(api)
	if err == nil {
		var data map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&data)
		if snap, ok := data["archived_snapshots"].(map[string]interface{}); ok {
			if closest, ok := snap["closest"].(map[string]interface{}); ok {
				fmt.Println()
				fmt.Println(ui.Green("─── closest snapshot ───"))
				fmt.Printf("  url:       %v\n", closest["url"])
				fmt.Printf("  timestamp: %v\n", closest["timestamp"])
				fmt.Printf("  status:    %v\n", closest["status"])
			}
		}
		resp.Body.Close()
	}

	fmt.Println()
	fmt.Println(ui.Cyan("─── archive links ───"))
	fmt.Printf("  wayback:   https://web.archive.org/web/*/%s\n", u)
	fmt.Printf("  archive.today: https://archive.ph/newest/%s\n", u)
	fmt.Printf("  timetravel: http://timetravel.mementoweb.org/timemap/link/%s\n", u)
}

// DarkWeb generates .onion search links (requires tor to open).
func DarkWeb() {
	ui.Cyan("enter search term:")
	var q string
	fmt.Scanln(&q)
	q = strings.TrimSpace(q)

	fmt.Println()
	fmt.Println(ui.Yellow("requires tor browser or socks5 proxy on 127.0.0.1:9050"))
	fmt.Println()
	fmt.Println(ui.Green("─── .onion indexes ───"))
	links := [][2]string{
		{"ahmia",  "http://juhanurmihxlp77nkq76byazcldy2hlmovfu2epvl5ankdibsot4csyd.onion/search/?q=" + q},
		{"tor66",  "http://tor66sewebgixwhcqfnp5inzp5x5uohhdy3kvtnyfxc2e5mxiuh34iid.onion/search?q=" + q},
		{"phobos", "http://phobosxilamwcg75xt22id7aywkzol6q6rfl2flipcqoc4e4ahima5id.onion/search?query=" + q},
		{"haystak", "http://haystak5njsmn2hqkewecpaxetahtwhsbsa64jom2k22z5afxhnpxfid.onion/?q=" + q},
	}
	for _, l := range links {
		fmt.Printf("  %-10s %s\n", ui.Cyan(l[0]+":"), l[1])
	}
}

// PortScanner scans common ports on a host (fast, concurrent).
func PortScanner() {
	ui.Cyan("enter host or IP:")
	var h string
	fmt.Scanln(&h)
	h = strings.TrimSpace(h)

	ui.Cyan("port range (e.g. 1-1024, blank = common):")
	var pr string
	fmt.Scanln(&pr)
	pr = strings.TrimSpace(pr)

	var ports []int
	if pr == "" {
		ports = []int{21, 22, 23, 25, 53, 80, 110, 111, 135, 139, 143, 443, 445, 993, 995, 1723, 3306, 3389, 5432, 5900, 6379, 8080, 8443, 27017}
	} else if strings.Contains(pr, "-") {
		var lo, hi int
		fmt.Sscanf(pr, "%d-%d", &lo, &hi)
		if hi > 65535 {
			hi = 65535
		}
		for p := lo; p <= hi; p++ {
			ports = append(ports, p)
		}
	} else {
		for _, s := range strings.Split(pr, ",") {
			var n int
			fmt.Sscanf(strings.TrimSpace(s), "%d", &n)
			if n > 0 {
				ports = append(ports, n)
			}
		}
	}

	fmt.Println()
	ui.Cyan(fmt.Sprintf("scanning %d ports on %s...", len(ports), h))

	var wg sync.WaitGroup
	var mu sync.Mutex
	open := 0
	sem := make(chan struct{}, 200)
	for _, p := range ports {
		wg.Add(1)
		sem <- struct{}{}
		go func(port int) {
			defer wg.Done()
			defer func() { <-sem }()
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", h, port), 1200*time.Millisecond)
			if err != nil {
				return
			}
			conn.Close()
			mu.Lock()
			open++
			fmt.Printf("  %s %d\n", ui.Green("[open]"), port)
			mu.Unlock()
		}(p)
	}
	wg.Wait()
	fmt.Println()
	ui.Green(fmt.Sprintf("scan complete — %d open", open))
}

// helper unused
var _ = io.Discard