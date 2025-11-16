# API Documentation

Base URL: `http://localhost:8080/api`

## Authentication

All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

---

## Public Endpoints

### Register User

**POST** `/register`

Create a new user account with $10,000 starting balance.

**Request Body:**
```json
{
  "email": "user@example.com",
  "username": "johndoe",
  "password": "password123"
}
```

**Response:** `201 Created`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "username": "johndoe",
    "balance": 10000.00,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Errors:**
- `400` - Invalid request data
- `409` - User already exists

---

### Login

**POST** `/login`

Authenticate a user and receive a JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:** `200 OK`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "username": "johndoe",
    "balance": 9500.00,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Errors:**
- `400` - Invalid request data
- `401` - Invalid credentials

---

## Protected Endpoints

### Get User Profile

**GET** `/user/profile`

Get the current user's profile information.

**Headers:**
```
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "username": "johndoe",
  "balance": 9500.00,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**Errors:**
- `401` - Unauthorized
- `404` - User not found

---

## Trading Endpoints

### Buy Asset

**POST** `/trade/buy`

Purchase an asset at the current market price.

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "asset_symbol": "BTC",
  "quantity": 0.1,
  "price": 45000.00
}
```

**Response:** `200 OK`
```json
{
  "message": "Trade executed successfully",
  "trade": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "asset_id": "770e8400-e29b-41d4-a716-446655440002",
    "trade_type": "BUY",
    "quantity": 0.1,
    "price": 45000.00,
    "total_amount": 4500.00,
    "solana_signature": "5J8t...",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "solana_signature": "5J8tK3pVqG8Lq...",
  "remaining_balance": 5500.00
}
```

**Errors:**
- `400` - Invalid request or insufficient balance
- `401` - Unauthorized
- `404` - Asset not found

---

### Sell Asset

**POST** `/trade/sell`

Sell an asset at the current market price.

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "asset_symbol": "BTC",
  "quantity": 0.05,
  "price": 46000.00
}
```

**Response:** `200 OK`
```json
{
  "message": "Trade executed successfully",
  "trade": {
    "id": "660e8400-e29b-41d4-a716-446655440003",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "asset_id": "770e8400-e29b-41d4-a716-446655440002",
    "trade_type": "SELL",
    "quantity": 0.05,
    "price": 46000.00,
    "total_amount": 2300.00,
    "solana_signature": "6K9u...",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "solana_signature": "6K9uL4qWrH9Mr...",
  "remaining_balance": 7800.00
}
```

**Errors:**
- `400` - Invalid request or insufficient quantity
- `401` - Unauthorized
- `404` - Asset not found

---

### Get Trade History

**GET** `/trades/history`

Get the user's trade history (last 100 trades).

**Headers:**
```
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
[
  {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "asset_id": "770e8400-e29b-41d4-a716-446655440002",
    "trade_type": "BUY",
    "quantity": 0.1,
    "price": 45000.00,
    "total_amount": 4500.00,
    "solana_signature": "5J8tK3pVqG8Lq...",
    "created_at": "2024-01-01T00:00:00Z",
    "asset": {
      "id": "770e8400-e29b-41d4-a716-446655440002",
      "symbol": "BTC",
      "name": "Bitcoin",
      "asset_type": "SPOT"
    }
  }
]
```

**Errors:**
- `401` - Unauthorized
- `500` - Server error

---

## Portfolio Endpoints

### Get Portfolio

**GET** `/portfolio`

Get the user's complete portfolio with current values and P&L.

**Headers:**
```
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "total_value": 10250.50,
  "cash": 5500.00,
  "pnl": 250.50,
  "holdings": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440004",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "asset_id": "770e8400-e29b-41d4-a716-446655440002",
      "quantity": 0.1,
      "avg_price": 45000.00,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "asset": {
        "id": "770e8400-e29b-41d4-a716-446655440002",
        "symbol": "BTC",
        "name": "Bitcoin",
        "asset_type": "SPOT"
      },
      "current_price": 47500.00,
      "value": 4750.50,
      "pnl": 250.50,
      "pnl_percent": 5.57
    }
  ]
}
```

**Errors:**
- `401` - Unauthorized
- `404` - User not found

---

### Get Holdings

**GET** `/portfolio/holdings`

Get the user's current asset holdings.

**Headers:**
```
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
[
  {
    "id": "880e8400-e29b-41d4-a716-446655440004",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "asset_id": "770e8400-e29b-41d4-a716-446655440002",
    "quantity": 0.1,
    "avg_price": 45000.00,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "asset": {
      "id": "770e8400-e29b-41d4-a716-446655440002",
      "symbol": "BTC",
      "name": "Bitcoin",
      "asset_type": "SPOT"
    }
  }
]
```

**Errors:**
- `401` - Unauthorized
- `500` - Server error

---

## Available Assets

The platform supports the following assets:

### Spot Assets
- **BTC** - Bitcoin (~$45,000)
- **ETH** - Ethereum (~$2,500)
- **SOL** - Solana (~$100)
- **USDC** - USD Coin ($1.00)
- **USDT** - Tether USD ($1.00)

### Futures Assets
- **BTC-PERP** - Bitcoin Perpetual Futures (~$45,000)
- **ETH-PERP** - Ethereum Perpetual Futures (~$2,500)
- **SOL-PERP** - Solana Perpetual Futures (~$100)

---

## Error Responses

All errors follow this format:

```json
{
  "error": "Error message description"
}
```

### Common HTTP Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request (invalid input)
- `401` - Unauthorized (missing or invalid token)
- `404` - Not Found
- `409` - Conflict (resource already exists)
- `500` - Internal Server Error

---

## Rate Limiting

Currently, there are no rate limits on API endpoints.

---

## Caching

Portfolio and holdings endpoints are cached in Redis for 30 seconds. Fresh data is automatically fetched when cache expires.

---

## Blockchain Integration

All trades are recorded on the Solana blockchain (Devnet). The `solana_signature` field in trade responses contains the transaction signature, which can be viewed on:

https://explorer.solana.com/tx/{signature}?cluster=devnet

---

## Example Usage with cURL

### Register
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","username":"testuser","password":"password123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### Buy Asset
```bash
curl -X POST http://localhost:8080/api/trade/buy \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"asset_symbol":"BTC","quantity":0.1,"price":45000}'
```

### Get Portfolio
```bash
curl -X GET http://localhost:8080/api/portfolio \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## WebSocket Support

WebSocket support for real-time price updates is planned for future releases.
