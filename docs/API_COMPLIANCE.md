# SwapKit API compliance

This document confirms that the SDK request/response types and API usage match the [SwapKit API](https://docs.swapkit.dev/).

## Endpoints and usage

| Endpoint | Method | SDK usage | Match |
|----------|--------|-----------|-------|
| `/providers` | GET | `GET /providers`, no query/body, `x-api-key` + `accept: application/json` | Yes |
| `/tokens` | GET | `GET /tokens?provider=<provider>` | Yes |
| `/swapFrom` | GET | `GET /swapFrom?buyAsset=<identifier>` (tokens you can sell to get buyAsset) | Yes |
| `/swapTo` | GET | `GET /swapTo?sellAsset=<identifier>` (tokens you can buy with sellAsset) | Yes |
| `/v3/quote` | POST | `POST /v3/quote`, JSON body | Yes |
| `/v3/swap` | POST | `POST /v3/swap`, JSON body | Yes |
| `/track` | POST | `POST /track`, JSON body | Yes |
| `/price` | POST | `POST /price`, JSON body | Yes |

Base URL: `https://api.swapkit.dev`. All requests send header `x-api-key`; POST requests send `Content-Type: application/json` and `accept: application/json`.

---

## Request/response formats

### GET /providers

- **Request:** No parameters.
- **Response:** Array of provider objects. SDK `types.Provider`: `name`, `provider`, `displayName`, `displayNameLong`, `keywords`, `count`, `logoURI`, `url`, `enabledChainIds`, `supportedActions`, `supportedChainIds` — all match docs and live snapshot.

### GET /tokens

- **Request:** Query `provider` (required in practice).
- **Response:** Single object. SDK `types.TokenList`: `provider`, `name`, `timestamp`, `version` (major/minor/patch), `keywords`, `count`, `tokens`. SDK `types.Token`: `chain`, `chainId`, `address`, `ticker`, `identifier`, `symbol`, `name`, `decimals`, `logoURI`, `coingeckoId` — match docs and sample.

### GET /swapFrom

- **Request:** Query `buyAsset` = token identifier (e.g. `BTC.BTC`). Returns list of token identifiers that can be **sold** to receive that asset.
- **Response:** Array of strings (token identifiers). SDK decodes as `[]string`. Match.

### GET /swapTo

- **Request:** Query `sellAsset` = token identifier. Returns list of token identifiers that can be **bought** with that asset.
- **Response:** Array of strings. SDK decodes as `[]string`. Match.

### POST /v3/quote

- **Request body:** `sellAsset`, `buyAsset`, `sellAmount` (string), `slippage` (number). Optional: `providers` (array of strings). SDK `QuoteRequest` matches.
- **Response:** `quoteId`, `routes[]`, optional `providerErrors[]`, optional `error`. Route: `routeId`, `provider`, `providers[]`, `sellAsset`, `sellAmount`, `buyAsset`, `expectedBuyAmount`, `expectedBuyAmountMaxSlippage`, `inboundAmount`, `outboundAmount`, `expiration`, `estimatedTime` (inbound, swap, outbound, total — numbers in seconds), `fees` (array; SDK `json.RawMessage`), `totalSlippage`, `totalSlippageBps`, `legs`, `warnings`, `meta`, `nextActions[]`. SDK types updated to match live API. `QuoteError`: `provider`, `errorCode`, `message`. Match.

### POST /v3/swap

- **Request body:** `routeId`, `sourceAddress`, `destinationAddress`. Optional: `disableBalanceCheck`, `disableBuildTx`, `overrideSlippage`, `disableEstimate`. SDK `SwapRequest` uses these exact JSON names (camelCase).
- **Response:** Quote route fields plus `targetAddress`, `inboundAddress`, `tx` (string or object), `meta.txType`. SDK `SwapResponse` embeds `QuoteRoute` and adds these; `tx` as `json.RawMessage` handles both string (e.g. PSBT) and object (e.g. EVM). Match.

### POST /track

- **Request body:** Either `hash` + `chainId`, or `depositAddress` (NEAR). SDK `TrackRequest`: `hash`, `chainId`, `depositAddress` with `omitempty`. Match.
- **Response:** `chainId`, `hash`, `block`, `type`, `status`, `trackingStatus`, `fromAsset`, `fromAmount`, `fromAddress`, `toAsset`, `toAmount`, `toAddress`, `finalisedAt`, `meta` (provider, providerAction, images), `payload`, `legs[]`. SDK `TrackResponse` and `TrackMeta` (including `images` as `map[string]string`) match the doc example.

### POST /price

- **Request body:** `tokens`: array of `{ "identifier": "..." }`, optional `metadata`: boolean. SDK `PriceRequest` / `PriceTokenInput` match.
- **Response:** Array of objects: `identifier`, `provider`, `price_usd`, `timestamp`, optional `cg` (CoinGecko). SDK `PriceResult`: `identifier`, `provider`, `price_usd`, `timestamp`, `cg`. `CoingeckoPrice` uses snake_case tags (`market_cap`, `total_volume`, `price_change_24h_usd`, `price_change_percentage_24h_usd`, `sparkline_in_7d`). Match.

---

## Error responses

4xx/5xx bodies: `message`, `error` (code), optional `data`. SDK `client.APIError` and `parseAPIError` use the same shape. Documented codes (`noRoutesFound`, `insufficientBalance`, `apiKeyInvalid`, etc.) are unwrapped via `Unwrap()` and sentinel errors. Match.

---

## Summary

- **Paths and methods:** All eight endpoints use the correct method and path (`/providers`, `/tokens`, `/swapFrom`, `/swapTo`, `/v3/quote`, `/v3/swap`, `/track`, `/price`).
- **Query parameters:** `provider`, `buyAsset`, `sellAsset` used as in the docs.
- **Request bodies:** JSON field names and structure match the documented request schemas.
- **Response types:** Structs and `json` tags align with the documented and sample response formats; optional fields use `omitempty` where appropriate.
- **Headers:** `x-api-key`, `accept: application/json`, and `Content-Type: application/json` (for POST with body) are set as required.

The SDK is aligned with the SwapKit API as documented at https://docs.swapkit.dev/.
