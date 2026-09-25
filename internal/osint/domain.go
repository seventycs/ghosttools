package osint

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"ghosttools/internal/ui"
)

// DNSLookup prints A, MX, TXT, NS, CNAME records for a domain.
func DNSLookup() {
	ui.Cyan("enter domain:")
	var d string
	fmt.Scanln(&d)
	d = strings.TrimSpace(d)

	fmt.Println()
	fmt.Println(ui.Green("A records:"))
	if ips, err := net.LookupHost(d); err == nil {
		for _, ip := range ips {
			fmt.Println("  " + ip)
		}
	} else {
		fmt.Println("  (none)")
	}
	fmt.Println(ui.Green("\nMX records:"))
	if mxs, err := net.LookupMX(d); err == nil {
		for _, mx := range mxs {
			fmt.Printf("  %s (pref %d)\n", mx.Host, mx.Pref)
		}
	} else {
		fmt.Println("  (none)")
	}
	fmt.Println(ui.Green("\nTXT records:"))
	if txts, err := net.LookupTXT(d); err == nil {
		for _, t := range txts {
			fmt.Println("  " + t)
		}
	} else {
		fmt.Println("  (none)")
	}
	fmt.Println(ui.Green("\nNS records:"))
	if nss, err := net.LookupNS(d); err == nil {
		for _, ns := range nss {
			fmt.Println("  " + ns.Host)
		}
	} else {
		fmt.Println("  (none)")
	}
	fmt.Println(ui.Green("\nCNAME:"))
	if cname, err := net.LookupCNAME(d); err == nil {
		fmt.Println("  " + cname)
	} else {
		fmt.Println("  (none)")
	}
}

// WhoisLookup queries whois servers for the given domain.
func WhoisLookup() {
	ui.Cyan("enter domain:")
	var d string
	fmt.Scanln(&d)
	d = strings.TrimSpace(d)

	// find whois server via IANA
	conn, err := net.Dial("tcp", "whois.iana.org:43")
	if err != nil {
		ui.Red("connect failed: " + err.Error())
		return
	}
	fmt.Fprintf(conn, "%s\r\n", d)
	buf := make([]byte, 4096)
	n, _ := conn.Read(buf)
	conn.Close()
	server := ""
	for _, line := range strings.Split(string(buf[:n]), "\n") {
		if strings.HasPrefix(strings.ToLower(line), "refer:") || strings.HasPrefix(strings.ToLower(line), "whois:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				server = strings.TrimSpace(parts[1])
				break
			}
		}
	}
	if server == "" {
		server = "whois.verisign-grs.com"
	}

	conn2, err := net.Dial("tcp", server+":43")
	if err != nil {
		ui.Red("whois connect failed: " + err.Error())
		return
	}
	defer conn2.Close()
	fmt.Fprintf(conn2, "domain %s\r\n", d)
	sc := bufio.NewScanner(conn2)
	for sc.Scan() {
		fmt.Println(sc.Text())
	}
}

// DomainRep shows reputation links and public data on a domain.
func DomainRep() {
	ui.Cyan("enter domain:")
	var d string
	fmt.Scanln(&d)
	d = strings.TrimSpace(d)

	fmt.Println()
	fmt.Println(ui.Green("─── reputation ───"))
	links := [][2]string{
		{"virustotal",  "https://www.virustotal.com/gui/domain/" + d},
		{"urlscan",     "https://urlscan.io/domain/" + d},
		{"shodan",      "https://www.shodan.io/domain/" + d},
		{"urlhaus",     "https://urlhaus.abuse.ch/browse.php?search=" + d},
		{"talos",       "https://talosintelligence.com/reputation_center/lookup?search=" + d},
		{"abuseipdb",   "https://www.abuseipdb.com/check/" + d},
		{"crt.sh",      "https://crt.sh/?q=" + d},
		{"securitytrails", "https://securitytrails.com/domain/" + d},
	}
	for _, l := range links {
		fmt.Printf("  %-16s %s\n", ui.Cyan(l[0]+":"), l[1])
	}
}

// SSLCert fetches cert transparency records via crt.sh.
func SSLCert() {
	ui.Cyan("enter domain:")
	var d string
	fmt.Scanln(&d)
	d = strings.TrimSpace(d)

	resp, err := http.Get("https://crt.sh/?q=" + d + "&output=json")
	if err != nil {
		ui.Red("error: " + err.Error())
		return
	}
	defer resp.Body.Close()

	var certs []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&certs); err != nil {
		ui.Red("parse error: " + err.Error())
		return
	}

	fmt.Println()
	fmt.Println(ui.Green(fmt.Sprintf("─── cert transparency (%d entries) ───", len(certs))))
	seen := map[string]bool{}
	for i, c := range certs {
		if i >= 30 {
			break
		}
		name := fmt.Sprint(c["name_value"])
		if seen[name] {
			continue
		}
		seen[name] = true
		fmt.Printf("  %s issued by %s\n", name, c["issuer_name"])
	}
}

// SubdomainEnum enumerates subdomains via crt.sh + hackertarget.
func SubdomainEnum() {
	ui.Cyan("enter domain:")
	var d string
	fmt.Scanln(&d)
	d = strings.TrimSpace(d)

	fmt.Println()
	fmt.Println(ui.Green("─── subdomains ───"))

	// hackertarget
	resp, err := http.Get("https://api.hackertarget.com/hostsearch/?q=" + d)
	if err == nil {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		for _, line := range strings.Split(string(body), "\n") {
			if line == "" {
				continue
			}
			fmt.Println("  " + strings.Split(line, ",")[0])
		}
	}

	// crt.sh
	resp2, err := http.Get("https://crt.sh/?q=%25." + d + "&output=json")
	if err == nil {
		defer resp2.Body.Close()
		var certs []map[string]interface{}
		if json.NewDecoder(resp2.Body).Decode(&certs) == nil {
			seen := map[string]bool{}
			for _, c := range certs {
				name := fmt.Sprint(c["name_value"])
				for _, n := range strings.Split(name, "\n") {
					if !seen[n] && strings.HasSuffix(n, d) {
						seen[n] = true
						fmt.Println("  " + n)
					}
				}
			}
		}
	}
}