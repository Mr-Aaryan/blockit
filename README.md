#blockit - Website Blocker CLI

blockit 🛡️
===========

**blockit** is a simple Go-based CLI tool to block/unblock websites by modifying your `/etc/hosts` file. It stores the blocked website list in a SQLite database at `/etc/blockit/data.db`.

📦 Installation
---------------

### 1\. Clone the repository

    git clone https://github.com/Mr-Aaryan/blockit.git
    cd blockit
    

### 2\. Build and install

    go build -o ~/go/bin/blockit
    # Or install globally (requires $GOBIN in PATH)
    go install .
    

Make sure `~/go/bin` is in your `$PATH`.

### 3\. Create required directory and database (run with sudo)

    sudo mkdir -p /etc/blockit
    sudo touch /etc/blockit/data.db
    

⚙️ Usage
--------

Run with `sudo` to allow editing `/etc/hosts`:

### Block a site

    sudo blockit block facebook.com

### Unblock a site

    sudo blockit unblock facebook.com

### Show all blocked sites

    sudo blockit list

🗂️ Data
--------

All blocked website entries are tracked in:

    /etc/blockit/data.db

Make sure your app has permission to read/write this file (usually run as `sudo`).

🧑‍💻 Requirements
------------------

*   Go 1.18+
*   SQLite (via Go package)
*   Linux (tested on Arch)

📝 License
----------

MIT — see [LICENSE](./LICENSE)