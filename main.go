package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ersinakyuz/SnapMan/internal/snapsys"
)

// ANSI Color Codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[37m"
	ColorBold   = "\033[1m"
)

func main() {
	fmt.Println(ColorBold + ColorBlue + "SnapMan starting..." + ColorReset)

	if !snapsys.CheckRoot() {
		fmt.Println("This tool requires root privileges. Please run with 'sudo'.")
		os.Exit(1)
	}

	fmt.Print("\nSystem scanning...")

	items, err := snapsys.GetDisabledSnaps()
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		os.Exit(1)
	}

	if len(items) == 0 {
		fmt.Println("✓ System is clean! No old snaps to remove.")
		return
	}

	// Calculate total reclaimable space
	var totalReclaim int64 = 0
	for _, item := range items {
		totalReclaim += item.SizeBytes
	}

	fmt.Printf("\nFound %s%d%s disabled revisions. Potential gain: %s%s%s\n\n",
		ColorBold, len(items), ColorReset,
		ColorBold+ColorYellow, formatTotal(totalReclaim), ColorReset)

	// Display header
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	// Table Header (Bold)
	fmt.Fprintln(w, "PACKAGE\tVERSION\tREV\tSIZE\tSTATUS")
	fmt.Fprintln(w, "-------\t-------\t---\t----\t------")

	for _, item := range items {
		statusStr := "Ready"
		if !item.FileFound {
			statusStr = "Missing (Ghost)"
		}

		// Renk kodları olmadığı için \t karakterleri görevini tam yapacak
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			item.Name,
			truncateString(item.Version, 25),
			item.Revision,
			item.SizeHuman,
			statusStr,
		)
	}

	w.Flush()
	fmt.Println("")

	fmt.Println("\nOperation Complete.")
}
func formatTotal(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func truncateString(str string, num int) string {
	if len(str) > num {
		return str[0:num-3] + "..."
	}
	return str
}
