package snapsys

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type SnapItem struct {
	Name     string
	Version  string
	Revision string
}

func CheckRoot() bool {
	return os.Geteuid() == 0
}

func GetDisabledSnaps() ([]SnapItem, error) {
	cmd := exec.Command("snap", "list", "--all")

	cmd.Env = append(os.Environ(), "LANG=en_US.UTF-8")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("snap command failed: %v", err)
	}

	var candidates []SnapItem

	scanner := bufio.NewScanner(bytes.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line) // Bosluklara gore bol

		if len(fields) >= 3 && strings.Contains(line, "disabled") {
			item := SnapItem{
				Name:     fields[0],
				Version:  fields[1],
				Revision: fields[2],
			}
			candidates = append(candidates, item)
		}
	}

	return candidates, nil
}

func RemoveSnap(item SnapItem) error {
	// Command: snap remove <name> --revision=<rev>
	cmd := exec.Command("snap", "remove", item.Name, fmt.Sprintf("--revision=%s", item.Revision))

	return cmd.Run()
}
