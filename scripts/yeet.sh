#!/bin/bash
set -e

echo "🧹 Cleaning up isolated interview environment..."

# Stop and remove any docker containers & volumes
if command -v docker &> /dev/null; then
  docker compose down -v 2>/dev/null || true
  docker rm -f interview_postgres 2>/dev/null || true
fi

echo "✅ Environment teardown complete."
echo "To completely delete the project directory run:"
echo "rm -rf ~/Desktop/interview-prep-app"
