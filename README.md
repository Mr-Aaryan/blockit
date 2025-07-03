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

    sudo blockit

Now you can select to block and unselect to unblock.

🗂️ Data
--------

All blocked website entries are tracked in:

    /etc/blockit/data.db

And `/etc/hosts` contains all the list of sites which are to be blocked and line number along with website name is added to `data.db` to perform block operation

### List of sites
```
# 127.0.0.1 www.reddit.com
# 127.0.0.1 reddit.com
# 127.0.0.1 www.linkedin.com
# 127.0.0.1 linkedin.com
# 127.0.0.1 x.com
# 127.0.0.1 www.x.com
# 127.0.0.1 twitter.com
# 127.0.0.1 www.twitter.com
# 127.0.0.1 www.facebook.com
# 127.0.0.1 facebook.com
# 127.0.0.1 netflix.com
# 127.0.0.1 www.netflix.com
# 127.0.0.1 twitch.tv
# 127.0.0.1 www.twitch.tv
# 127.0.0.1 youtube.com
# 127.0.0.1 www.youtube.com
# 127.0.0.1 instagram.com
# 127.0.0.1 www.instagram.com
# 127.0.0.1 tiktok.com
# 127.0.0.1 www.tiktok.com
# 127.0.0.1 pinterest.com
# 127.0.0.1 www.pinterest.com
# 127.0.0.1 snapchat.com
# 127.0.0.1 www.snapchat.com
# 127.0.0.1 discord.com
# 127.0.0.1 www.discord.com
# 127.0.0.1 roblox.com
# 127.0.0.1 www.roblox.com
# 127.0.0.1 steamcommunity.com
# 127.0.0.1 www.steamcommunity.com
# 127.0.0.1 amazon.com
# 127.0.0.1 www.amazon.com
# 127.0.0.1 ebay.com
# 127.0.0.1 www.ebay.com
# 127.0.0.1 aliexpress.com
# 127.0.0.1 www.aliexpress.com
```


Make sure your app has permission to read/write this file (usually run as `sudo`).

🧑‍💻 Requirements
------------------

*   Go 1.18+
*   SQLite (via Go package)
*   Linux (tested on Arch)

📝 License
----------

MIT — see [LICENSE](./LICENSE)
