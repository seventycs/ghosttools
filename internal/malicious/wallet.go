package malicious

import (
    "crypto/sha256"
    "fmt"
    "strings"

    "ghosttools/internal/ui"
)

func WalletBrute() {
    ui.Cyan("enter mnemonic phrase (space separated):")
    var phrase string
    fmt.Scanln(&phrase)
    words := strings.Fields(phrase)
    ui.Cyan(fmt.Sprintf("checking %d words against bip39...", len(words)))
    // simplified — real check requires wordlist
    h := sha256.Sum256([]byte(phrase))
    fmt.Printf("  seed hash: %x\n", h[:8])
    ui.Yellow("full bip39 validation requires full wordlist — dropping wordlist into ./wordlist.txt")
}