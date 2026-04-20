#!/bin/bash

# ==========================================
#             GitS Installer
# ==========================================
# Author: ramadanny
# Repository: github.com/ramadanny/gits
# ==========================================

REPO="ramadanny/gits"
BINARY_NAME="gits"

echo -e "\033[1;36m==========================================\033[0m"
echo -e "\033[1;36m             GitS Installer               \033[0m"
echo -e "\033[1;36m==========================================\033[0m"

# 1. Dependency Check
check_dep() {
    if ! command -v $1 &> /dev/null; then
        echo -e "\x1b[33m[!] Dependency '$1' is missing.\x1b[0m"
        read -p "Install it now? [Y/n]: " choice </dev/tty
        case "$choice" in 
            y|Y|"" ) 
                echo -e "\x1b[36m[*] Installing $1...\x1b[0m"
                if [ -d "/data/data/com.termux" ]; then pkg install $1 -y
                elif [[ "$OSTYPE" == "linux-gnu"* ]]; then sudo apt update && sudo apt install $1 -y
                elif [[ "$OSTYPE" == "darwin"* ]]; then brew install $1
                fi ;;
            * ) echo -e "\x1b[31m[!] Aborted. '$1' is required.\x1b[0m"; exit 1 ;;
        esac
    fi
}

check_dep "curl"
check_dep "git"

# 2. Fetch Latest Release Data
echo -e "\x1b[36m[*] Fetching latest assets from GitHub...\x1b[0m"
LATEST_JSON=$(curl -s "https://api.github.com/repos/$REPO/releases/latest")
LATEST_TAG=$(echo "$LATEST_JSON" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo -e "\x1b[31m[!] Failed to fetch release data.\x1b[0m"
    exit 1
fi

# 3. Parse Assets
ASSETS=($(echo "$LATEST_JSON" | grep '"name":' | grep 'gits-' | sed -E 's/.*"([^"]+)".*/\1/'))

if [ ${#ASSETS[@]} -eq 0 ]; then
    echo -e "\x1b[31m[!] No binary assets found for $LATEST_TAG.\x1b[0m"
    exit 1
fi

# 4. Interactive Selection Menu
echo -e "\n\x1b[32mAvailable Binaries (Version $LATEST_TAG):\x1b[0m"
for i in "${!ASSETS[@]}"; do
    echo -e "  \x1b[33m$((i+1)).\x1b[0m ${ASSETS[$i]}"
done

echo ""
read -p "Select the number corresponding to your OS/Arch [1-${#ASSETS[@]}]: " SELECTION </dev/tty

if ! [[ "$SELECTION" =~ ^[0-9]+$ ]] || [ "$SELECTION" -lt 1 ] || [ "$SELECTION" -gt "${#ASSETS[@]}" ]; then
    echo -e "\x1b[31m[!] Invalid selection. Aborting.\x1b[0m"
    exit 1
fi

INDEX=$((SELECTION-1))
SELECTED_ASSET="${ASSETS[$INDEX]}"
URL="https://github.com/$REPO/releases/download/$LATEST_TAG/$SELECTED_ASSET"

# 5. Determine Installation Path
if [[ "$SELECTED_ASSET" == *"android"* ]]; then
    INSTALL_PATH="${PREFIX:-/data/data/com.termux/files/usr}/bin"
    SUDO=""
elif [[ "$SELECTED_ASSET" == *"windows"* ]]; then
    INSTALL_PATH="$PWD"
    SUDO=""
    BINARY_NAME="gits.exe"
    echo -e "\x1b[33m[*] Windows Mode: Downloading to current directory.\x1b[0m"
else
    INSTALL_PATH="/usr/local/bin"
    SUDO="sudo"
fi

# 6. Execution
echo -e "\x1b[36m[*] Downloading $SELECTED_ASSET...\x1b[0m"
TMP_FILE="/tmp/$SELECTED_ASSET"
curl -L -q --show-progress "$URL" -o "$TMP_FILE"

if [ ! -s "$TMP_FILE" ] || grep -q "Not Found" "$TMP_FILE"; then
    echo -e "\x1b[31m[!] Download failed.\x1b[0m"
    rm -f "$TMP_FILE"
    exit 1
fi

chmod +x "$TMP_FILE"
echo -e "\x1b[36m[*] Moving binary to $INSTALL_PATH/$BINARY_NAME...\x1b[0m"

if [ -n "$SUDO" ] && command -v sudo &> /dev/null; then
    $SUDO mv -f "$TMP_FILE" "$INSTALL_PATH/$BINARY_NAME"
else
    mv -f "$TMP_FILE" "$INSTALL_PATH/$BINARY_NAME"
fi

echo -e "\033[1;32m[+] GitS $LATEST_TAG installed successfully!\033[0m"
echo -e "Run '\x1b[33m$BINARY_NAME --help\x1b[0m' to begin."