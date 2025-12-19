<div align="center">

# 📦 SnapMan
### The Missing Garbage Collector for Ubuntu Snaps

[![Go Version](https://img.shields.io/github/go-mod/go-version/ersinakyuz/snapman?style=flat-square&color=00ADD8)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=flat-square)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-linux-lightgrey?style=flat-square&logo=linux)](https://ubuntu.com/)
[![Status](https://img.shields.io/badge/status-active-success?style=flat-square)]()

<p align="center">
  <b>SnapMan</b> is a lightweight CLI tool written in Go to safely clean up old, disabled Snap revisions <br> 
  and reclaim disk space on your Ubuntu system.
</p>

[Installation](#-installation) • [Usage](#-usage) • [Roadmap](#-roadmap)

</div>

---

## 🚀 Why SnapMan?

Ubuntu keeps older versions of snaps (by default 3 revisions) for safety during updates. However, over time, these disabled revisions can accumulate and occupy gigabytes of disk space.

SnapMan helps you by:
* 🧹 Parsing `snap list --all` output efficiently.
* 🔍 Identifying only **disabled** revisions.
* 🛡️ **Safe Execution:** Targets specific revisions without affecting active packages.
* 🌍 **Locale Safe:** Works correctly even if your system language is set to German, Turkish, etc.

## 🛠️ Installation

You can build SnapMan from source. Ensure you have **Go 1.20+** installed.

```bash
# 1. Clone the repository
git clone [https://github.com/ersinakyuz/snapman.git](https://github.com/ersinakyuz/snapman.git)
cd snapman

# 2. Tidy dependencies
go mod tidy

# 3. Build the binary
go build -o snapman ./cmd/snapman

# 4. (Optional) Install system-wide
sudo mv snapman /usr/local/bin/
```

## 💻 Usage

Managing snaps requires root privileges to execute removal commands.
1. Scan and Clean

Run the tool with sudo:
Bash

`sudo ./snapman`

(If you moved the binary to /usr/local/bin, simply run sudo snapman)
Example Output
Plaintext
```bash
--- SnapMan v1.0: Cleanup Tool ---
Scanning system...

3 disabled revisions found:

PACKAGE  VERSION   REVISION
-------  -------   --------
spotify  1.1.84    x1
vlc      3.0.18    22
code     1.74.0    old_ver

Starting cleanup...
Removing: spotify (Rev: x1)... ✓ OK
Removing: vlc (Rev: 22)... ✓ OK
Removing: code (Rev: old_ver)... ✓ OK

Operation Complete.
```

## 📂 Project Structure

This project follows the Standard Go Project Layout to ensure maintainability.

```text
snapman/
├── cmd/
│   └── snapman/     # Main application entry point
├── internal/
│   └── snapsys/     # Business logic (System calls, parsing)
├── go.mod           # Module definition
└── README.md        # Documentation
```

## 🗺️ Roadmap

    [x] Basic CLI (List and Remove)

    [x] Locale-independent parsing (LC_ALL=C logic)

    [ ] Dry-Run mode (--dry-run flag for simulation)

    [ ] Interactive Confirmation (Yes/No prompt)

    [ ] TUI (Terminal User Interface) - Planned with Bubbletea

    [ ] GUI (Graphical User Interface) - Planned with Fyne or Wails

## 🤝 Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

    Fork the Project
    
    Create your Feature Branch (`git checkout -b feat/AmazingFeature`)
    
    Commit your Changes (`git commit -m 'feat: Add some AmazingFeature'`)
    
    Push to the Branch (`git push origin feat/AmazingFeature`)
    
    Open a Pull Request
    

## 📝 License

Distributed under the **MIT** License. See **LICENSE** for more information.
