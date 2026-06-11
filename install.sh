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

echo -e "${BLUE}[*] Fetching available releases from GitHub...${NC}"
RELEASES_JSON=$(curl -s "https://api.github.com/repos/$REPO/releases")
TAGS=($(echo "$RELEASES_JSON" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'))

if [ ${#TAGS[@]} -eq 0 ]; then
    echo -e "${RED}[!] Error: Could not retrieve release tags.${NC}"
    exit 1
fi

echo -e "\nAvailable versions:"
MAX_RELEASES=5
if [ ${#TAGS[@]} -lt $MAX_RELEASES ]; then
    MAX_RELEASES=${#TAGS[@]}
fi

for i in $(seq 0 $((MAX_RELEASES-1))); do
    echo -e "  $((i+1)). ${TAGS[$i]}"
done

echo ""
read -p "Select version number [1-$MAX_RELEASES]: " VER_SELECTION </dev/tty

if ! [[ "$VER_SELECTION" =~ ^[0-9]+$ ]] || [ "$VER_SELECTION" -lt 1 ] || [ "$VER_SELECTION" -gt "$MAX_RELEASES" ]; then
    echo -e "${RED}[!] Invalid version selection.${NC}"; exit 1
fi

INDEX=$((VER_SELECTION-1))
SELECTED_TAG="${TAGS[$INDEX]}"

LATEST_JSON=$(curl -s "https://api.github.com/repos/$REPO/releases/tags/$SELECTED_TAG")
ASSETS=($(echo "$LATEST_JSON" | grep '"name":' | grep 'gits-' | sed -E 's/.*"([^"]+)".*/\1/'))

if [ ${#ASSETS[@]} -eq 0 ]; then
    echo -e "${RED}[!] Error: No binaries found for version $SELECTED_TAG.${NC}"
    exit 1
fi

OS_TYPE="linux"
if [[ "$OSTYPE" == "darwin"* ]]; then OS_TYPE="darwin"; fi
if [[ "$OSTYPE" == "msys"* || "$OSTYPE" == "cygwin"* || "$OSTYPE" == "win32" ]]; then OS_TYPE="windows"; fi
if [ -d "/data/data/com.termux" ]; then OS_TYPE="android"; fi

ARCH_TYPE="amd64"
if [[ "$(uname -m)" == "arm64" || "$(uname -m)" == "aarch64" ]]; then ARCH_TYPE="arm64"; fi

EXT=""
if [ "$OS_TYPE" == "windows" ]; then EXT=".exe"; fi

DETECTED_ASSET="gits-${OS_TYPE}-${ARCH_TYPE}${EXT}"

ASSET_FOUND=false
for asset in "${ASSETS[@]}"; do
    if [ "$asset" == "$DETECTED_ASSET" ]; then
        ASSET_FOUND=true
        break
    fi
done

USE_AUTO="n"
if [ "$ASSET_FOUND" = true ]; then
    read -p "Use auto-detected binary '$DETECTED_ASSET'? [Y/n]: " AUTO_CHOICE </dev/tty
    case "$AUTO_CHOICE" in
        y|Y|"" ) USE_AUTO="y" ;;
        * ) USE_AUTO="n" ;;
    esac
fi

if [ "$USE_AUTO" == "y" ]; then
    SELECTED_ASSET="$DETECTED_ASSET"
else
    echo -e "\nAvailable binaries for $SELECTED_TAG:"
    for i in "${!ASSETS[@]}"; do
        echo -e "  $((i+1)). ${ASSETS[$i]}"
    done

    echo ""
    read -p "Select binary number [1-${#ASSETS[@]}]: " SELECTION </dev/tty

    if ! [[ "$SELECTION" =~ ^[0-9]+$ ]] || [ "$SELECTION" -lt 1 ] || [ "$SELECTION" -gt "${#ASSETS[@]}" ]; then
        echo -e "${RED}[!] Invalid binary selection.${NC}"; exit 1
    fi

    INDEX=$((SELECTION-1))
    SELECTED_ASSET="${ASSETS[$INDEX]}"
fi

URL="https://github.com/$REPO/releases/download/$SELECTED_TAG/$SELECTED_ASSET"

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

echo -e "\n${BLUE}[*] Downloading $SELECTED_ASSET...${NC}"
TMP_FILE="./$SELECTED_ASSET"

curl -L -q -# "$URL" -o "$TMP_FILE"

if [ ! -f "$TMP_FILE" ] || [ ! -s "$TMP_FILE" ]; then
    echo -e "${RED}[!] Download failed.${NC}"; rm -f "$TMP_FILE"; exit 1
fi

echo -e "${BLUE}[*] Installing to $INSTALL_PATH/$BINARY_NAME...${NC}"

if [ -n "$SUDO" ] && command -v sudo &> /dev/null; then
    $SUDO mv -f "$TMP_FILE" "$INSTALL_PATH/$BINARY_NAME"
else
    mv -f "$TMP_FILE" "$INSTALL_PATH/$BINARY_NAME"
fi

echo -e "${GREEN}[+] GitS $SELECTED_TAG installed successfully.${NC}"