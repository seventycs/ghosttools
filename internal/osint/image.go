package osint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"ghosttools/internal/ui"
)

// ImageGeo — extracts EXIF, calls picarta.ai if PICARTA_API_KEY is set,
// then prints reverse image search URLs. Location from photo.
func ImageGeo() {
	ui.Cyan("enter path to image:")
	var path string
	fmt.Scanln(&path)
	path = strings.Trim(strings.TrimSpace(path), "\"")

	if _, err := os.Stat(path); err != nil {
		ui.Red("file not found")
		return
	}

	// 1. EXIF extraction
	fmt.Println()
	fmt.Println(ui.Green("─── exif data ───"))
	if _, err := exec.LookPath("exiftool"); err == nil {
		out, _ := exec.Command("exiftool", path).Output()
		for _, line := range strings.Split(string(out), "\n") {
			lower := strings.ToLower(line)
			if strings.Contains(lower, "gps") || strings.Contains(lower, "date") ||
				strings.Contains(lower, "camera") || strings.Contains(lower, "model") ||
				strings.Contains(lower, "location") || strings.Contains(lower, "software") {
				fmt.Println("  " + line)
			}
		}
	} else {
		ui.Yellow("exiftool not installed (winget install exiftool) — skipping raw exif")
	}

	// 2. picarta AI geolocation
	fmt.Println()
	fmt.Println(ui.Green("─── ai geolocation (picarta) ───"))
	apiKey := os.Getenv("PICARTA_API_KEY")
	if apiKey == "" {
		ui.Yellow("set PICARTA_API_KEY env var to enable AI geolocation")
		ui.Yellow("  register free at picarta.ai")
	} else {
		f, err := os.Open(path)
		if err != nil {
			ui.Red(err.Error())
			return
		}
		defer f.Close()
		var buf bytes.Buffer
		w := multipart.NewWriter(&buf)
		part, _ := w.CreateFormFile("image", filepath.Base(path))
		io.Copy(part, f)
		w.Close()

		req, _ := http.NewRequest("POST", "https://picarta.ai/ip/api/v1/photo", &buf)
		req.Header.Set("Content-Type", w.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+apiKey)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			ui.Red("picarta request failed: " + err.Error())
		} else {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			var d map[string]interface{}
			json.Unmarshal(body, &d)
			out, _ := json.MarshalIndent(d, "", "  ")
			fmt.Println(ui.Green(string(out)))
			if lat, ok := d["latitude"]; ok {
				if lon, ok := d["longitude"]; ok {
					fmt.Printf("  %s https://www.google.com/maps?q=%v,%v\n", ui.Cyan("map:"), lat, lon)
				}
			}
		}
	}

	// 3. reverse image search links
	fmt.Println()
	fmt.Println(ui.Green("─── reverse image search ───"))
	fmt.Println("  upload the image to:")
	fmt.Println("    google lens: https://lens.google.com/")
	fmt.Println("    yandex:      https://yandex.com/images/")
	fmt.Println("    tineye:      https://tineye.com/")
	fmt.Println("    bing visual: https://www.bing.com/visualsearch")
	fmt.Println("    pimeyes:     https://pimeyes.com/en  (face search)")
}

// ReverseImage alias — same as ImageGeo.
func ReverseImage() { ImageGeo() }

// MetadataExtract dumps full exiftool output.
func MetadataExtract() {
	ui.Cyan("enter file path:")
	var path string
	fmt.Scanln(&path)
	path = strings.Trim(strings.TrimSpace(path), "\"")
	if _, err := os.Stat(path); err != nil {
		ui.Red("file not found")
		return
	}
	if _, err := exec.LookPath("exiftool"); err != nil {
		ui.Yellow("exiftool not installed (winget install exiftool)")
		return
	}
	out, _ := exec.Command("exiftool", path).Output()
	fmt.Println(string(out))
}