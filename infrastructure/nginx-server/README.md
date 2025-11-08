# NGINX Server Setup

We are hosting an Ubuntu VM in Digital Ocean that runs NGINX that serves as our entry point to the application and our load balancer. The following are notes on how to set it up should you need to recreate the server and/or perform maintenance on it.

## Update the system and install NGINX

```
sudo apt update && sudo apt upgrade -y
sudo apt install -y software-properties-common curl nginx
```

## Check that it's running

```
sudo systemctl enable --now nginx
sudo systemctl status nginx
```

Test HTTP access:

```
curl -I http://localhost
```

You should get a 200 OK from the default NGINX page.

## Create a directory for your site config

```
sudo mkdir -p /etc/nginx/sites-available
sudo mkdir -p /etc/nginx/sites-enabled
```

## Create NGINX reverse proxy configuration

Let's assume your app server runs on port 3000 on the same machine for now. Later you can add multiple upstream servers.

```
sudo vim /etc/nginx/sites-available/parallax.conf
```

Paste this:

```
# Define upstream for future load balancing

upstream parallax_service {
    server 127.0.0.1:3000;
    # server 192.168.1.101:5000;
    # server 192.168.1.102:5000;
}

server {
    listen 80;
    server_name parallax.com www.parallax.com;

    # Reverse proxy
    location / {
        proxy_pass http://parallax_service;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Optional: increase proxy buffer sizes if needed
    proxy_buffers 16 16k;
    proxy_buffer_size 32k;

}
```

## Enable the site

```
sudo ln -s /etc/nginx/sites-available/parallax.conf /etc/nginx/sites-enabled/parallax.conf
```

Remove the default site if you don't need it:

```
sudo rm /etc/nginx/sites-enabled/default
```

## Test NGINX configuration

```
sudo nginx -t
```

You should see:

```
nginx: syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
```

## Reload NGINX

```
sudo systemctl reload nginx
```

## Install Certbot

```
sudo apt install -y certbot python3-certbot-nginx
```

## Obtain SSL Cert with Certbot

```
sudo certbot --nginx -d parallax.com -d www.parallax.com
```

## Verify NGINX HTTPS Configuration

```
sudo nginx -t
sudo systemctl reload nginx
systemctl list-timers | grep certbot
sudo certbot renew --dry-run # Renew Manually
```

## Enforce Strongger SSL Settings

```
sudo vim /etc/nginx/sites-available/parallax.conf
```

Add this to the server block:

```
ssl_prefer_server_ciphers on;
```

Then reload:

```
sudo systemctl reload nginx
```
