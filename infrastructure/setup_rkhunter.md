# Setup RKHUNTER

## Install rkhunter

```
sudo apt update
sudo apt install -y rkhunter mailutils
```

Choose LOCAL ONLY

## Update rkhunter database

```
sudo rkhunter --update
sudo rkhunter --propupd
```

The --propupd step records the current system state as a baseline.

## Configure rkhunter

Edit the config file:

```
sudo vim /etc/rkhunter.conf
```

Make sure these lines are set (uncomment or add them if missing):

```
WEB_CMD=""
```

## Create a log file

```
sudo touch /var/log/rkhunter.log
sudo chown root:root /var/log/rkhunter.log
sudo chmod 600 /var/log/rkhunter.log
```

## Create a daily cron job

Instead of relying on the default cron.d entry (sometimes disabled on Ubuntu), we’ll create our own:

```
sudo vim /etc/cron.daily/rkhunter-check
```

Paste this script:

```
#!/bin/bash

# Define log file
LOGFILE="/var/log/rkhunter.log"

# Run rkhunter cron job
/usr/bin/rkhunter --cronjob --report-warnings-only --skip-keypress > "$LOGFILE" 2>&1
```

Save & exit.

## Make the script executable

```
sudo chmod +x /etc/cron.daily/rkhunter-check
```

Ubuntu's /etc/cron.daily scripts run once a day via cron (typically early morning).

## Test the script manually

Run it once to confirm it works:

```
sudo /etc/cron.daily/rkhunter-check
```

Then check the log file:

```
cat /var/log/rkhunter.log
```
