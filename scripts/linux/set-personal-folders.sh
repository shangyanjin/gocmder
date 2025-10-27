#!/bin/bash

#
# Personal Folders Redirection Tool for Ubuntu/Linux
#
# This script will:
# 1. Create target directories under specified root
# 2. Update user shell configuration files (.bashrc, .zshrc)
# 3. Create symbolic links or move existing directories
# 4. Update XDG user directories configuration
#

set -e

TARGET_ROOT="${1:-$HOME/Personal}"
NO_MOVE="${2:-false}"

# Colors for output
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}Personal Folders Redirection Tool${NC}"
echo -e "${CYAN}Target Root: $TARGET_ROOT${NC}"

# Expand ~ to home directory
TARGET_ROOT="${TARGET_ROOT/#\~/$HOME}"

# Create base directory
if [ ! -d "$TARGET_ROOT" ]; then
    echo "Creating base target directory: $TARGET_ROOT"
    mkdir -p "$TARGET_ROOT"
fi

# Define folders to redirect
declare -A FOLDERS=(
    ["Desktop"]="Desktop"
    ["Documents"]="Documents"
    ["Downloads"]="Downloads"
    ["Music"]="Music"
    ["Pictures"]="Pictures"
    ["Videos"]="Videos"
)

# XDG user dirs config file
XDG_CONFIG="$HOME/.config/user-dirs.dirs"

# Process each folder
for name in "${!FOLDERS[@]}"; do
    old_path="$HOME/$name"
    new_path="$TARGET_ROOT/$name"
    
    echo -e "\n${CYAN}Processing [$name]${NC}"
    echo "  Old Path: $old_path"
    echo "  New Path: $new_path"
    
    # Create target directory if missing
    if [ ! -d "$new_path" ]; then
        mkdir -p "$new_path"
    fi
    
    # Move files if source exists and not the same
    if [ -d "$old_path" ] && [ "$old_path" != "$new_path" ]; then
        if [ "$NO_MOVE" != "true" ]; then
            if [ "$(ls -A "$old_path")" ]; then
                echo "  Moving files..."
                mv "$old_path"/* "$new_path/" 2>/dev/null || true
                rmdir "$old_path" 2>/dev/null || true
            fi
        fi
        
        # Create symlink
        if [ ! -e "$old_path" ]; then
            ln -s "$new_path" "$old_path"
            echo "  Symlink created: $old_path -> $new_path"
        fi
    fi
done

# Update XDG user dirs configuration
if [ -f "$XDG_CONFIG" ]; then
    echo -e "\n${CYAN}Updating XDG user directories...${NC}"
    
    for name in "${!FOLDERS[@]}"; do
        new_path="$TARGET_ROOT/$name"
        xdg_key="XDG_${name^^}_DIR"
        
        if grep -q "^$xdg_key=" "$XDG_CONFIG"; then
            sed -i "s|^$xdg_key=.*|$xdg_key=\"$new_path\"|" "$XDG_CONFIG"
        else
            echo "$xdg_key=\"$new_path\"" >> "$XDG_CONFIG"
        fi
    done
    
    echo "  XDG configuration updated"
fi

echo -e "\n${GREEN}Completed. Please log out and log back in to apply changes.${NC}"
