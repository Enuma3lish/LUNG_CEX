# LUNG CEX - Virtual Trading Platform on Solana

A full-stack virtual cryptocurrency exchange built with Go/Gin backend, React.js frontend, and Solana blockchain integration.

## Features

- **Virtual Trading**: Start with $10,000 USD virtual balance
- **Multiple Assets**: Trade BTC, ETH, SOL, USDC, USDT
- **Spot & Futures Trading**: Support for both spot and perpetual futures
- **Blockchain Integration**: All trades recorded on Solana blockchain
- **Real-time Portfolio**: Track your holdings and P&L
- **Redis Caching**: Fast data retrieval with intelligent caching
- **PostgreSQL Storage**: Persistent data storage
- **JWT Authentication**: Secure user authentication

## Tech Stack

### Backend
- **Go 1.21+** with Gin framework
- **PostgreSQL** for persistent storage
- **Redis** for caching
- **Solana Go SDK** for blockchain integration
- **JWT** for authentication
- **GORM** for ORM

### Frontend
- **React 18** with Hooks
- **Vite** for build tooling
- **TailwindCSS** for styling
- **React Query** for data fetching
- **React Router** for navigation
- **Recharts** for charts and graphs

### Blockchain
- **Solana** blockchain (Devnet)
- Custom program for trade recording

## Project Structure

```
LUNG_CEX/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go              # Main entry point
│   ├── internal/
│   │   ├── models/                  # Data models
│   │   ├── handlers/                # HTTP handlers
│   │   │   ├── auth.go              # Authentication
│   │   │   ├── trade.go             # Trading operations
│   │   │   └── portfolio.go         # Portfolio management
│   │   ├── middleware/              # Middleware (auth, etc.)
│   │   ├── services/                # Business logic
│   │   └── database/                # Database setup
│   ├── pkg/
│   │   ├── blockchain/              # Solana integration
│   │   ├── redis/                   # Redis client
│   │   └── utils/                   # Utilities
│   ├── .env.example                 # Environment variables template
│   └── go.mod                       # Go dependencies
├── frontend/
│   ├── src/
│   │   ├── components/              # React components
│   │   ├── pages/                   # Page components
│   │   │   ├── Login.jsx
│   │   │   ├── Register.jsx
│   │   │   ├── Dashboard.jsx
│   │   │   ├── Trading.jsx
│   │   │   ├── Portfolio.jsx
│   │   │   └── History.jsx
│   │   ├── contexts/                # React contexts
│   │   ├── services/                # API services
│   │   └── utils/                   # Utilities
│   ├── package.json
│   └── vite.config.js
├── solana-program/                  # Solana smart contracts
└── README.md
```

## Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- PostgreSQL 13 or higher
- Redis 6 or higher
- Solana CLI (optional, for blockchain deployment)

## Setup Instructions

### 1. Database Setup

#### PostgreSQL

```bash
# Create database
createdb lung_cex

# Or using psql
psql -U postgres
CREATE DATABASE lung_cex;
```

#### Redis

```bash
# Start Redis server
redis-server

# Or using Docker
docker run -d -p 6379:6379 redis:alpine
```

### 2. Backend Setup

```bash
# Navigate to backend directory
cd backend

# Copy environment file
cp .env.example .env

# Edit .env with your configuration
nano .env

# Install dependencies
go mod download
go mod tidy

# Run the server
go run cmd/api/main.go
```

The backend server will start on `http://localhost:8080`

#### Backend Environment Variables

```env
# Server Configuration
PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=lung_cex

# Redis Configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

# JWT Configuration
JWT_SECRET=your-secret-key-change-this-in-production

# Solana Configuration
SOLANA_RPC_URL=https://api.devnet.solana.com
SOLANA_PRIVATE_KEY=
```

### 3. Frontend Setup

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

The frontend will start on `http://localhost:3000`

### 4. Solana Setup (Optional)

For blockchain integration, you need a Solana keypair:

