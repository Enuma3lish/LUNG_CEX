# Quick Start Guide

Get LUNG CEX running in 5 minutes!

## Option 1: Using Make (Recommended)

### Step 1: Install Dependencies

```bash
make setup
```

### Step 2: Start Database Services

```bash
make docker-up
```

This will start PostgreSQL and Redis using Docker.

### Step 3: Configure Backend

```bash
cd backend
cp .env.example .env
```

The default values should work if you're using Docker.

### Step 4: Start Backend

```bash
# In one terminal
make backend
```

Wait for the message: "Server starting on port 8080"

### Step 5: Start Frontend

```bash
# In another terminal
make frontend
```

### Step 6: Access the Application

Open your browser and go to:
```
http://localhost:3000
```

## Option 2: Manual Setup

### Prerequisites

Install these first:
- Go 1.21+
- Node.js 18+
- Docker & Docker Compose

### Step 1: Start Database

```bash
docker-compose up -d
```

### Step 2: Setup Backend

```bash
cd backend
go mod download
cp .env.example .env
go run cmd/api/main.go
```

### Step 3: Setup Frontend

```bash
cd frontend
npm install
cp .env.example .env
npm run dev
```

### Step 4: Create Account

1. Go to http://localhost:3000
2. Click "Register here"
3. Fill in:
   - Email: your@email.com
   - Username: yourname
   - Password: password123
4. Click Register

You'll get $10,000 virtual USD automatically!

## First Trade

1. After login, click "Trading" in the navbar
2. Select an asset (e.g., BTC)
3. Click "Buy"
4. Enter quantity (e.g., 0.1)
5. Click "BUY BTC"

Your trade will be recorded on Solana blockchain!

## View Portfolio

Click "Portfolio" to see:
- Your total value
- All holdings
- Profit/Loss
- Asset allocation chart

## Troubleshooting

### Backend won't start

```bash
# Check if databases are running
docker ps

# Should see lung_cex_postgres and lung_cex_redis
```

### Frontend shows connection error

Make sure backend is running on port 8080:
```bash
curl http://localhost:8080/api/login
# Should return 405 Method Not Allowed (which means it's running)
```

### Database connection failed

```bash
# Check PostgreSQL
docker exec -it lung_cex_postgres psql -U postgres -c "SELECT 1"

# Check Redis
docker exec -it lung_cex_redis redis-cli ping
```

## Default Configuration

### Backend (Port 8080)
- Database: PostgreSQL on localhost:5432
- Cache: Redis on localhost:6379
- Blockchain: Solana Devnet

### Frontend (Port 3000)
- API: http://localhost:8080/api

## Next Steps

1. Explore the Trading page
2. Try buying different assets
3. Check your trade history
4. View your portfolio allocation
5. See trades on Solana Explorer (click "View on Explorer" in History)

## Need Help?

Check the main README.md for detailed documentation.

## Clean Up

To stop everything:

```bash
# Stop Docker containers
make docker-down

# Or manually
docker-compose down
```

To remove all data:

```bash
docker-compose down -v
```

Happy Trading! 🚀
