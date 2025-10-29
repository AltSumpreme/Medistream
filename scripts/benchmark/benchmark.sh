#!/usr/bin/env bash
# =========================================
#  Medistream Load Test Script
# =========================================

BASE_URL=${1:-http://localhost:8080}
THREADS=${2:-8}
CONNECTIONS=${3:-200}
DURATION=${4:-30s}

echo "🚀 Medistream Load Test"
echo "Base URL: $BASE_URL"
echo "Threads: $THREADS | Connections: $CONNECTIONS | Duration: $DURATION"
echo "==============================================="

echo "📦 Testing /auth/signup ..."
wrk -t$THREADS -c$CONNECTIONS -d$DURATION -s ./signup_benchmark.lua $BASE_URL/auth/signup
echo "==============================================="

echo "🔑 Testing /auth/login ..."
wrk -t$THREADS -c$CONNECTIONS -d$DURATION -s ./login_benchmark.lua $BASE_URL/auth/login
echo "==============================================="
echo "✅ Benchmark completed"
