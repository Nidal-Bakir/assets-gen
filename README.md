# assets-gen

[![assets-gen](https://github.com/Nidal-Bakir/assets-gen/actions/workflows/go.yml/badge.svg)](https://github.com/Nidal-Bakir/assets-gen/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/Nidal-Bakir/assets-gen.svg)](https://pkg.go.dev/github.com/Nidal-Bakir/assets-gen)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A flexible **Go** CLI for generating Android & iOS app icons, Android notification icons, Google Play store logos, and Android image assets at all DPIs.

---

## 🚀 Features

- **Android App Icon**
  🎨 Rounded corners, gradient or solid-color backgrounds, whitespace trimming & padding, DPI‐specific mipmaps.
- **Android Notification Icon**
  🔔 Generate all notification‐icon mipmaps (ic_stat\_), with alpha‐threshold and whitespace trim.
- **Android Asset Generator**
  📐 Produce drawable assets across all DPIs from a single image.
- **Android Google Play Logo**
  🛒 Create 512×512 Play Store logos with backgrounds, padding, and trim options.
- **iOS App Icon**
  📱 Export all AppIcon sizes into your `AppIcon.appiconset`, with padding, and BG options.
- **Generate All (`all`)**
  🎉 Run all of the above tasks in parallel using a single source image.

---

## 📦 Installation

### via go (requires Go 1.18+)

```bash
# via go (requires Go 1.18+)
go install github.com/Nidal-Bakir/assets-gen/cmd/assetsgen@latest
```

### released binary

You can download the released binary for your OS from the [releases](https://github.com/Nidal-Bakir/assets-gen/releases/latest) page.

### from source

```bash
# or clone & build using Makefile
git clone https://github.com/Nidal-Bakir/assets-gen.git
cd assets-gen
make        # builds executable as `assetsgen`
cd build
./assetsgen --help
```

Ensure your PATH includes the build output (`./build/assetsgen`) or rename the binary accordingly.

---

## 🛠 Usage

Run `assetsgen --help` (or any subcommand `--help`) for a full list of flags and examples. Below are the main commands:

---

### 1. Android App Icon (`aai`)

Generate launcher icons in all mipmap-density folders.

```bash
# help:
assetsgen android-app-icon --help

# minimal: just supply a source PNG
assetsgen android-app-icon ./ic_launcher.png
# alias
assetsgen aai ./ic_launcher.png

# solid-color BG:
assetsgen aai --color "#3498db" ./ic_launcher.png

# linear-gradient BG with stops & degree:
assetsgen aai \
  --bg linear-gradient \
  --colors "#FF0000,#00FF00,#0000FF" \
  --stops "0.0,0.5,1.0" \
  --degree 90 \
  ./ic_launcher.png

# image as a background:
assetsgen aai \
  --bg image \
  --bg-path ./bg_image.png \
  ./ic_launcher.png

# trim whitespace, add padding, custom output & apply directly:
assetsgen aai --trim --padding 0.1 --corner-radius 0.5 -o "app_icon" --apply ./ic_launcher.png
```

---

### 2. Android Notification Icon (`ani`)

Produce notification icons (ic_stat\_) across all densities.

```bash
# help:
assetsgen android-notification-icon --help

# basic:
assetsgen android-notification-icon ./notif.png
# alias
assetsgen ani ./notif.png

# with trim, custom name, and apply:
assetsgen ani --trim -o "ic_stat_notification" --apply ./notif.png
```

---

### 3. Android Asset Generator (`aag`)

Generate drawable image assets (all DPIs) from a single source.

```bash
# help:
assetsgen android-asset-gen --help

# generate into default `drawable-<dpi>` folders:
assetsgen android-asset-gen ./image.png
# alias
assetsgen aag ./image.png

# trim whitespace, and apply:
assetsgen aag --trim --apply ./image.png
```

---

### 4. Android Google Play Logo (`agpl`)

Create a 512×512 Play Store logo with optional BG styling.

```bash
# help:
assetsgen android-google-play-logo --help

# basic:
assetsgen android-google-play-logo ./playlogo.png
# alias
assetsgen agpl ./playlogo.png

# gradient background + trim + padding:
assetsgen agpl \
  --bg linear-gradient \
  --colors "#ffa500,#ff4500" \
  --stops "0.0,1.0" \
  --degree 45 \
  --trim \
  --padding 0.05 \
  --apply \
  -o "play_store_logo" \
  ./playlogo.png
```

---

### 5. iOS App Icon (`iai`)

Export all required iOS app-icon sizes into your Xcode `AppIcon.appiconset`.

```bash
# help:
assetsgen ios-app-icon --help

# basic:
assetsgen ios-app-icon ./appicon.png
# alias
assetsgen iai ./appicon.png

# solid color bg + padding + trim:
assetsgen iai --color "#8e44ad" --padding 0.1 --trim --apply ./appicon.png

# image as a background:
assetsgen iai \
  --bg image \
  --bg-path ./bg_image.png \
  ./appicon.png

# radial-gradient BG with stops & degree:
assetsgen iai \
  --bg radial-gradient \
  --colors "#FF0000,#00FF00,#0000FF" \
  --stops "0.0,0.5,1.0" \
  ./appicon.png

```

---

### 6. Generate All (`all`)

Run **all** generation tasks in parallel using a single image. Useful in CI or build scripts.

```bash
# help:
assetsgen all --help

# basic:
assetsgen all ./master_image.png

# customize options (same flags as individual commands):
assetsgen all \
  --bg linear-gradient \
  --colors "#123456,#abcdef" \
  --stops "0.0,1.0" \
  --degree 60 \
  --trim \
  --padding 0.1 \
  --corner-radius 0.2 \
  --alpha-threshold 0.8 \
  --apply \
  ./master_image.png
```

_(Use `assetsgen all --help` for full flag list.)_

---

## 💡 Tips & Tricks

- **Aliases**: `aai`, `ani`, `aag`, `agpl`, `iai`, `all` for quick commands.
- **Dry-run**: Omit `--apply` to preview outputs in `assets-gen-out/` without moving into your project.
- **Color Formats**: Hex strings must start with `#`; for gradients provide comma-separated lists.

---

## 🤝 Contributing

1. Fork it
2. Create your feature branch (`git checkout -b feature/x`)
3. Commit your changes (`git commit -m "feat: add foo"`)
4. Push to your branch (`git push origin feature/x`)
5. Open a Pull Request

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

---

## 📄 License

This project is licensed under the **MIT License**. See [LICENSE](LICENSE) for details.

---

> Built with ❤️ by [Nidal Bakir](https://github.com/Nidal-Bakir)
