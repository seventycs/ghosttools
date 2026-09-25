package osint

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"ghosttools/internal/ui"
)

// EmailLookup checks an email against many services to detect registrations.
func EmailLookup() {
	ui.Cyan("enter email address:")
	var e string
	fmt.Scanln(&e)
	e = strings.TrimSpace(strings.ToLower(e))

	// email validation
	if !strings.Contains(e, "@") {
		ui.Red("invalid email")
		return
	}

	hash := md5.Sum([]byte(e))
	hashStr := hex.EncodeToString(hash[:])

	fmt.Println()
	fmt.Println(ui.Green("─── email intel ───"))
	fmt.Printf("  %-14s %s\n", "email:", e)
	fmt.Printf("  %-14s %s\n", "md5:", hashStr)
	fmt.Printf("  %-14s https://www.gravatar.com/avatar/%s\n", "gravatar:", hashStr)

	// holehe-style checks against public signup/password-reset endpoints
	sites := []struct {
		name string
		url  string
	}{
		{"instagram", "https://www.instagram.com/accounts/account_recovery_send_ajax/"},
		{"twitter",   "https://api.twitter.com/i/users/email_available.json?email=" + e},
		{"github",    "https://github.com/password_reset"},
		{"facebook",  "https://www.facebook.com/login/identify"},
		{"linkedin",  "https://www.linkedin.com/checkpoint/lg/login-submit"},
		{"pinterest", "https://www.pinterest.com/resource/EmailExistsResource/get/"},
		{"spotify",   "https://www.spotify.com/api/signup/validate"},
		{"snapchat",  "https://accounts.snapchat.com/accounts/merlin/login"},
	}

	client := &http.Client{Timeout: 8 * time.Second}
	fmt.Println()
	fmt.Println(ui.Cyan("─── checking registrations ───"))

	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, s := range sites {
		wg.Add(1)
		go func(name, url string) {
			defer wg.Done()
			req, err := http.NewRequest("POST", url, nil)
			if err != nil {
				return
			}
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			resp, err := client.Do(req)
			if err != nil {
				return
			}
			resp.Body.Close()
			// 200 or 302 without "not found" usually means "exists"
			status := resp.StatusCode
			mu.Lock()
			defer mu.Unlock()
			switch status {
			case 200, 302:
				fmt.Printf("  %-14s %s\n", ui.Green("[?]"), name)
			default:
				fmt.Printf("  %-14s %s\n", ui.Red("[-]"), name)
			}
		}(s.name, s.url)
	}
	wg.Wait()

	fmt.Println()
	fmt.Println(ui.Cyan("─── osint links ───"))
	links := [][2]string{
		{"haveibeenpwned", "https://haveibeenpwned.com/account/" + e},
		{"hunter",         "https://hunter.io/email-verifier/" + e},
		{"epieos",         "https://epieos.com/?q=" + e},
		{"ghunt",          "https://github.com/mxrch/GHunt"},
		{"emailrep",       "https://emailrep.io/" + e},
		{"dehashed",       "https://dehashed.com/search?query=" + e},
	}
	for _, l := range links {
		fmt.Printf("  %-16s %s\n", ui.Cyan(l[0]+":"), l[1])
	}
}