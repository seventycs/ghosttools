package osint

import (
	"fmt"
	"strings"

	"ghosttools/internal/ui"
)

// PhoneLookup aggregates phone OSINT across public directories.
func PhoneLookup() {
	ui.Cyan("enter phone number with country code (+15551234567):")
	var p string
	fmt.Scanln(&p)
	p = strings.TrimSpace(p)
	clean := strings.ReplaceAll(p, " ", "")
	clean = strings.ReplaceAll(clean, "-", "")
	clean = strings.ReplaceAll(clean, "(", "")
	clean = strings.ReplaceAll(clean, ")", "")

	// guess country from prefix
	country := "unknown"
	prefixes := map[string]string{
		"+1":  "USA / Canada",
		"+44": "UK",
		"+33": "France",
		"+49": "Germany",
		"+34": "Spain",
		"+39": "Italy",
		"+61": "Australia",
		"+81": "Japan",
		"+86": "China",
		"+91": "India",
		"+7":  "Russia",
		"+55": "Brazil",
		"+52": "Mexico",
		"+31": "Netherlands",
		"+46": "Sweden",
		"+47": "Norway",
		"+45": "Denmark",
		"+48": "Poland",
		"+90": "Turkey",
		"+20": "Egypt",
		"+27": "South Africa",
		"+971": "UAE",
		"+972": "Israel",
		"+966": "Saudi Arabia",
	}
	for pre, name := range prefixes {
		if strings.HasPrefix(clean, pre) {
			country = name
			break
		}
	}

	fmt.Println()
	fmt.Println(ui.Green("─── phone intel ───"))
	fmt.Printf("  %-14s %s\n", "number:", p)
	fmt.Printf("  %-14s %s\n", "normalized:", clean)
	fmt.Printf("  %-14s %s\n", "country:", country)
	fmt.Printf("  %-14s %s\n", "carrier:", "run external lookup (see links)")
	fmt.Println()
	fmt.Println(ui.Cyan("─── osint links ───"))
	links := [][2]string{
		{"truecaller", "https://www.truecaller.com/search/" + strings.TrimPrefix(clean, "+")},
		{"sync.me",    "https://sync.me/search/?number=" + clean},
		{"spamcalls",  "https://spamcalls.net/en/search?q=" + clean},
		{"whocalld",   "https://whocalld.com/" + clean},
		{"numlookup",  "https://www.numlookup.com/?number=" + clean},
		{"callapp",    "https://callapp.com/search/" + clean},
		{"fraudrecord","https://www.fraudrecord.com/search.asp?phone=" + clean},
		{"reverse-gen", "https://www.reversephonelookup.com/?n=" + clean},
	}
	for _, l := range links {
		fmt.Printf("  %-14s %s\n", ui.Cyan(l[0]+":"), l[1])
	}

	// optional numlookup if API key set (environment)
	ui.Cyan("\nfor live carrier data, register at numlookup.com and set NUM_LOOKUP_KEY env var.")
}