package osint

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"ghosttools/internal/ui"
)

// UsernameSearch checks a username across 40+ sites concurrently.
func UsernameSearch() {
	ui.Cyan("enter username:")
	var u string
	fmt.Scanln(&u)
	u = strings.TrimSpace(u)

	sites := map[string]string{
		"instagram":  "https://www.instagram.com/%s/",
		"twitter":    "https://twitter.com/%s",
		"github":     "https://github.com/%s",
		"reddit":     "https://www.reddit.com/user/%s",
		"tiktok":     "https://www.tiktok.com/@%s",
		"youtube":    "https://www.youtube.com/@%s",
		"twitch":     "https://www.twitch.tv/%s",
		"roblox":     "https://www.roblox.com/user.aspx?username=%s",
		"steam":      "https://steamcommunity.com/id/%s",
		"spotify":    "https://open.spotify.com/user/%s",
		"pinterest":  "https://www.pinterest.com/%s",
		"telegram":   "https://t.me/%s",
		"snapchat":   "https://www.snapchat.com/add/%s",
		"facebook":   "https://www.facebook.com/%s",
		"linkedin":   "https://www.linkedin.com/in/%s",
		"medium":     "https://medium.com/@%s",
		"soundcloud": "https://soundcloud.com/%s",
		"gitlab":     "https://gitlab.com/%s",
		"bitbucket":  "https://bitbucket.org/%s",
		"devto":      "https://dev.to/%s",
		"hackernews": "https://news.ycombinator.com/user?id=%s",
		"keybase":    "https://keybase.io/%s",
		"replit":     "https://replit.com/@%s",
		"codepen":    "https://codepen.io/%s",
		"behance":    "https://www.behance.net/%s",
		"dribbble":   "https://dribbble.com/%s",
		"flickr":     "https://www.flickr.com/people/%s",
		"tumblr":     "https://%s.tumblr.com",
		"vk":         "https://vk.com/%s",
		"weibo":      "https://weibo.com/%s",
		"quora":      "https://www.quora.com/profile/%s",
		"goodreads":  "https://www.goodreads.com/%s",
		"patreon":    "https://www.patreon.com/%s",
		"onlyfans":   "https://onlyfans.com/%s",
		"cashapp":    "https://cash.app/$%s",
		"venmo":      "https://venmo.com/u/%s",
		"gravatar":   "https://gravatar.com/%s",
		"about.me":   "https://about.me/%s",
		"producthunt": "https://www.producthunt.com/@%s",
		"angel.co":   "https://angel.co/u/%s",
	}

	fmt.Println()
	ui.Cyan(fmt.Sprintf("checking %d sites...", len(sites)))

	client := &http.Client{Timeout: 8 * time.Second}
	var wg sync.WaitGroup
	var mu sync.Mutex
	sem := make(chan struct{}, 20)

	results := make(map[string]int) // name -> status code
	for name, tmpl := range sites {
		wg.Add(1)
		sem <- struct{}{}
		go func(name, tmpl string) {
			defer wg.Done()
			defer func() { <-sem }()
			url := fmt.Sprintf(tmpl, u)
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
			resp, err := client.Do(req)
			if err != nil {
				mu.Lock()
				results[name] = 0
				mu.Unlock()
				return
			}
			resp.Body.Close()
			mu.Lock()
			results[name] = resp.StatusCode
			mu.Unlock()
		}(name, tmpl)
	}
	wg.Wait()

	fmt.Println()
	fmt.Println(ui.Green("─── results ───"))
	for name, code := range results {
		url := fmt.Sprintf(sites[name], u)
		switch code {
		case 200:
			fmt.Printf("  %s %-14s %s\n", ui.Green("[FOUND]"), name, url)
		case 404:
			fmt.Printf("  %s %-14s %s\n", ui.Red("[none] "), name, url)
		case 0:
			fmt.Printf("  %s %-14s %s\n", ui.Yellow("[err]  "), name, url)
		default:
			fmt.Printf("  %s %-14s %s\n", ui.Yellow("[?]    "), name, url)
		}
	}
}

// SocialScraper alias — same as UsernameSearch.
func SocialScraper() { UsernameSearch() }

// DoxTracker generates dox-tracking search links.
func DoxTracker() {
	ui.Cyan("enter query (name/email/username/phone):")
	var q string
	fmt.Scanln(&q)
	q = strings.TrimSpace(q)

	fmt.Println()
	fmt.Println(ui.Green("─── dox trackers ───"))
	links := [][2]string{
		{"intelx",     "https://intelx.io/?s=" + q},
		{"dehashed",   "https://dehashed.com/search?query=" + q},
		{"epieos",     "https://epieos.com/?q=" + q},
		{"pipl",       "https://pipl.com/search/?q=" + q},
		{"beenverified", "https://www.beenverified.com/app/optout/search?firstName=&lastName=" + q},
		{"whitepages", "https://www.whitepages.com/name/" + q},
		{"spokeo",     "https://www.spokeo.com/search?q=" + q},
		{"peekyou",    "https://www.peekyou.com/" + q},
		{"truepeoplesearch", "https://www.truepeoplesearch.com/results?name=" + q},
		{"fastpeoplesearch", "https://fastpeoplesearch.com/name/" + q},
	}
	for _, l := range links {
		fmt.Printf("  %-20s %s\n", ui.Cyan(l[0]+":"), l[1])
	}
	fmt.Println()
	fmt.Println(ui.Green("─── google dorks ───"))
	fmt.Printf("  site:pastebin.com \"%s\"\n", q)
	fmt.Printf("  site:ghostbin.com \"%s\"\n", q)
	fmt.Printf("  intext:\"%s\" filetype:pdf\n", q)
	fmt.Printf("  \"%s\" (site:linkedin.com OR site:facebook.com)\n", q)
	fmt.Printf("  \"%s\" \"password\" OR \"email\"\n", q)
}

// DoxCreator helps you build a dossier.
func DoxCreator() {
	ui.Cyan("─── dox dossier builder ───")
	fmt.Println("run each of these tools to gather fields, then compile:")
	fmt.Println()
	fmt.Println("  1. IPGeo          — target IP location")
	fmt.Println("  2. PhoneLookup    — carrier and links")
	fmt.Println("  3. EmailLookup    — registrations + breach")
	fmt.Println("  4. UsernameSearch — cross-platform presence")
	fmt.Println("  5. BreachCheck    — leak data")
	fmt.Println("  6. ImageGeo       — location from photo")
	fmt.Println("  7. MetadataExtract — exif from images")
	fmt.Println()
	ui.Green("compile into a markdown or json file manually.")
}