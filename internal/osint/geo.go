package osint

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"ghosttools/internal/ui"
)

// IPGeo does a full IP intel lookup: location, ISP, ASN, reverse DNS, common ports.
func IPGeo() {
	ui.Cyan("enter IP (blank = your own):")
	var ip string
	fmt.Scanln(&ip)
	ip = strings.TrimSpace(ip)

	url := "http://ip-api.com/json/" + ip + "?fields=status,message,country,regionName,city,zip,lat,lon,timezone,isp,org,as,asname,reverse,query"
	resp, err := http.Get(url)
	if err != nil {
		ui.Red("error: " + err.Error())
		return
	}
	defer resp.Body.Close()

	var d map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&d)

	if d["status"] == "fail" {
		ui.Red("lookup failed: " + fmt.Sprint(d["message"]))
		return
	}

	fmt.Println()
	fmt.Println(ui.Green("─── ip intel ───"))
	keys := []string{"query", "country", "regionName", "city", "zip", "lat", "lon", "timezone", "isp", "org", "as", "asname", "reverse"}
	for _, k := range keys {
		if v, ok := d[k]; ok {
			fmt.Printf("  %-12s %v\n", ui.Cyan(k+":"), v)
		}
	}

	lat := fmt.Sprint(d["lat"])
	lon := fmt.Sprint(d["lon"])
	fmt.Println()
	fmt.Printf("  %s https://www.google.com/maps?q=%s,%s\n", ui.Cyan("map:"), lat, lon)

	// reverse DNS
	ui.Cyan("\nresolving reverse dns...")
	if names, err := net.LookupAddr(ip); err == nil && len(names) > 0 {
		for _, n := range names {
			fmt.Println("  " + n)
		}
	} else if ip == "" {
		fmt.Println("  (skipped — no ip)")
	} else {
		fmt.Println("  (no reverse dns)")
	}

	// quick port scan on common ports
	if ip != "" {
		ui.Cyan("\nscanning common ports...")
		common := []int{21, 22, 23, 25, 53, 80, 110, 143, 443, 445, 3306, 3389, 5432, 5900, 6379, 8080, 8443, 27017}
		var wg sync.WaitGroup
		var mu sync.Mutex
		var open []int
		for _, p := range common {
			wg.Add(1)
			go func(port int) {
				defer wg.Done()
				conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 1500*time.Millisecond)
				if err == nil {
					conn.Close()
					mu.Lock()
					open = append(open, port)
					mu.Unlock()
				}
			}(p)
		}
		wg.Wait()
		if len(open) == 0 {
			fmt.Println("  none open")
		} else {
			for _, p := range open {
				fmt.Printf("  %s %d open\n", ui.Green("[+]"), p)
			}
		}
	}
}