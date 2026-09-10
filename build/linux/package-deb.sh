#!/usr/bin/env bash
set -e

# ==============================================================================
# IT Toolbox - Debian (.deb) Package Builder
# Menghasilkan installer .deb siap pakai (tinggal double-click untuk install)
# dengan auto-shortcut di Desktop & Start Menu aplikasi.
# ==============================================================================

VERSION="${1:-1.0.0}"
ARCH="amd64"
PKG_NAME="it-toolbox"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"
DIST_DIR="${ROOT_DIR}/dist"
STAGE_DIR=$(mktemp -d /tmp/it-toolbox-deb-build.XXXXXX)

# Trap untuk membersihkan staging directory di /tmp saat selesai/gagal
trap 'rm -rf "${STAGE_DIR}"' EXIT

echo "========================================================"
echo " Building Debian Package: ${PKG_NAME} v${VERSION} (${ARCH})"
echo "========================================================"

# Bersihkan staging lama
rm -rf "${STAGE_DIR}"
mkdir -p "${STAGE_DIR}/DEBIAN"
mkdir -p "${STAGE_DIR}/usr/bin"
mkdir -p "${STAGE_DIR}/usr/share/applications"
mkdir -p "${STAGE_DIR}/usr/share/icons/hicolor/512x512/apps"
mkdir -p "${STAGE_DIR}/usr/share/icons/hicolor/256x256/apps"
mkdir -p "${STAGE_DIR}/usr/share/icons/hicolor/128x128/apps"
mkdir -p "${STAGE_DIR}/usr/share/pixmaps"

# 1. Compile Go binary dengan flag optimal
echo "[1/5] Mengompilasi binary IT Toolbox (Linux x86_64)..."
cd "${ROOT_DIR}"
go build -tags x11 -ldflags "-s -w" -o "${STAGE_DIR}/usr/bin/${PKG_NAME}" .
chmod 755 "${STAGE_DIR}/usr/bin/${PKG_NAME}"

# 2. Salin Desktop Entry & Icon
echo "[2/5] Menyalin icon dan desktop shortcut template..."
cp "${ROOT_DIR}/assets/it-toolbox.desktop" "${STAGE_DIR}/usr/share/applications/${PKG_NAME}.desktop"
chmod 644 "${STAGE_DIR}/usr/share/applications/${PKG_NAME}.desktop"

cp "${ROOT_DIR}/assets/icon.png" "${STAGE_DIR}/usr/share/icons/hicolor/512x512/apps/${PKG_NAME}.png"
chmod 644 "${STAGE_DIR}/usr/share/icons/hicolor/512x512/apps/${PKG_NAME}.png"
cp "${ROOT_DIR}/assets/icon.png" "${STAGE_DIR}/usr/share/icons/hicolor/256x256/apps/${PKG_NAME}.png"
chmod 644 "${STAGE_DIR}/usr/share/icons/hicolor/256x256/apps/${PKG_NAME}.png"
cp "${ROOT_DIR}/assets/icon.png" "${STAGE_DIR}/usr/share/icons/hicolor/128x128/apps/${PKG_NAME}.png"
chmod 644 "${STAGE_DIR}/usr/share/icons/hicolor/128x128/apps/${PKG_NAME}.png"

# Link ke pixmaps agar kompatibel dengan seluruh desktop manager (XFCE, LXDE, Cinnamon, MATE)
cp "${ROOT_DIR}/assets/icon.png" "${STAGE_DIR}/usr/share/pixmaps/${PKG_NAME}.png"
chmod 644 "${STAGE_DIR}/usr/share/pixmaps/${PKG_NAME}.png"

# 3. Buat file DEBIAN/control
echo "[3/5] Membuat DEBIAN/control metadata..."
cat <<EOF > "${STAGE_DIR}/DEBIAN/control"
Package: ${PKG_NAME}
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: ${ARCH}
Maintainer: Yudz <yudz@example.com>
Installed-Size: $(du -ks "${STAGE_DIR}/usr" | cut -f1)
Depends: libgl1, libx11-6, libxcursor1, libxrandr2, libxinerama1, libxi6, libxxf86vm1
Description: IT Toolbox - All-in-one Developer and Network Utility Suite
 IT Toolbox provides networking diagnostics, data conversion, encoding/decoding,
 JWT inspection, hashing tools, and system utilities in a modern GUI.
