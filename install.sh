#!/bin/bash

# ==========================================
#             GitS Installer
#        Aesthetic & Pro Edition
# ==========================================

REPO="ramadanny/gits"
BINARY_NAME="gits"

# Colors
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
RED='\033[0;31m'
BOLD='\033[1m'
NC='\033[0m' # No Color

clear

# Aesthetic ASCII Art Header
echo -e "${MAGENTA}${BOLD}"
echo "    ________  _________  ____"
echo "   / ____/ / / / ___/ / / / /"
echo "  / / __/ / / /\__ \ / / / / "
echo " / /_/ / /_/ /___/ // /_/ /  "
echo " \____/\____//____/ \____/   "
echo -e "      ${CYAN}GitS Toolchain Installer${NC}"
echo -e "${BLUE}  ----------------------------------${NC}"

# 1. Dependency Check
check_dep() {
    if ! command -v $1 &> /dev/null; then
        echo -e "${YELLOW}[!]${NC} Component ${BOLD}'$1'${NC} is missing."
        echo -ne "${CYAN}[?]${NC} Install it now? (y/n): "
        read -r choice </dev/tty
        case "$choice" in 
            y|Y|"" ) 
                echo -e "${BLUE}[*]${NC} Deploying $1..."
                if [ -d "/data/data/com.termux" ]; then pkg install $1 -y
                elif [[ "$OSTYPE" == "linux-gnu"* ]]; then sudo apt update && sudo apt install $1 -y
                elif [[ "$OSTYPE" == "darwin"* ]]; then brew install $1
                fi ;;
            * ) echo -e "${RED}[!]${NC} Aborted. '$1' is required."; exit 1 ;;
        esac
    fi
}

check_dep "curl"
check_dep "git"

# 2. Fetch Data
echo -e "${BLUE}[*]${NC} Communicating with GitHub API..."
LATEST_JSON=$(curl -s "https://api.github.com/repos/$REPO/releases/latest")
LATEST_TAG=$(echo "$LATEST_JSON" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo -e "${RED}[!] Gagal sinkronisasi data rilis.${NC}"
    exit 1
fi

# 3. Assets Menu
ASSETS=($(echo "$LATEST_JSON" | grep '"name":' | grep 'gits-' | sed -E 's/.*"([^"]+)".*/\1/'))

echo -e "\n${MAGENTA}┌─${NC} ${BOLD}Select Environment${NC} (${LATEST_TAG})"
for i in "${!ASSETS[@]}"; do
    echo -e "${MAGENTA}├─ [${NC}${CYAN}$((i+1))${NC}${MAGENTA}]${NC} ${ASSETS[$i]}"
done
echo -e "${MAGENTA}└───────────────────────────────${NC}"

echo -ne "${CYAN}>>${NC} Choose Index: "
read -r SELECTION </dev/tty

if ! [[ "$SELECTION" =~ ^[0-9]+$ ]] || [ "$SELECTION" -lt 1 ] || [ "$SELECTION" -gt "${#ASSETS[@]}" ]; then
    echo -e "${RED}[!] Invalid choice.${NC}"; exit 1
fi

INDEX=$((SELECTION-1))
SELECTED_ASSET="${ASSETS[$INDEX]}"
URL="https://github.com/$REPO/releases/download/$LATEST_TAG/$SELECTED_ASSET"

# 4. Path Configuration
if [[ "$SELECTED_ASSET" == *"android"* ]]; then
    INSTALL_PATH="${PREFIX:-/data/data/com.termux/files/usr}/bin"
    SUDO=""
elif [[ "$SELECTED_ASSET" == *"windows"* ]]; then
    INSTALL_PATH="$PWD"
    SUDO=""
    BINARY_NAME="gits.exe"
else
    INSTALL_PATH="/usr/local/bin"
    SUDO="sudo"
fi

# 5. Execution
echo -e "\n${BLUE}[*]${NC} Downloading ${BOLD}$SELECTED_ASSET${NC}..."
TMP_FILE="/tmp/$SELECTED_ASSET"

# Progress bar modern
curl -L -q -# "$URL" -o "$TMP_FILE"

if [ ! -s "$TMP_FILE" ] || grep -q "Not Found" "$TMP_FILE"; then
    echo -e "${RED}[!] Transmission error.${NC}"; rm -f "$TMP_FILE"; exit 1
fi

chmod +x "$TMP_FILE"
echo -e "${BLUE}[*]${NC} Installing to ${BOLD}$INSTALL_PATH${NC}..."

if [ -n "$SUDO" ] && command -v sudo &> /dev/null; then
    $SUDO mv -f "$TMP_FILE" "$INSTALL_PATH/$BINARY_NAME"
else
    mv -f "$TMP_FILE" "$INSTALL_PATH/$BINARY_NAME"
fi

# 6. Success Message
echo -e "\n${GREEN}${BOLD}SUCCESS!${NC} GitS ${LATEST_TAG} is now active."
echo -e "${MAGENTA}----------------------------------${NC}"
echo -e "Type '${CYAN}${BINARY_NAME} --help${NC}' to explore commands."