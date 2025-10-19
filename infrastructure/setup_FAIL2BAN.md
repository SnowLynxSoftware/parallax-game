# Setup FAIL2BAN On Ubuntu Server

## Install Fail2Ban

```
sudo apt update
sudo apt install -y fail2ban
```

## Copy the default configuration (best practice)

Never edit jail.conf directly — use jail.local or a custom .local file.

```
sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local
```

## Configure the SSHD jail

Open the config file in your editor:

```
sudo vim /etc/fail2ban/jail.local
```

Find the [sshd] section and set the following:

```
[sshd]
enabled = true
port = ssh
filter = sshd
logpath = /var/log/auth.log
maxretry = 3
bantime = 1h
findtime = 10m

maxretry = 3 → lock out after 3 failed attempts

findtime = 10m → look back 10 minutes for failures

bantime = 1h → ban for 1 hour (you can change to 24h or -1 for permanent ban)
```

## Restart and enable Fail2Ban

```
sudo systemctl enable --now fail2ban
sudo systemctl restart fail2ban
```

## Check Fail2Ban status

```
sudo systemctl status fail2ban
sudo fail2ban-client status
sudo fail2ban-client status sshd
```

The last command shows active bans for SSHD.

## (Optional) Unban an IP

```
sudo fail2ban-client set sshd unbanip <IP_ADDRESS>
```
