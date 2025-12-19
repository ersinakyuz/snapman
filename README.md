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
git clone https://github.com/ersinakyuz/snapman.git
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
SnapMan starting...

System scanning...
Found 15 disabled revisions. Potential gain: 3.57 GB

PACKAGE                 VERSION                     REV    SIZE        STATUS
-------                 -------                     ---    ----        ------
brave                   1.85.117                    578    186.37 MB   Ready
chromium-ffmpeg         120726-120170-119605-1...   88     7.27 MB     Ready
core18                  20251001                    2959   55.49 MB    Ready
core20                  20250822                    2682   63.77 MB    Ready
core22                  20250923                    2139   73.91 MB    Ready
core24                  20251001                    1225   66.84 MB    Ready
firefox                 145.0.2-1                   7423   250.58 MB   Ready
gnome-46-2404           0+git.4ca00c0-sdk0+git...   125    618.26 MB   Ready
libreoffice             25.8.3.2                    362    1.17 GB     Ready
lxqt-support            2025-10                     7      4.21 MB     Ready
mesa-2404               24.2.8-snap185              912    290.77 MB   Ready
opera                   125.0.5729.21               416    179.11 MB   Ready
telegram-desktop        6.3.8                       6880   82.02 MB    Ready
thunderbird             140.5.0esr-2                915    226.32 MB   Ready
wine-platform-runtime   v1.0                        399    346.89 MB   Ready


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
