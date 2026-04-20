# 📻 RadioGo v2

**RadioGo v2** is a modern, aesthetic, and lightweight terminal-based radio player built with Go and the Charmbracelet Bubble Tea framework. It allows you to stream your favorite online radio stations, manage your station list, and keep track of songs you love—all from the comfort of your terminal.

![Aesthetic TUI](https://img.shields.io/badge/UI-Lipgloss-pink)
![Audio Engine](https://img.shields.io/badge/Audio-FFmpeg-blueviolet)
![Persistence](https://img.shields.io/badge/DB-SQLite-blue)
![License](https://img.shields.io/badge/License-Apache%202.0-purple)

---

## ✨ Features

- 🎹 **Modern TUI**: A beautiful, responsive interface built with `bubbletea`, `lipgloss`, and `huh`.
- 📡 **Station Management**: Easily add, edit, and delete radio stations directly within the app.
- 🎵 **Real-time Metadata**: Automatically parses ICY metadata (Artist - Title) so you always know what's playing.
- ❤️ **"Like" System**: Found a track you love? Save it to your personal "Liked Songs" library with a single keypress.
- 📋 **Clipboard Integration**: Quickly copy liked song information to your clipboard for easy searching later.
- 💾 **Local Persistence**: Uses a CGO-free SQLite database to store your stations and liked songs locally.
- 🎨 **Aesthetic Design**: A carefully crafted purple/pink theme that's easy on the eyes.

---

## 🛠️ Prerequisites

RadioGo v2 relies on **FFmpeg** for audio decoding and playback.

- **Go**: 1.25.3 or higher
- **FFmpeg**: `ffplay` must be available in your system's PATH.

### Installing FFmpeg
- **macOS**: `brew install ffmpeg`
- **Linux**: `sudo apt install ffmpeg` (or your distro's equivalent)
- **Windows**: `choco install ffmpeg` or download from [ffmpeg.org](https://ffmpeg.org/download.html)

---

## 🚀 Installation

You can install RadioGo v2 directly using `go install`:

```bash
go install github.com/JbrownWFU/Radio-Go/cmd/radiogo@latest
```

Alternatively, build it from source:

```bash
git clone https://github.com/JbrownWFU/Radio-Go.git
cd Radio-Go
go build -o radiogo.exe ./cmd/radiogo/main.go
./radiogo.exe
```

---

## 📡 Sample Stations

Here are some stations to get you started:

| Station | Stream URL | Website |
| :--- | :--- | :--- |
| **WFMU** - Freeform Radio | http://stream0.wfmu.org/freeform-high.aac | [wfmu.org](https://wfmu.org/) |
| **Delicious Agony** - Progressive Rock | http://deliciousagony.streamguys1.com/ | [deliciousagony.com](https://www.deliciousagony.com/) |
| **KEXP** - Seattle | https://kexp.streamguys1.com/kexp160.aac | [kexp.org](https://www.kexp.org/) |

---

## ⌨️ Keybindings

| Key | Action |
| :--- | :--- |
| `Tab` | Switch between Home, Likes, and About tabs |
| `Enter` | Play selected station / View song details |
| `a` | Add a new radio station |
| `d` | Delete the selected station or liked song |
| `l` | "Like" the currently playing song |
| `s` | Stop audio playback |
| `c` | Copy liked song info to clipboard |
| `?` | Toggle help menu |
| `q` / `Esc` | Quit the application |

---

## 📂 Project Structure

- `cmd/radiogo`: The main entry point of the application.
- `internal/player`: FFplay-based audio engine and metadata parsing logic.
- `internal/storage`: SQLite database management for stations and likes.
- `internal/ui`: Bubble Tea models and Lipgloss styles for the TUI.

---

## 📄 License

This project is licensed under the **Apache License 2.0**. See the [LICENSE](LICENSE) file for details.

---

*Built with ❤️ using [Bubble Tea](https://github.com/charmbracelet/bubbletea)*

---

> [!NOTE]
> This project and its documentation were co-authored with the assistance of **Gemini CLI**.
