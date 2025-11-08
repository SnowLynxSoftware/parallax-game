# Setup Docker

## Install Deps

```
sudo apt install -y apt-transport-https ca-certificates curl software-properties-common gnupg lsb-release
```

## Add docker's official GPG Key

```
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
```

## Add the docker repository

```
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" \
  | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
```

## Install Docker

```
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
sudo docker version
sudo docker info
```

## Enable and Start Docker Service

```
sudo systemctl enable --now docker
sudo systemctl status docker

```

## Setup Docker Credentials

### Log in to GHCR using Docker

```
echo "YOUR_GITHUB_PAT" | docker login ghcr.io -u YOUR_GITHUB_USERNAME --password-stdin
```

### Build the docker image

```
docker build -t ghcr.io/snowlynxsoftware/parallax-game:latest .
```

### Push the image to the GHCR

```
docker push ghcr.io/snowlynxsoftware/parallax-game:latest
```

### Run the container

```
docker run -d \
  --name parallax-game-2025-09-30 \
  --env-file /root/deployments/.env.production \
  -p 3000:3000 \
  ghcr.io/snowlynxsoftware/parallax-game:latest
```
