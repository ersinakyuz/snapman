package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
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
	dryRun := flag.Bool("dry-run", false, "List disabled revisions without removing them")
	var assumeYes bool
	flag.BoolVar(&assumeYes, "assume-yes", false, "Proceed without confirmation prompt")
	flag.BoolVar(&assumeYes, "yes", false, "Alias for --assume-yes")
	flag.Parse()

	fmt.Println(ColorBold + ColorBlue + "SnapMan starting..." + ColorReset)

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

	if *dryRun {
		fmt.Println("Dry run enabled: no removals executed.")
		fmt.Println("\nOperation Complete.")
		return
	}

	if !snapsys.CheckRoot() {
		fmt.Println("This tool requires root privileges. Please run with 'sudo'.")
		os.Exit(1)
	}

	if !confirmProceed(len(items), assumeYes) {
		fmt.Println("Aborted. No changes made.")
		return
	}

	var removedCount, skippedMissing int
	var reclaimedBytes int64

	for _, item := range items {
		if !item.FileFound {
			skippedMissing++
			fmt.Printf("Skipping %s (revision %s): file missing on disk.\n", item.Name, item.Revision)
			continue
		}

		if err := snapsys.RemoveSnap(item); err != nil {
			fmt.Printf("Failed to remove %s (revision %s): %v\n", item.Name, item.Revision, err)
			continue
		}

		removedCount++
		reclaimedBytes += item.SizeBytes
		fmt.Printf("Removed %s (revision %s) [%s]\n", item.Name, item.Revision, item.SizeHuman)
	}

	fmt.Printf("\nSummary: removed %d, skipped (missing) %d, reclaimed %s.\n",
		removedCount, skippedMissing, formatTotal(reclaimedBytes))

	fmt.Println("\nOperation Complete.")
}

func confirmProceed(count int, assumeYes bool) bool {
	if assumeYes {
		return true
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Proceed to remove %d revisions? [y/N]: ", count)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	return response == "y" || response == "yes"
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
