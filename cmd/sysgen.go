package cmd

import (
	"ghosttools/internal/sysgen"
	"ghosttools/internal/ui"
)

func init() {
	ui.RegisterCategory(ui.Category{
		Name: "SYS/GEN",
		Items: []ui.Item{
			{"Base64 Tool",     "encode or decode strings",                sysgen.Base64Codec},
			{"System Info",      "inspect local hardware specs",            sysgen.SystemInfo},
			{"IP Pinger",        "ping target IPs for latency",             sysgen.IPPinger},
			{"Python Obfuscator",       "obfuscate python to prevent reverse",     sysgen.Obfuscator},
			{"Metadata Scan",    "analyze and remove EXIF data",            sysgen.MetadataScan},
			{"About",         "display version, license, developers",    sysgen.AppInfo},
			{"Config",       "manage themes, updates, startup",         sysgen.AppConfig},
			{"Disable AV",       "add defender exclusions",                 sysgen.DisableAV},
			{"Windows Debloat","launch utility to debloat and optimize",  sysgen.Debloater},
			{"Proxy Pull",    "scrape HTTP, SOCKS4, SOCKS5",             sysgen.ProxyScraper},
			{"Proxy Check",    "test scraped proxies for validity",       sysgen.ProxyChecker},
			{"Registry Edit",  "safe registry tweak manager",             sysgen.RegistryEditor},
			{"Service Manager",  "list and control windows services",       sysgen.ServiceManager},
			{"Startup Edit",  "manage startup programs",                 sysgen.StartupManager},
			{"File Wipe",    "securely delete files",                   sysgen.FileShredder},
			{"Password Gen",     "cryptographically secure passwords",      sysgen.PasswordGen},
			{"Hash Crack",     "crack md5/sha1/sha256 hashes",            sysgen.HashCracker},
			{"Port Forward",     "set up local port forwarding",            sysgen.PortForward},
			{"Packet Sniff",   "capture local network packets",           sysgen.PacketSniffer},
			{"WiFi Scan",     "scan nearby wifi networks",               sysgen.WiFiScanner},
			{"Wipe Traces",      "delete all ghosttools traces from this pc", ui.WipeLocal},
		},
	})
}