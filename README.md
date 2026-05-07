# ACT Backup Server

This is basic FTP server used to sync files to disk or a USB volume to a Windows XP VM in UTM.

Unfortunately, Windows XP is old enough that the typical methods of mounting a disk didn't work for me.

Instead, FTP is still well supported so I'm using that.

## Setup

### 1. Get the latest release

Download the latest release and put in ~/bin and act-backup-server with +x

### 2. Load in launchd

Save this to `~/bin/com.act.act-backup-server.plist`.

```
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
 "http://www.apple.com/DTDs/PropertyList-1.0.dtd">

<plist version="1.0">
<dict>

    <key>Label</key>
    <string>com.act.act-backup-server</string>

    <key>ProgramArguments</key>
    <array>
        <string>/path/to/bin/act-backup-server</string>
        <string>/Volumes/ACT BACKUP</string>
    </array>

    <key>RunAtLoad</key>
    <true/>

    <key>KeepAlive</key>
    <true/>

    <key>StandardOutPath</key>
    <string>/tmp/act-backup-server.log</string>

    <key>StandardErrorPath</key>
    <string>/tmp/act-backup-server.err</string>

</dict>
</plist>
```

```sh
launchctl load com.sean.act-backup-server.plist
```

### 3. Unload if needed

```sh
launchctl unload com.sean.act-backup-server.plist
```