```bash
# Install Solana CLI
sh -c "$(curl -sSfL https://release.solana.com/stable/install)"

# Generate a new keypair
solana-keygen new --outfile ~/.config/solana/devnet.json

# Get your public key
solana-keygen pubkey ~/.config/solana/devnet.json

# Airdrop some SOL (Devnet only)
solana airdrop 2 --url devnet

# Export private key to base58 and add to .env
# Use the base58 encoded private key in SOLANA_PRIVATE_KEY
```

## API Endpoints

### Authentication

- `POST /api/register` - Register new user
- `POST /api/login` - Login user
- `GET /api/user/profile` - Get user profile (protected)

### Trading

- `POST /api/trade/buy` - Buy asset (protected)
- `POST /api/trade/sell` - Sell asset (protected)
- `GET /api/trades/history` - Get trade history (protected)

### Portfolio

- `GET /api/portfolio` - Get portfolio summary (protected)
- `GET /api/portfolio/holdings` - Get current holdings (protected)

## Usage

### 1. Register an Account

- Navigate to `http://localhost:3000/register`
- Create an account with email, username, and password
- You'll automatically receive $10,000 USD virtual balance

### 2. View Dashboard

- After login, you'll see your dashboard with:
  - Total portfolio value
  - Cash balance
  - Holdings overview
  - Quick actions

### 3. Trade Assets

- Go to the Trading page
- Select an asset (BTC, ETH, SOL, USDC, USDT)
- Choose Buy or Sell
- Enter quantity
- Execute trade
- Trade will be recorded on Solana blockchain

### 4. Monitor Portfolio

- View your portfolio allocation
- Check profit/loss
- See detailed holdings with current prices

### 5. Check History

- View all your trades
- See blockchain transaction signatures
- Click to view on Solana Explorer

## Caching Strategy

The application uses Redis for intelligent caching:

- **Portfolio Data**: Cached for 30 seconds
- **Holdings Data**: Cached for 30 seconds
- **Price Data**: Cached for 5 seconds

When data is fresh (within TTL), it's served from Redis. Expired data is fetched from PostgreSQL and re-cached.

## Blockchain Integration

Each trade is recorded on the Solana blockchain:

1. Trade executed in database
2. Transaction created with trade details
3. Signed with configured keypair
4. Sent to Solana network
5. Signature stored in database
6. Viewable on Solana Explorer

## Development

### Run Backend Tests

```bash
cd backend
go test ./...
```

### Run Frontend in Development

```bash
cd frontend
npm run dev
```

### Build for Production

#### Backend

```bash
cd backend
go build -o lung-cex cmd/api/main.go
./lung-cex
```

#### Frontend

```bash
cd frontend
npm run build
npm run preview
```

## Docker Setup (Optional)

Create a `docker-compose.yml` for easy deployment:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: lung_cex
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

Run with:
```bash
docker-compose up -d
```

## Security Notes

⚠️ **Important Security Considerations:**

1. Change `JWT_SECRET` in production
2. Use strong database passwords
3. Enable SSL for PostgreSQL in production
4. Use environment variables, never commit secrets
5. This is a VIRTUAL trading platform - no real money involved
6. Solana private key should be kept secure

## Troubleshooting

### Backend won't start

- Check PostgreSQL is running: `psql -U postgres -c "SELECT 1"`
- Check Redis is running: `redis-cli ping`
- Verify environment variables are set correctly

### Frontend can't connect to backend

- Ensure backend is running on port 8080
- Check CORS settings in backend
- Verify proxy settings in `vite.config.js`

### Blockchain transactions failing

- Ensure you have SOL in your devnet wallet
- Check Solana network status
- Verify RPC URL is correct
- Airdrop more SOL: `solana airdrop 2 --url devnet`

## Future Enhancements

- [ ] Real-time price updates using WebSockets
- [ ] Advanced order types (limit, stop-loss)
- [ ] Leverage trading for futures
- [ ] Social trading features
- [ ] Mobile app
- [ ] More trading pairs
- [ ] Technical analysis tools
- [ ] Trading bots/algorithms

## Contributing

This is a virtual trading platform for educational purposes. Contributions are welcome!

## License

MIT License

## Support

For issues and questions, please open an issue on GitHub.

---

**Disclaimer**: This is a virtual trading platform for educational purposes only. No real cryptocurrency or money is involved. All trades are simulated.
