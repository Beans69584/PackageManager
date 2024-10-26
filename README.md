# PackageManager

PackageManager is a simple CLI-based package manager written in Go. It allows users to install, uninstall, and manage software packages distributed as `.tar.gz` archives. It was only really created for personal use, it probably has a lot of flaws and isn't practical for most people.

## Features

- **Install Packages:** Extracts `.tar.gz` archives, creates symlinks for executables, and generates `.desktop` files for application launchers like Wofi.
- **Uninstall Packages:** Removes installed packages, symlinks, and corresponding `.desktop` files.
- **List Installed Packages:** Displays all currently installed packages with their details.

## Installation

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/Beans69584/PackageManager.git
   cd PackageManager

2. **Build the Package**

    ```bash
    go build -o PackageManager main.go
    ```

3. **Move Executable to Bin Dir**

    ```bash
    sudo mv PackageManager /usr/local/bin/
    ```
