package roblox

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"ghosttools/internal/ui"
)

func UserInfo() {
	ui.Cyan("enter roblox username:")
	var u string
	fmt.Scanln(&u)
	body := `{"usernames":["` + u + `"],"excludeBannedUsers":false}`
	r, err := http.Post("https://users.roblox.com/v1/usernames/users", "application/json", strings.NewReader(body))
	if err != nil {
		ui.Red(err.Error())
		return
	}
	defer r.Body.Close()
	var out map[string]interface{}
	json.NewDecoder(r.Body).Decode(&out)
	b, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(b))
}

func CookieInfo() {
	ui.Cyan("enter .ROBLOSECURITY cookie:")
	var c string
	fmt.Scanln(&c)
	req, _ := http.NewRequest("GET", "https://users.roblox.com/v1/users/authenticated", nil)
	req.Header.Set("Cookie", ".ROBLOSECURITY="+c)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ui.Red(err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		ui.Red(fmt.Sprintf("invalid cookie (status %d)", resp.StatusCode))
		return
	}
	var out map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&out)
	b, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(ui.Green(string(b)))
}

func CookieLogin() {
	ui.Yellow("validation only — no session is stored")
	CookieInfo()
}

func GroupInfo() {
	ui.Cyan("enter group id:")
	var g string
	fmt.Scanln(&g)
	resp, err := http.Get("https://groups.roblox.com/v1/groups/" + g)
	if err != nil {
		ui.Red(err.Error())
		return
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
}

func AssetDownload() {
	ui.Cyan("enter asset id:")
	var id string
	fmt.Scanln(&id)
	ui.Cyan("enter .ROBLOSECURITY (blank to skip):")
	var c string
	fmt.Scanln(&c)
	url := "https://assetdelivery.roblox.com/v1/asset/?id=" + id
	req, _ := http.NewRequest("GET", url, nil)
	if c != "" {
		req.Header.Set("Cookie", ".ROBLOSECURITY="+c)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ui.Red(err.Error())
		return
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	out := "asset_" + id + ".rbxm"
	os.WriteFile(out, data, 0644)
	ui.Green("saved to " + out)
}

func NameHistory() {
	ui.Cyan("enter user id:")
	var id string
	fmt.Scanln(&id)
	resp, err := http.Get("https://users.roblox.com/v1/users/" + id + "/username-history?limit=100")
	if err != nil {
		ui.Red(err.Error())
		return
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
}

func UsernameChecker() {
	ui.Cyan("enter username:")
	var u string
	fmt.Scanln(&u)
	body := `{"usernames":["` + u + `"],"excludeBannedUsers":false}`
	r, _ := http.Post("https://users.roblox.com/v1/usernames/users", "application/json", strings.NewReader(body))
	if r == nil {
		ui.Red("failed")
		return
	}
	defer r.Body.Close()
	var out map[string]interface{}
	json.NewDecoder(r.Body).Decode(&out)
	data, _ := out["data"].([]interface{})
	if len(data) == 0 {
		ui.Green("available")
	} else {
		ui.Red("taken")
	}
}

func CookieRefresher() {
	ui.Yellow("cookie refresh requires valid authenticated session")
}

func GameInfo() {
	ui.Cyan("enter universe id:")
	var id string
	fmt.Scanln(&id)
	resp, err := http.Get("https://games.roblox.com/v1/games?universeIds=" + id)
	if err != nil {
		ui.Red(err.Error())
		return
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
}

func InventoryDump() {
	ui.Cyan("enter user id:")
	var id string
	fmt.Scanln(&id)
	resp, err := http.Get("https://inventory.roblox.com/v2/users/" + id + "/inventory?assetTypes=Hat&limit=100&sortOrder=Asc")
	if err != nil {
		ui.Red(err.Error())
		return
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
}

func TradeScanner() {
	ui.Cyan("enter user id:")
	var id string
	fmt.Scanln(&id)
	resp, err := http.Get("https://trades.roblox.com/v1/trades/" + id)
	if err != nil {
		ui.Red(err.Error())
		return
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
}

func LimitedSniper() {
	ui.Cyan("enter limited item id:")
	var id string
	fmt.Scanln(&id)
	resp, err := http.Get("https://economy.roblox.com/v2/assets/" + id + "/details")
	if err != nil {
		ui.Red(err.Error())
		return
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
}

func GroupCloner() {
	ui.Yellow("group cloning requires authenticated session and api scope — not supported")
}

func BotFollower() {
	ui.Yellow("bot following requires valid accounts with cookies")
}