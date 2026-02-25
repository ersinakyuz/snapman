package snapsys

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// SnapItem holds the metadata and physical status of a snap revision
type SnapItem struct {
	Name      string
	Version   string
	Revision  string
	SizeBytes int64
	SizeHuman string
	FileFound bool
	FilePath  string // The actual path found on the filesystem
}

var (
	snapNamePattern     = regexp.MustCompile(`^[a-z0-9-]+$`)
	snapRevisionPattern = regexp.MustCompile(`^[0-9]+$`)
)

// CheckRoot checks if the user has root privileges (required for snap operations)
func CheckRoot() bool {
	return os.Geteuid() == 0
}

// formatSize converts raw bytes into a human-readable string (MB, GB, etc.)
func formatSize(bytes int64) string {
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

func getSnapBinaryPath() (string, error) {
	// Use fixed, trusted locations instead of resolving "snap" from PATH.
	candidates := []string{"/usr/bin/snap", "/snap/bin/snap"}
	for _, p := range candidates {
		info, err := os.Stat(p)
		if err != nil || info.IsDir() {
			continue
		}
		if info.Mode()&0111 != 0 {
			return p, nil
		}
	}

	return "", fmt.Errorf("snap binary not found in trusted locations")
}

func isValidSnapID(name, revision string) bool {
	return snapNamePattern.MatchString(name) && snapRevisionPattern.MatchString(revision)
}

// findSnapFile attempts to locate the actual .snap file using Pattern Matching (Globbing)
// instead of relying on a hardcoded naming convention.
func findSnapFile(name, revision string) (string, int64, bool) {
	baseDir := "/var/lib/snapd/snaps"

	// METHOD 1: Try the standard naming convention first (Fastest)
	// Standard format: packagename_revision.snap
	standardName := fmt.Sprintf("%s_%s.snap", name, revision)
	standardPath := filepath.Join(baseDir, standardName)

	info, err := os.Stat(standardPath)
	if err == nil && !info.IsDir() {
		return standardPath, info.Size(), true
	}

	// METHOD 2: If file not found, use 'Glob' search (Safe/Fallback method)
	// We look for any file ending with "_{revision}.snap".
	// This handles cases where 'snap list' truncates the name or uses a different versioning scheme.
	pattern := fmt.Sprintf("%s/*_%s.snap", baseDir, revision)
	matches, err := filepath.Glob(pattern)

	if err != nil || len(matches) == 0 {
		return "", 0, false
	}

	// If there are multiple matches (unlikely but possible), verify the prefix
	// to ensure we picked the file belonging to the correct package.
	for _, match := range matches {
		filename := filepath.Base(match)
		// Does the filename start with the package name?
		if strings.HasPrefix(filename, name) {
			info, err := os.Stat(match)
			if err == nil {
				return match, info.Size(), true
			}
		}
	}

	// File really doesn't exist (Ghost snap)
	return "", 0, false
}

// GetDisabledSnaps runs 'snap list', parses the output, and performs filesystem verification
func GetDisabledSnaps() ([]SnapItem, error) {
	snapPath, err := getSnapBinaryPath()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, snapPath, "list", "--all")

	// Force locale to English to ensure parsing consistency (avoids "disabled" translation issues)
	cmd.Env = append(os.Environ(), "LANG=en_US.UTF-8")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute snap command: %v", err)
	}

	var candidates []SnapItem
	scanner := bufio.NewScanner(bytes.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		// We need at least 5 fields (Name, Version, Rev, Tracking, Publisher/Notes)
		// and Notes should include "disabled".
		if len(fields) >= 5 {
			name := fields[0]
			rev := fields[2]
			notes := strings.Join(fields[4:], " ")
			if !strings.Contains(notes, "disabled") || !isValidSnapID(name, rev) {
				continue
			}

			// Perform smart file lookup
			realPath, sizeBytes, found := findSnapFile(name, rev)

			humanSize := "Ghost"
			if found {
				humanSize = formatSize(sizeBytes)
			}

			item := SnapItem{
				Name:      name,
				Version:   fields[1],
				Revision:  rev,
				SizeBytes: sizeBytes,
				SizeHuman: humanSize,
				FileFound: found,
				FilePath:  realPath,
			}
			candidates = append(candidates, item)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse snap list output: %v", err)
	}

	return candidates, nil
}

// RemoveSnap executes the snap remove command for a specific revision
func RemoveSnap(item SnapItem) error {
	if !isValidSnapID(item.Name, item.Revision) {
		return fmt.Errorf("invalid snap name or revision")
	}

	snapPath, err := getSnapBinaryPath()
	if err != nil {
		return err
	}

	// snap remove <name> --revision=<id>
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, snapPath, "remove", item.Name, fmt.Sprintf("--revision=%s", item.Revision))
	return cmd.Run()
}
