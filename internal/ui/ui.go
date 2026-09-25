package ui

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/term"
)

type Item struct {
	Name    string
	Desc    string
	Handler func()
}

type Category struct {
	Name  string
	Items []Item
}

var (
	categories   []Category
	currentCat   int
	currentPage  int
	itemsPerPage = 9
)

func RegisterCategory(c Category) { categories = append(categories, c) }

func CurrentCategory() *Category {
	if currentCat < 0 || currentCat >= len(categories) {
		return nil
	}
	return &categories[currentCat]
}

func CurrentCategoryName() string {
	c := CurrentCategory()
	if c == nil {
		return "unknown"
	}
	return c.Name
}

func CategoryLeft() {
	currentCat--
	if currentCat < 0 {
		currentCat = len(categories) - 1
	}
	currentPage = 0
	redraw()
}

func CategoryRight() {
	currentCat++
	if currentCat >= len(categories) {
		currentCat = 0
	}
	currentPage = 0
	redraw()
}

func PagePrev() {
	if currentPage > 0 {
		currentPage--
	}
	redraw()
}

func PageNext() {
	c := CurrentCategory()
	if c == nil {
		return
	}
	maxPage := (len(c.Items) - 1) / itemsPerPage
	if currentPage < maxPage {
		currentPage++
	}
	redraw()
}

func redraw() {
	Clear()
	PrintBanner()
	RenderMainMenu()
}

// ---------- ansi helpers ----------

