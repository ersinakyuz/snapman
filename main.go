package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ersinakyuz/SnapMan/snapsys"
)

func main() {
	fmt.Println("starting...")

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

	fmt.Printf("\n%d disabled snaps found:\n", len(items))
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "Package\tVersion\tRevision")
	fmt.Fprintln(w, "-----\t--------\t--------")
	for _, item := range items {
		fmt.Fprintf(w, "%s\t%s\t%s\n", item.Name, item.Version, item.Revision)
	}
	w.Flush()
	fmt.Println("")

	fmt.Println("\nfinished.")
}
