<div align="center">

# 👻 ghosttools

**advanced multi-tool framework — discord · osint · malicious · roblox · sys/gen · simulators**

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Platform](https://img.shields.io/badge/Platform-Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)](https://www.microsoft.com/windows)
[![License](https://img.shields.io/badge/License-MIT-00ffcc?style=for-the-badge)](LICENSE)
[![Version](https://img.shields.io/badge/version-1.0.0-00ffcc?style=for-the-badge)]()

*toolkit for recon.*

*discord: https://discord.gg/GustAS9SGY*

</div>

---

## ⚡ features

| category | tools |
|----------|-------|
| **discord** | webhook spam/delete · token tools · account nuker · server info · bot invite gen · selfbot (40+ commands) · server cloner · nuke bot · username checker · token checker · guild backup · message logger · voice spam |
| **osint** | ip geo + port scan · reverse dns · **ai image geolocation** (picarta) · exif extract · 40-site username search · email registration check · phone carrier lookup · dns (a/mx/txt/ns/cname) · whois · subdomain enum · ssl cert transparency · shodan/censys search · wayback machine · breach check · dark web indexes · google dorks |
| **malicious** | email bomber · crypto clipper · vuln scanner · stress tester · browser stealer · keylogger · rat builder · wallet brute · reverse shell (multiple formats) · phishing pages · persistence · av evasion tips · c2 server · botnet stub · ransom sim |
| **roblox** | user info · cookie validation · cookie login · group info · asset download · name history · username checker · cookie refresher · game info · inventory dump · trade scanner · limited sniper · group cloner · bot follower |
| **sys/gen** | base64 codec · system info · ip pinger · python obfuscator · metadata scan · app info · config · defender exclusion · **windows debloater** · proxy scraper · proxy checker · registry editor · service manager · startup manager · file shredder · password gen · hash cracker · port forward · packet sniffer · wifi scanner |
| **simulators** | fake token gen · fake mail · fake identity · fake credit card · social botter · fake paypal otp · fake fortnite · fake exodus · hacker terminal · fake bruteforcer · qr code gen · fake discord/steam/instagram logins |

---

## 🚀 install

### requirements

- **go** 1.22 or newer → [go.dev/dl](https://go.dev/dl)
- **python** 3.10+ → [python.org](https://python.org)
- **windows** 10/11 (linux/macos partially supported)

### build

```bash
git clone https://github.com/YOUR_USERNAME/ghosttools.git
cd ghosttools
go mod tidy
go build -o ghosttools.exe .
./ghosttools.exe
```

### python dependencies (for selfbot + advanced collectors)

```bash
pip install requests pillow discord.py-self opencv-python sounddevice scipy pycaw comtypes pywin32 pycryptodome
```

---

## 🎮 usage

```
a / d   switch category
p / n   previous / next page
[nn]    run the tool at that index
99      exit
```

every tool prints clear prompts. no config needed for most — api keys are optional (picarta, shodan) and set via env vars.

---

## 🏗 structure

```
ghosttools/
├── main.go                     entry point
├── cmd/                        category registrations
│   ├── root.go                 main loop
│   ├── discord.go
│   ├── osint.go
│   ├── malicious.go
│   ├── roblox.go
│   ├── sysgen.go
│   └── simulators.go
├── internal/
│   ├── ui/                     startup, banner, menu
│   ├── discord/                tools + selfbot.py
│   ├── osint/                  all recon modules
│   ├── malicious/              payloads, builders, servers
│   ├── roblox/                 roblox api wrappers
│   ├── sysgen/                 system utilities
│   ├── simulators/             fake generators
│   └── utils/                  crypto + net helpers
└── scripts/                    build + install scripts
```

---

## 🔧 configuration

optional — set env vars to unlock extra features:

| var | what it does |
|-----|--------------|
| `PICARTA_API_KEY` | ai image geolocation (free tier at picarta.ai) |
| `SHODAN_API_KEY` | live shodan lookups (paid) |
| `NUM_LOOKUP_KEY` | live phone carrier lookups |

---

## ⚠️ disclaimer

some tools included in this project interact with systems, accounts, and networks that you may not own or have explicit permission to access. using them against targets without **written, informed consent** is illegal in most jurisdictions and can result in criminal prosecution.

this project is published **for educational and research purposes only**. the author (**seventycs**) does not condone, encourage, or support any illegal use of the software. you are solely responsible for how you use it and for ensuring your actions comply with all applicable local, national, and international laws.

**by using ghosttools, you agree that:**
- you will only run these tools on systems you own or have permission to test
- you will not use them to harm, defraud, or surveil anyone
- the author is not liable for any damages, legal consequences, or misuse caused by this software

if you're unsure whether something is legal in your area — **don't run it.**


---

## 📜 license

mit — do what you want.

---

<div align="center">

**made by [seventycs](https://github.com/YOUR_USERNAME)**

*for every hacker that needs a good toolkit.*

</div>
