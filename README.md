#GOLS

GOLS is a command-line tool for listing files and folders in a directory with additional features such as colored icons and summary statistics.

## Features

- List files and folders in a directory.
- Recursive listing with the `--deep` flag.
- Colored icons based on file extensions.
- Summary statistics including the total number of files, folders, and total size.

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/MyGoApp.git

    Navigate to the project directory:

    bash

cd MyGoApp

Build the executable:

bash

go build -o mygoapp

Optionally, move the executable to a directory in your system's PATH:

bash

    mv mygoapp /usr/local/bin

Usage

bash

mygoapp [directory] [--deep]

    If no directory is provided, the current directory will be used.
    Use the --deep flag for recursive listing.

Configuration

MyGoApp uses a configuration file (config.json) to define icons and colors for file extensions. Customize this file to match your preferences.

Example config.json:

json

{
  "icons": {
    ".go": "🐹",
    ".exe": "🚀",
    "folder": "📁"
  },
  "colors": {
    ".exe": "fgYellow"
  }
}

    Icons and colors are applied based on file extensions.
    The default configuration is used if no matching entry is found.

Examples

bash

# List files and folders in the current directory
mygoapp

# List files and folders in a specific directory
mygoapp /path/to/directory

# Recursive listing
mygoapp --deep

License
