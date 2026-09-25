package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"ghosttools/internal/ui"
)

var reader = bufio.NewReader(os.Stdin)

func Run() {
	ui.Clear()
	ui.PrintBanner()
	ui.RenderMainMenu()

	for {
		fmt.Printf("\n%s@%s:~# ", ui.Green("ghosttools"), ui.Cyan(ui.CurrentCategoryName()))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "":
			// empty enter — redraw and continue
			ui.Clear()
			ui.PrintBanner()
			ui.RenderMainMenu()
		case "a", "A":
			ui.CategoryLeft()
		case "d", "D":
			ui.CategoryRight()
		case "p", "P":
			ui.PagePrev()
			ui.Clear()
			ui.PrintBanner()
			ui.RenderMainMenu()
		case "n", "N":
			ui.PageNext()
			ui.Clear()
			ui.PrintBanner()
			ui.RenderMainMenu()
		case "99":
			fmt.Println()
			ui.Cyan("closing session. stay ghost.")
			os.Exit(0)
		default:
			handleCategoryChoice(input)
			ui.Clear()
			ui.PrintBanner()
			ui.RenderMainMenu()
		}
	}
}

func handleCategoryChoice(choice string) {
	cat := ui.CurrentCategory()
	if cat == nil {
		ui.Red("invalid category")
		return
	}
	idx, err := parseInt(choice)
	if err != nil || idx < 0 || idx >= len(cat.Items) {
		ui.Red("unknown command: " + choice)
		pause()
		return
	}

	item := cat.Items[idx]

		// enter tool screen
	ui.Clear()
	ui.PrintBanner()
	ui.Screen(item.Name, item.Desc)

	// run the tool
	item.Handler()

	// done prompt — obvious and pushed to bottom
	pause()
}

func pause() {
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println(ui.Green("  ═══════════════════════════════════════════════════════════════"))
	fmt.Println(ui.Green("     ✓  done  ·  press ENTER to return to menu"))
	fmt.Println(ui.Green("  ═══════════════════════════════════════════════════════════════"))
	reader.ReadString('\n')
}

func parseInt(s string) (int, error) {
	n := 0
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}