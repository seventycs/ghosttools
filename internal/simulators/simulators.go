package simulators

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"ghosttools/internal/ui"
)

func FakeTokenGen() {
	randHex := func(n int) string {
		b := make([]byte, n)
		rand.Read(b)
		return hex.EncodeToString(b)
	}
	token := base64.RawURLEncoding.EncodeToString([]byte(randHex(16))) + "." +
		base64.RawURLEncoding.EncodeToString([]byte(randHex(6))) + "." +
		base64.RawURLEncoding.EncodeToString([]byte(randHex(20)))
	fmt.Println(ui.Green(token))
}

func FakeMailGen() {
	domains := []string{"gmail.com", "outlook.com", "yahoo.com", "proton.me", "icloud.com"}
	names := []string{"james", "sarah", "michael", "emma", "david", "olivia", "john", "ava"}
	for i := 0; i < 5; i++ {
		n := names[int(time.Now().UnixNano()+int64(i))%len(names)]
		d := domains[i%len(domains)]
		fmt.Printf("  %s%d@%s\n", n, time.Now().UnixNano()%1000, d)
	}
}

func FakeIdentity() {
	first := []string{"John", "Sarah", "Michael", "Emma", "David", "Olivia"}
	last := []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia"}
	streets := []string{"Maple St", "Oak Ave", "Pine Rd", "Cedar Ln", "Elm Blvd"}
	cities := []string{"Springfield", "Riverside", "Franklin", "Clinton", "Georgetown"}
	states := []string{"CA", "NY", "TX", "FL", "IL"}
	seed := time.Now().UnixNano()
	pick := func(l []string) string { return l[int(seed)%len(l)] }
	fmt.Printf("  name:    %s %s\n", pick(first), pick(last))
	fmt.Printf("  address: %d %s\n", seed%9000+1000, pick(streets))
	fmt.Printf("  city:    %s, %s %05d\n", pick(cities), pick(states), seed%99999)
	fmt.Printf("  phone:   (%03d) %03d-%04d\n", seed%900+100, seed%900+100, seed%9000+1000)
	fmt.Printf("  dob:     %d/%d/%d\n", seed%12+1, seed%28+1, 1960+seed%50)
}

func FakeCreditCard() {
	prefix := "4"
	body := fmt.Sprintf("%014d", time.Now().UnixNano()%100000000000000)
	num := prefix + body
	sum := 0
	alt := false
	for i := len(num) - 1; i >= 0; i-- {
		d := int(num[i] - '0')
		if alt {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		alt = !alt
	}
	check := (10 - (sum % 10)) % 10
	fmt.Printf("  %s%d\n", num, check)
	fmt.Printf("  exp: %02d/%d\n", time.Now().Month(), time.Now().Year()+2)
	fmt.Printf("  cvv: %03d\n", time.Now().UnixNano()%900+100)
}

func SocialBotter() {
	ui.Cyan("enter target username:")
	var u string
	fmt.Scanln(&u)
	fmt.Println("  simulating follower growth...")
	for i := 0; i < 10; i++ {
		fmt.Printf("  +%d followers\n", (i+1)*100)
		time.Sleep(300 * time.Millisecond)
	}
	ui.Green("simulated")
}

func FakePaypalOTP() {
	fmt.Println(ui.Cyan("paypal-otp simulation — for phish templates only"))
	fmt.Println("  code: " + fmt.Sprintf("%06d", time.Now().UnixNano()%1000000))
}

func FakeFortnite() {
	fmt.Println(ui.Cyan("fortnite account simulation"))
	fmt.Println("  vbucks:  " + fmt.Sprintf("%d", time.Now().UnixNano()%99999))
	fmt.Println("  level:   " + fmt.Sprintf("%d", time.Now().UnixNano()%500))
}

func FakeExodus() {
	fmt.Println(ui.Cyan("exodus wallet simulation"))
	b := make([]byte, 32)
	rand.Read(b)
	fmt.Println("  seed: " + hex.EncodeToString(b))
}

func HackerTerminal() {
	lines := []string{
		"initializing exploit framework...",
		"scanning target network...",
		"found 47 open ports",
		"bypassing firewall...",
		"escalating privileges...",
		"access granted",
		"downloading /etc/shadow...",
		"cracking hashes...",
		"root access acquired",
	}
	for _, l := range lines {
		fmt.Println(ui.Green("[+] ") + l)
		time.Sleep(500 * time.Millisecond)
	}
}

func FakeBruteforcer() {
	fmt.Println(ui.Cyan("bruteforce simulation"))
	for i := 0; i < 20; i++ {
		fmt.Printf("  trying %s...\n", randString(12))
		time.Sleep(100 * time.Millisecond)
	}
	ui.Green("password found: hunter2")
}

func QRCodeGen() {
	ui.Cyan("enter text/url:")
	var t string
	fmt.Scanln(&t)
	fmt.Println(ui.Yellow("qr generation via online api:"))
	fmt.Printf("  https://api.qrserver.com/v1/create-qr-code/?data=%s\n", t)
}

func FakeDiscord() {
	fmt.Println(ui.Cyan("discord login page template"))
	fmt.Println("  https://discord.com/login")
}

func FakeSteam() {
	fmt.Println(ui.Cyan("steam login page template"))
	fmt.Println("  https://store.steampowered.com/login")
}

func FakeInstagram() {
	fmt.Println(ui.Cyan("instagram login page template"))
	fmt.Println("  https://www.instagram.com/accounts/login")
}

func randString(n int) string {
	const l = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	rand.Read(b)
	for i := range b {
		b[i] = l[int(b[i])%len(l)]
	}
	return string(b)
}

var _ = big.NewInt
var _ = exec.Command
var _ = runtime.GOOS
var _ = strings.Join