# Create a 1 GB swapfile on Ubuntu Server

Run these commands as a user with sudo (or as root). I keep each step short with the precise commands to copy/paste.

## Check current swap and free space

```
sudo swapon --show
free -h
df -h /
```

## Create a 1 GB file for swap

```
sudo fallocate -l 1G /swapfile
```

## Secure the file permissions

```
sudo chmod 600 /swapfile
sudo chown root:root /swapfile
```

## Set up the swap area and enable it

```
sudo mkswap /swapfile
sudo swapon /swapfile
```

## Make the swap permanent across reboots

This will only append the line if it isn’t already present.

```
grep -q '^/swapfile ' /etc/fstab || echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
```

## Verify the swap is active

```
sudo swapon --show
free -h
cat /proc/swaps
```

## Tuning

Set swappiness to 10 (less aggressive swapping):

```
sudo sysctl vm.swappiness=10
echo 'vm.swappiness=10' | sudo tee /etc/sysctl.d/99-swappiness.conf
```
