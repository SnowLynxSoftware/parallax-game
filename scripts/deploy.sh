#!/bin/bash
set -euo pipefail

echo "📦 Pulling latest image..."
docker pull ghcr.io/snowlynxsoftware/parallax-game:latest

# Identify existing container (if any)
OLD_CONTAINER=$(docker ps -aq --filter "name=parallax-game-")

if [ -n "$OLD_CONTAINER" ]; then
  echo "🛑 Stopping old container: $OLD_CONTAINER"
  docker stop "$OLD_CONTAINER"
else
  echo "ℹ️ No previous container found."
fi

# Start new container with timestamped name
NEW_CONTAINER="parallax-game-$(date +%Y-%m-%d-%H%M%S)"
echo "🚀 Starting new container: $NEW_CONTAINER"

docker run -d \
  --name "$NEW_CONTAINER" \
  --env-file /root/deployments/.env.production \
  -p 3000:3000 \
  ghcr.io/snowlynxsoftware/parallax-game:latest

echo "⏳ Waiting for container $NEW_CONTAINER to be running..."
if timeout 30 bash -c \
  'until [ "$(docker inspect -f "{{.State.Status}}" '"$NEW_CONTAINER"')" = "running" ]; do sleep 2; done'; then
  echo "✅ Container $NEW_CONTAINER is running."
else
  echo "❌ Container $NEW_CONTAINER failed to start within 30s."
  echo "📜 Logs from failed container:"
  docker logs "$NEW_CONTAINER" || true
  echo "♻️ Restarting previous container..."
  if [ -n "$OLD_CONTAINER" ]; then
    docker start "$OLD_CONTAINER"
  else
    echo "⚠️ No previous container to restart!"
  fi
  exit 1
fi

# If we got here, the new container is good
if [ -n "$OLD_CONTAINER" ]; then
  echo "🧹 Removing old container: $OLD_CONTAINER"
  docker rm "$OLD_CONTAINER"
fi

echo "✅ Deployment complete."
