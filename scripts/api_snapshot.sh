#!/usr/bin/env bash
# Capture request/response formats for all SwapKit API endpoints.
# API key below is for local inspection only; production should use SWAPKIT_API_KEY env var.
# Override with: SWAPKIT_API_KEY=your_key ./scripts/api_snapshot.sh
set -e
BASE_URL="${SWAPKIT_BASE_URL:-https://api.swapkit.dev}"
API_KEY="${SWAPKIT_API_KEY:-488c6f98-b890-4ea6-9ac5-183966c547f1}"
SNAPSHOT_DIR="$(dirname "$0")/snapshots"
mkdir -p "$SNAPSHOT_DIR"

echo "=== GET $BASE_URL/providers ==="
echo "Request: GET (no body)"
curl -sS -X GET "$BASE_URL/providers" \
  -H "accept: application/json" \
  -H "x-api-key: $API_KEY" | jq . > "$SNAPSHOT_DIR/providers.json" 2>/dev/null || curl -sS -X GET "$BASE_URL/providers" -H "accept: application/json" -H "x-api-key: $API_KEY" > "$SNAPSHOT_DIR/providers.json"
echo "Response saved to $SNAPSHOT_DIR/providers.json"
echo ""

echo "=== GET $BASE_URL/tokens?provider=THORCHAIN ==="
echo "Request: GET ?provider=THORCHAIN"
curl -sS -X GET "$BASE_URL/tokens?provider=THORCHAIN" \
  -H "accept: application/json" \
  -H "x-api-key: $API_KEY" | jq . > "$SNAPSHOT_DIR/tokens.json" 2>/dev/null || curl -sS -X GET "$BASE_URL/tokens?provider=THORCHAIN" -H "accept: application/json" -H "x-api-key: $API_KEY" > "$SNAPSHOT_DIR/tokens.json"
echo "Response saved to $SNAPSHOT_DIR/tokens.json"
echo ""

echo "=== GET $BASE_URL/swapFrom?buyAsset=BTC.BTC ==="
echo "Request: GET ?buyAsset=BTC.BTC"
curl -sS -X GET "$BASE_URL/swapFrom?buyAsset=BTC.BTC" \
  -H "accept: application/json" \
  -H "x-api-key: $API_KEY" | jq . > "$SNAPSHOT_DIR/swapFrom.json" 2>/dev/null || curl -sS -X GET "$BASE_URL/swapFrom?buyAsset=BTC.BTC" -H "accept: application/json" -H "x-api-key: $API_KEY" > "$SNAPSHOT_DIR/swapFrom.json"
echo "Response saved to $SNAPSHOT_DIR/swapFrom.json"
echo ""

echo "=== GET $BASE_URL/swapTo?sellAsset=ETH.ETH ==="
echo "Request: GET ?sellAsset=ETH.ETH"
curl -sS -X GET "$BASE_URL/swapTo?sellAsset=ETH.ETH" \
  -H "accept: application/json" \
  -H "x-api-key: $API_KEY" | jq . > "$SNAPSHOT_DIR/swapTo.json" 2>/dev/null || curl -sS -X GET "$BASE_URL/swapTo?sellAsset=ETH.ETH" -H "accept: application/json" -H "x-api-key: $API_KEY" > "$SNAPSHOT_DIR/swapTo.json"
echo "Response saved to $SNAPSHOT_DIR/swapTo.json"
echo ""

echo "=== POST $BASE_URL/v3/quote ==="
echo 'Request body: {"sellAsset":"ETH.ETH","buyAsset":"SOL.SOL","sellAmount":"0.01","slippage":2}'
curl -sS -X POST "$BASE_URL/v3/quote" \
  -H "Content-Type: application/json" \
  -H "accept: application/json" \
  -H "x-api-key: $API_KEY" \
  -d '{"sellAsset":"ETH.ETH","buyAsset":"SOL.SOL","sellAmount":"0.01","slippage":2}' | jq . > "$SNAPSHOT_DIR/quote.json" 2>/dev/null || curl -sS -X POST "$BASE_URL/v3/quote" -H "Content-Type: application/json" -H "accept: application/json" -H "x-api-key: $API_KEY" -d '{"sellAsset":"ETH.ETH","buyAsset":"SOL.SOL","sellAmount":"0.01","slippage":2}' > "$SNAPSHOT_DIR/quote.json"
echo "Response saved to $SNAPSHOT_DIR/quote.json"
echo ""

echo "=== POST $BASE_URL/track ==="
echo 'Request body: {"hash":"0x0","chainId":"1"}'
curl -sS -X POST "$BASE_URL/track" \
  -H "Content-Type: application/json" \
  -H "accept: application/json" \
  -H "x-api-key: $API_KEY" \
  -d '{"hash":"0x0","chainId":"1"}' | jq . > "$SNAPSHOT_DIR/track.json" 2>/dev/null || curl -sS -X POST "$BASE_URL/track" -H "Content-Type: application/json" -H "accept: application/json" -H "x-api-key: $API_KEY" -d '{"hash":"0x0","chainId":"1"}' > "$SNAPSHOT_DIR/track.json"
echo "Response saved to $SNAPSHOT_DIR/track.json"
echo ""

echo "=== POST $BASE_URL/price ==="
echo 'Request body: {"tokens":[{"identifier":"ETH.ETH"}],"metadata":false}'
curl -sS -X POST "$BASE_URL/price" \
  -H "Content-Type: application/json" \
  -H "accept: application/json" \
  -H "x-api-key: $API_KEY" \
  -d '{"tokens":[{"identifier":"ETH.ETH"}],"metadata":false}' | jq . > "$SNAPSHOT_DIR/price.json" 2>/dev/null || curl -sS -X POST "$BASE_URL/price" -H "Content-Type: application/json" -H "accept: application/json" -H "x-api-key: $API_KEY" -d '{"tokens":[{"identifier":"ETH.ETH"}],"metadata":false}' > "$SNAPSHOT_DIR/price.json"
echo "Response saved to $SNAPSHOT_DIR/price.json"
echo ""

echo "Done. Snapshots in $SNAPSHOT_DIR/"
