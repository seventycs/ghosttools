package cmd

import (
	"ghosttools/internal/osint"
	"ghosttools/internal/ui"
)

func init() {
	ui.RegisterCategory(ui.Category{
		Name: "OSINT",
		Items: []ui.Item{
			{"Port Scan",       "scan target hosts for open ports", osint.PortScanner},
			{"Whois",       "domain registration details",       osint.WhoisLookup},
			{"DNS Lookup",         "A, MX, TXT, NS, CNAME records",     osint.DNSLookup},
			{"Dox Lookup",        "lookup dox information database",   osint.DoxTracker},
			{"Dox Builder",        "create custom doxing profiles",     osint.DoxCreator},
			{"Phone Lookup",       "carrier and location of phone",     osint.PhoneLookup},
			{"Email Lookup",       "OSINT data associated with email",  osint.EmailLookup},
			{"IP Locator",     "geolocate any IP + port scan",      osint.IPGeo},
			{"Image Locator",  "geospy-style AI location from photo", osint.ImageGeo},
			{"Username Search",    "search username across 40+ sites",  osint.UsernameSearch},
			{"Breach Lookup",       "check email against breach databases", osint.BreachCheck},
			{"Social Scraper",     "pull public social media data",     osint.SocialScraper},
			{"Domain Check",  "threat intel on domains",           osint.DomainRep},
			{"SSL Lookup",      "inspect SSL certificates (crt.sh)", osint.SSLCert},
			{"Reverse Image",      "reverse image search engine",       osint.ReverseImage},
			{"EXIF Reader",   "pull EXIF from images",             osint.MetadataExtract},
			{"Subdomain Scan",     "enumerate subdomains (crt.sh)",     osint.SubdomainEnum},
			{"Shodan Search",      "search shodan + censys + zoomeye",  osint.ShodanLookup},
			{"Wayback",    "historic snapshots of URLs",        osint.Wayback},
			{"Onion Search",    "search .onion indexes (tor)",       osint.DarkWeb},
		},
	})
}