func stripANSI(s string) string {
	var b strings.Builder
	inEsc := false
	for _, r := range s {
		if r == 0x1b {
			inEsc = true
			continue
		}
		if inEsc {
			if r == 'm' {
				inEsc = false
			}
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func visLen(s string) int { return utf8.RuneCountInString(stripANSI(s)) }

func termWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 100
	}
	if w < 70 {
		w = 70
	}
	if w > 130 {
		w = 130
	}
	return w
}

const bannerArt = ` ██████╗ ██╗  ██╗ ██████╗ ███████╗████████╗██╗     ██╗███╗   ██╗███████╗
██╔════╝ ██║  ██║██╔═══██╗██╔════╝╚══██╔══╝██║     ██║████╗  ██║██╔════╝
██║  ███╗███████║██║   ██║███████╗   ██║   ██║     ██║██╔██╗ ██║█████╗
██║   ██║██╔══██║██║   ██║╚════██║   ██║   ██║     ██║██║╚██╗██║██╔══╝
╚██████╔╝██║  ██║╚██████╔╝███████║   ██║   ███████╗██║██║ ╚████║███████╗
 ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚══════╝   ╚═╝   ╚══════╝╚═╝╚═╝  ╚═══╝╚══════╝`

// ---------- startup ----------

func Startup() {
	Clear()
	for i := 0; i < 4; i++ {
		Clear()
		fmt.Println(glitch(bannerArt, i))
		time.Sleep(time.Duration(40+i*15) * time.Millisecond)
	}
	Clear()
	fmt.Println()
	fmt.Println(NeonGreen(bannerArt))
	fmt.Println()
	fmt.Println("            " + DimCyan("advanced multi-tool   ·   v1.0   ·   made by seventycs"))
	fmt.Println()
	time.Sleep(300 * time.Millisecond)
	fmt.Println("            " + Border("──────────────────────────────────────────────────────────────"))
	time.Sleep(1400 * time.Millisecond)
}

func glitch(art string, n int) string {
	out := []rune(art)
	glyphs := []rune{'░', '▒', '▓', '█', '/', '\\', '|', '·', '+', '*'}
	count := 15 + n*30
	for i := 0; i < count; i++ {
		idx := rand.Intn(len(out))
		if out[idx] != '\n' && out[idx] != ' ' {
			out[idx] = glyphs[rand.Intn(len(glyphs))]
		}
	}
	switch n % 3 {
	case 0:
		return Green(string(out))
	case 1:
		return Cyan(string(out))
	default:
		return NeonGreen(string(out))
	}
}

// ---------- header ----------

func PrintBanner() {
	w := termWidth()

	fmt.Println()
	fmt.Println(NeonGreen(bannerArt))
	fmt.Println()

	left := "  " + DimCyan("advanced multi-tool")
	right := DimCyan("v1.0  ·  made by seventycs") + "  "
	gap := w - visLen(left) - visLen(right)
	if gap < 1 {
		gap = 1
	}
	fmt.Println(left + strings.Repeat(" ", gap) + right)
	fmt.Println(DimCyan("  " + strings.Repeat("─", w-4)))
	fmt.Println()
}

// ---------- main menu ----------

func RenderMainMenu() {
	c := CurrentCategory()
	if c == nil {
		return
	}
	w := termWidth()

	// category pills
	fmt.Println("  " + renderCategories())
	fmt.Println()
	fmt.Println(DimCyan("  " + strings.Repeat("┄", w-4)))
	fmt.Println()

	// section title
	maxPage := (len(c.Items) - 1) / itemsPerPage
	if maxPage < 0 {
		maxPage = 0
	}
	title := strings.ToUpper(c.Name)
	page := fmt.Sprintf("%d / %d", currentPage+1, maxPage+1)
	titleLine := "  " + NeonGreen("◈ ") + "\033[1;92m" + title + "\033[0m"
	gap := w - visLen(titleLine) - visLen(page) - 4
	if gap < 1 {
		gap = 1
	}
	fmt.Println(titleLine + strings.Repeat(" ", gap) + DimCyan(page) + "  ")
	fmt.Println()

	// items — 2-line format
	start := currentPage * itemsPerPage
	end := start + itemsPerPage
	if end > len(c.Items) {
		end = len(c.Items)
	}

	for i := start; i < end; i++ {
		item := c.Items[i]
		num := fmt.Sprintf("%02d", i)

		// line 1: "   ##  NAME"
		nameUpper := strings.ToUpper(item.Name)
		line1 := "   " + Cyan(num) + "    " + "\033[1;97m" + nameUpper + "\033[0m"
		fmt.Println(line1)

		// line 2: "       description"
		line2 := "         " + DimCyan(item.Desc)
		fmt.Println(line2)

		// blank spacer between items (except last)
		if i < end-1 {
			fmt.Println()
		}
	}
	fmt.Println()

	// footer
	fmt.Println(DimCyan("  " + strings.Repeat("─", w-4)))
	footer := "  " + Green("[A/D]") + " " + DimCyan("category") + "    " +
		Green("[P/N]") + " " + DimCyan("page") + "    " +
		Green("[#]") + " " + DimCyan("run") + "    " +
		Green("[99]") + " " + DimCyan("exit")
	fmt.Println(footer)
	fmt.Println()
}

func renderCategories() string {
	var parts []string
	for i, c := range categories {
		name := strings.ToUpper(c.Name)
		if i == currentCat {
			parts = append(parts, "\033[1;92m▐ "+name+" ▌\033[0m")
		} else {
			parts = append(parts, DimCyan(name))
		}
	}
	return strings.Join(parts, "   ")
}

// ---------- tool screen ----------

func Screen(title, desc string) {
	w := termWidth()
	fmt.Println()
	fmt.Println("  " + NeonGreen("◈ ") + "\033[1;92m" + strings.ToUpper(title) + "\033[0m")
	fmt.Println("  " + DimCyan(desc))
	fmt.Println(DimCyan("  " + strings.Repeat("─", w-4)))
	fmt.Println()
}

func Clear() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[2J\033[H")
	}
}

// ---------- colors ----------

func Green(s string) string     { return "\033[92m" + s + "\033[0m" }
func NeonGreen(s string) string { return "\033[38;5;46m" + s + "\033[0m" }
func Cyan(s string) string      { return "\033[96m" + s + "\033[0m" }
func DimCyan(s string) string   { return "\033[38;5;65m" + s + "\033[0m" }
func White(s string) string     { return "\033[97m" + s + "\033[0m" }
func Red(s string) string       { return "\033[91m" + s + "\033[0m" }
func Yellow(s string) string    { return "\033[93m" + s + "\033[0m" }
func Border(s string) string    { return "\033[38;5;29m" + s + "\033[0m" }

// AnimatedBanner kept for main.go compat.
func AnimatedBanner() { PrintBanner() }