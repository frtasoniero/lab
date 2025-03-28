#!/bin/bash

echo "Backing up Neovim config and related directories..."

backup_folder() {
  local target="$1"
  local backup="${target}.bak"

  if [ -d "$target" ]; then
    if [ -d "$backup" ]; then
      echo "Removing existing backup: $backup"
      rm -rf "$backup"
    fi
    echo "Backing up $target -> $backup"
    mv "$target" "$backup"
  else
    echo "Directory $target does not exist, skipping..."
  fi
}

# Required
backup_folder "$HOME/.config/nvim"

# Optional but recommended
backup_folder "$HOME/.local/share/nvim"
backup_folder "$HOME/.local/state/nvim"
backup_folder "$HOME/.cache/nvim"

echo "Done. All backups created with .bak suffix."

