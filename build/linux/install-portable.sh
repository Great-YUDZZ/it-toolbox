#!/usr/bin/env bash
# ==============================================================================
# IT Toolbox - 1-Click User Installer for Linux (No sudo required)
# ==============================================================================
# Skrip ini menginstal IT Toolbox ke direktori pengguna (~/.local/bin)
# dan langsung memunculkan shortcut di Desktop dan Application Menu.
# Bisa dijalankan dengan double-click atau via terminal.
# ==============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"
APP_NAME="IT Toolbox"
BIN_NAME="it-toolbox"
INSTALL_BIN_DIR="${HOME}/.local/bin"
INSTALL_APPS_DIR="${HOME}/.local/share/applications"
INSTALL_ICON_DIR="${HOME}/.local/share/icons/hicolor/128x128/apps"

# Deteksi GUI dialog (zenity atau kdialog)
HAS_ZENITY=false
if command -v zenity >/dev/null 2>&1; then
    HAS_ZENITY=true
fi

show_info() {
    local title="$1"
    local message="$2"
    if [ "$HAS_ZENITY" = true ]; then
        zenity --info --title="$title" --text="$message" --width=350 2>/dev/null || true
    else
        echo -e "\n[$title] $message"
    fi
}

show_question() {
    local title="$1"
    local message="$2"
    if [ "$HAS_ZENITY" = true ]; then
        zenity --question --title="$title" --text="$message" --width=350 2>/dev/null
        return $?
    else
        echo -e "\n[$title] $message (y/n)"
        read -r ans
        case "$ans" in
            [Yy]* ) return 0 ;;
            * ) return 1 ;;
        esac
    fi
}

# Dialog konfirmasi instalasi (Tinggal klik-klik)
if ! show_question "Instalasi IT Toolbox" "Apakah Anda ingin menginstal ${APP_NAME}?\n\nShortcut akan otomatis dibuat di Desktop dan Menu Aplikasi."; then
    echo "Instalasi dibatalkan oleh pengguna."
    exit 0
fi

# Pastikan binary sudah ter-compile
SOURCE_BIN="${ROOT_DIR}/${BIN_NAME}"
if [ ! -f "${SOURCE_BIN}" ]; then
    echo "Mengompilasi binary IT Toolbox..."
    (cd "${ROOT_DIR}" && go build -tags x11 -ldflags "-s -w" -o "${SOURCE_BIN}" .)
fi

# 1. Buat direktori tujuan
mkdir -p "${INSTALL_BIN_DIR}"
mkdir -p "${INSTALL_APPS_DIR}"
mkdir -p "${INSTALL_ICON_DIR}"

# 2. Salin binary & icon
cp -f "${SOURCE_BIN}" "${INSTALL_BIN_DIR}/${BIN_NAME}"
chmod 755 "${INSTALL_BIN_DIR}/${BIN_NAME}"

if [ -f "${ROOT_DIR}/assets/icon.png" ]; then
    cp -f "${ROOT_DIR}/assets/icon.png" "${INSTALL_ICON_DIR}/${BIN_NAME}.png"
fi

# 3. Buat Desktop Entry
DESKTOP_FILE="${INSTALL_APPS_DIR}/${BIN_NAME}.desktop"
cat <<EOF > "${DESKTOP_FILE}"
[Desktop Entry]
Version=1.0
Type=Application
Name=${APP_NAME}
GenericName=Developer & Network Utilities
Comment=All-in-one IT Utility suite for networking, security, data conversion, and developer tools
Exec=${INSTALL_BIN_DIR}/${BIN_NAME}
Icon=${INSTALL_ICON_DIR}/${BIN_NAME}.png
Terminal=false
Categories=Utility;Development;Network;
Keywords=it;toolbox;converter;network;hash;jwt;ping;dns;
StartupNotify=true
StartupWMClass=com.yudz.it-toolbox
EOF
chmod 755 "${DESKTOP_FILE}"

# 4. Salin langsung ke Desktop pengguna
DESKTOP_DIR="${HOME}/Desktop"
if [ ! -d "${DESKTOP_DIR}" ] && [ -d "${HOME}/desktop" ]; then
    DESKTOP_DIR="${HOME}/desktop"
fi

if [ -d "${DESKTOP_DIR}" ]; then
    TARGET_DESKTOP="${DESKTOP_DIR}/${BIN_NAME}.desktop"
    cp -f "${DESKTOP_FILE}" "${TARGET_DESKTOP}"
    chmod +x "${TARGET_DESKTOP}"
    
    # Set metadata trusted untuk GNOME/Ubuntu desktop
    if command -v gio >/dev/null 2>&1; then
        gio set "${TARGET_DESKTOP}" metadata::trusted true 2>/dev/null || true
    fi
fi

# 5. Update cache desktop
if command -v update-desktop-database >/dev/null 2>&1; then
    update-desktop-database -q "${INSTALL_APPS_DIR}" || true
fi
if command -v gtk-update-icon-cache >/dev/null 2>&1; then
    gtk-update-icon-cache -q -t -f "${HOME}/.local/share/icons/hicolor" || true
fi

# 6. Tampilkan pesan sukses dan opsi jalankan sekarang
if show_question "Instalasi Selesai!" "Instalasi ${APP_NAME} berhasil!\n\nShortcut telah muncul di Desktop dan Menu Aplikasi Anda.\n\nApakah Anda ingin menjalankan ${APP_NAME} sekarang?"; then
    "${INSTALL_BIN_DIR}/${BIN_NAME}" &
fi
