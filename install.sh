#!/bin/bash

# ==========================================
#             GitS Installer
# ==========================================
# Author: ramadanny
# Repository: github.com/ramadanny/gits
# ==========================================

REPO="ramadanny/gits"
BINARY_NAME="gits"

BLUE='\033[0;34m'
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}[*] Initializing GitS Installer...${NC}"

# 1. Dependency Check
check_dep() {
    if ! command -v $1 &> /dev/null; then
        echo -e "${RED}[!] Dependency '$1' is missing.${NC}"
        read -p "Install $1? [Y/n]: " choice </dev/tty
        case "$choice" in 
            y|Y|"" ) 
                if [ -d "/data/data/com.termux" ]; then pkg install $1 -y
                elif [[ "$OSTYPE" == "linux-gnu"* ]]; then sudo apt update && sudo apt install $1 -y
                elif [[ "$OSTYPE" == "darwin"* ]]; then brew install $1
                fi ;;
            * ) echo -e "${RED}[!] Aborted. '$1' is required.${NC}"; exit 1 ;;
        esac
    fi
}

check_dep "curl"
check_dep "git"

# 2. Fetch Release Data
echo -e "${BLUE}[*] Fetching latest release from GitHub...${NC}"
LATEST_JSON=$(curl -s "https://api.github.com/repos/$REPO/releases/latest")
LATEST_TAG=$(echo "$LATEST_JSON" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo -e "${RED}[!] Error: Could not retrieve release data.${NC}"
    exit 1
fi

# 3. Assets Selection
ASSETS=($(echo "$LATEST_JSON" | grep '"name":' | grep 'gits-' | sed -E 's/.*"([^"]+)".*/\1/'))

if [ ${#ASSETS[@]} -eq 0 ]; then
    echo -e "${RED}[!] Error: No binaries found for version $LATEST_TAG.${NC}"
    exit 1
fi

echo -e "\nAvailable binaries for $LATEST_TAG:"
for i in "${!ASSETS[@]}"; do
    echo -e "  $((i+1)). ${ASSETS[$i]}"
done

echo ""
read -p "Select binary number [1-${#ASSETS[@]}]: " SELECTION </dev/tty

if ! [[ "$SELECTION" =~ ^[0-9]+$ ]] || [ "$SELECTION" -lt 1 ] || [ "$SELECTION" -gt "${#ASSETS[@]}" ]; then
    echo -e "${RED}[!] Invalid selection.${NC}"; exit 1
fi

INDEX=$((SELECTION-1))
SELECTED_ASSET="${ASSETS[$INDEX]}"
URL="https://github.com/$REPO/releases/download/$LATEST_TAG/$SELECTED_ASSET"

# 4. Path Logic
if [[ "$SELECTED_ASSET" == *"android"* ]]; then
    INSTALL_PATH="${PREFIX:-/data/data/com.termux/files/usr}/bin"
    SUDO=""
elif [[ "$SELECTED_ASSET" == *"windows"* ]]; then
    INSTALL_PATH="."
    SUDO=""
    BINARY_NAME="gits.exe"
else
    INSTALL_PATH="/usr/local/bin"
    SUDO="sudo"
fi

# 5. Download and Move
echo -e "\n${BLUE}[*] Downloading $SELECTED_ASSET...${NC}"
# Langsung download ke folder saat ini (.)
TMP_FILE="./$SELECTED_ASSET"

curl -L -q -# "$URL" -o "$TMP_FILE"

if [ ! -f "$TMP_FILE" ] || [ ! -s "$TMP_FILE" ]; then
    echo -e "${RED}[!] Download failed.${NC}"; rm -f "$TMP_FILE"; exit 1
fi

chmod +x "$TMP_FILE"
echo -e "${BLUE}[*] Installing to $INSTALL_PATH/$BINARY_NAME...${NC}"

# Pindahkan dari folder saat ini ke path instalasi
if [ -n "$SUDO" ] && command -v sudo &> /dev/null; then
    $SUDO mv -f "$TMP_FILE" "$INSTALL_PATH/$BINARY_NAME"
else
    mv -f "$TMP_FILE" "$INSTALL_PATH/$BINARY_NAME"
fi

# 6. Final Status
echo -e "${GREEN}[+] GitS $LATEST_TAG installed successfully.${NC}"