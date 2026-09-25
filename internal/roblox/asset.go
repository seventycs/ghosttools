package roblox

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"ghosttools/internal/ui"
)

// AssetFetch downloads a roblox asset with optional auth cookie.
func AssetFetch(assetID, cookie string) {
	url := "https://assetdelivery.roblox.com/v1/asset/?id=" + assetID
	req, _ := http.NewRequest("GET", url, nil)
	if cookie != "" {
		req.Header.Set("Cookie", ".ROBLOSECURITY="+cookie)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ui.Red(err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		ui.Red(fmt.Sprintf("failed: status %d", resp.StatusCode))
		return
	}
	data, _ := io.ReadAll(resp.Body)
	out := "asset_" + assetID + ".rbxm"
	os.WriteFile(out, data, 0644)
	ui.Green("downloaded: " + out + fmt.Sprintf(" (%d bytes)", len(data)))
}