EOF
chmod 644 "${STAGE_DIR}/DEBIAN/control"

# 4. Buat script DEBIAN/postinst (Auto-create Desktop Shortcut langsung muncul)
echo "[4/5] Mengonfigurasi post-installation auto-shortcut trigger..."
cat <<'EOF' > "${STAGE_DIR}/DEBIAN/postinst"
#!/bin/sh
set -e

# Update cache aplikasi & icon sistem
if command -v update-desktop-database >/dev/null 2>&1; then
    update-desktop-database -q /usr/share/applications || true
fi
if command -v gtk-update-icon-cache >/dev/null 2>&1; then
    gtk-update-icon-cache -q -t -f /usr/share/icons/hicolor || true
fi

# Buat shortcut langsung di Desktop untuk setiap user aktif
for user_dir in /home/*; do
    if [ -d "$user_dir" ]; then
        user_name=$(basename "$user_dir")
        
        # Cek folder Desktop atau desktop
        for desk in "$user_dir/Desktop" "$user_dir/desktop"; do
            if [ -d "$desk" ]; then
                target_file="$desk/it-toolbox.desktop"
                cp -f /usr/share/applications/it-toolbox.desktop "$target_file"
                chown "$user_name:$user_name" "$target_file" || true
                chmod +x "$target_file" || true
                
                # Trust desktop file di GNOME / Ubuntu agar bisa langsung di-double click
                if command -v gio >/dev/null 2>&1; then
                    su - "$user_name" -c "gio set '$target_file' metadata::trusted true" 2>/dev/null || true
                fi
            fi
        done
    fi
done

exit 0
EOF
chmod 755 "${STAGE_DIR}/DEBIAN/postinst"

# 5. Buat script DEBIAN/postrm (Clean up shortcut saat di-uninstall)
cat <<'EOF' > "${STAGE_DIR}/DEBIAN/postrm"
#!/bin/sh
set -e

if [ "$1" = "remove" ] || [ "$1" = "purge" ]; then
    # Hapus shortcut di desktop user
    for user_dir in /home/*; do
        for desk in "$user_dir/Desktop" "$user_dir/desktop"; do
            if [ -f "$desk/it-toolbox.desktop" ]; then
                rm -f "$desk/it-toolbox.desktop" || true
            fi
        done
    done
    
    if command -v update-desktop-database >/dev/null 2>&1; then
        update-desktop-database -q /usr/share/applications || true
    fi
    if command -v gtk-update-icon-cache >/dev/null 2>&1; then
        gtk-update-icon-cache -q -t -f /usr/share/icons/hicolor || true
    fi
fi

exit 0
EOF
chmod 755 "${STAGE_DIR}/DEBIAN/postrm"

# 6. Build .deb package
OUTPUT_DEB="${DIST_DIR}/${PKG_NAME}_${VERSION}_${ARCH}.deb"
mkdir -p "${DIST_DIR}"
echo "[5/5] Mengemas package .deb menggunakan dpkg-deb..."
dpkg-deb --build --root-owner-group "${STAGE_DIR}" "${OUTPUT_DEB}"

# Bersihkan staging
rm -rf "${DIST_DIR}/deb-staging"

echo ""
echo "========================================================"
echo " SUKSES! Installer Linux telah dibuat:"
echo " File: ${OUTPUT_DEB}"
echo " Ukuran: $(du -h "${OUTPUT_DEB}" | cut -f1)"
echo "========================================================"
echo "Cara Penggunaan oleh User (Klik-klik aja):"
echo "1. Double-click file '${PKG_NAME}_${VERSION}_${ARCH}.deb'"
echo "2. Klik tombol 'Install' pada Software Center / GDebi"
echo "3. Shortcut 'IT Toolbox' akan LANGSUNG muncul di Desktop"
echo "   dan di Application Menu / Launcher!"
echo "========================================================"
