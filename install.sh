#!/bin/sh

# Universal Installer Script for "nextver"
#
# This script is designed to be run directly from the web, for example:
# sh <(curl -sL https://your-domain.com/install.sh)
#
# It automatically detects the OS and architecture, then downloads the
# appropriate binary from a hypothetical GitHub releases page.

# --- Configuration ---
# Replace these with your actual GitHub repository and app name.
GITHUB_REPO="chonla/nextver"
APP_NAME="nextver"
INSTALL_DIR="/usr/local/bin"

# --- Style and Logging Functions ---
# Use color codes if the terminal supports it
if [ -t 1 ]; then
  RED=$(printf '\033[31m')
  GREEN=$(printf '\033[32m')
  YELLOW=$(printf '\033[33m')
  BLUE=$(printf '\033[34m')
  BOLD=$(printf '\033[1m')
  RESET=$(printf '\033[0m')
else
  RED=""
  GREEN=""
  YELLOW=""
  BLUE=""
  BOLD=""
  RESET=""
fi

info() {
  printf "%s${BLUE}%s%s\\n" "${BOLD}" "$1" "${RESET}"
}

success() {
  printf "%s${GREEN}%s%s\\n" "${BOLD}" "$1" "${RESET}"
}

warn() {
  printf "%s${YELLOW}WARN: %s%s\\n" "${BOLD}" "$1" "${RESET}"
}

error() {
  printf "%s${RED}ERROR: %s%s\\n" "${BOLD}" "$1" "${RESET}" >&2
  exit 1
}

# --- Main Installation Logic ---

main() {
  # 1. Check for required tools
  info "1. Checking for required tools (curl, tar, install)..."
  for tool in curl tar install; do
    if ! command -v "$tool" >/dev/null 2>&1; then
      error "'$tool' is required but not found. Please install it first."
    fi
  done
  success "All required tools are available."

  # 2. Determine OS and Architecture
  info "2. Detecting Operating System and Architecture..."
  OS_TYPE=$(uname -s | tr '[:upper:]' '[:lower:]')
  ARCH_TYPE=$(uname -m)

  case "$OS_TYPE" in
    linux)
      OS="linux"
      ;;
    darwin)
      OS="darwin" # For macOS
      ;;
    *)
      error "Unsupported OS: $OS_TYPE. This script supports Linux and macOS."
      ;;
  esac

  case "$ARCH_TYPE" in
    x86_64 | amd64)
      ARCH="amd64"
      ;;
    aarch64 | arm64)
      ARCH="arm64"
      ;;
    *)
      error "Unsupported architecture: $ARCH_TYPE. This script supports x86_64/amd64 and aarch64/arm64."
      ;;
  esac
  success "Detected OS: $OS, Architecture: $ARCH"

  # 3. Get the latest version tag from GitHub API
  info "3. Fetching the latest version of '$APP_NAME'..."
  # We use curl to hit the GitHub API for the latest release.
  # The 'jq' tool would be ideal here, but to avoid dependencies, we use grep/sed.
  LATEST_TAG_URL="https://api.github.com/repos/${GITHUB_REPO}/releases/latest"
  LATEST_VERSION=$(curl -sL "$LATEST_TAG_URL" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

  if [ -z "$LATEST_VERSION" ]; then
    error "Could not fetch the latest version tag from GitHub. Check repository name and network."
  fi
  success "Latest version is: $LATEST_VERSION"


  # 4. Construct download URL and download the release
  FILENAME="${APP_NAME}-${LATEST_VERSION}-${OS}-${ARCH}.tar.gz"
  DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${LATEST_VERSION}/${FILENAME}"

  info "4. Downloading from: $DOWNLOAD_URL"

  # Create a temporary directory for the download
  TMP_DIR=$(mktemp -d)
  if [ $? -ne 0 ]; then
    error "Failed to create a temporary directory."
  fi

  # Clean up the temporary directory on script exit
  trap 'rm -rf "$TMP_DIR"' EXIT

  # Download the file
  curl --progress-bar -L "$DOWNLOAD_URL" -o "$TMP_DIR/$FILENAME"
  if [ $? -ne 0 ]; then
    error "Download failed. Please check the URL and your network connection."
  fi
  success "Download complete."


  # 5. Install the binary
  info "5. Installing '$APP_NAME' to $INSTALL_DIR..."

  # Check for write permissions to the install directory
  SUDO_CMD=""
  if ! [ -w "$INSTALL_DIR" ]; then
      warn "Write permission to $INSTALL_DIR is required."
      # Check if sudo is available
      if command -v sudo >/dev/null 2>&1; then
          info "Attempting to use sudo..."
          SUDO_CMD="sudo"
      else
          error "Cannot write to $INSTALL_DIR. Please run this script as root or with sudo."
      fi
  fi
  
  # Extract the binary and install it
  cd "$TMP_DIR"

  tar -xzf "$FILENAME"
  
  # Find the binary in the extracted files
  if [ ! -f "$APP_NAME" ]; then
      error "Could not find '$APP_NAME' binary in the downloaded archive."
  fi
  
  # Use the 'install' command which handles permissions correctly
  $SUDO_CMD install -m 755 "$APP_NAME" "$INSTALL_DIR/$APP_NAME"
  if [ $? -ne 0 ]; then
      error "Installation failed. Could not move binary to $INSTALL_DIR."
  fi
  
  # 6. Final verification
  if ! command -v "$APP_NAME" >/dev/null 2>&1; then
      error "Installation failed. '$APP_NAME' command not found in PATH."
  fi

  # --- Final Message ---
  echo ""
  success "'$APP_NAME' ($LATEST_VERSION) was installed successfully to $INSTALL_DIR/$APP_NAME"
  info "You can now run the '$APP_NAME' command from your terminal."
  info "To get started, try running: $APP_NAME --help"
}

# Execute the main function
main
