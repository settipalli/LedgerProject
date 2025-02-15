# Ledger Project - Double-Entry Ledger System

A double-entry accounting system built in Go, designed to handle financial transactions with strict consistency
guarantees and proper currency handling.


## Features

- Support double-entry accounting system
- Thread-safe transaction processing
- ISO 4217 currency validation
- RESTful API interface
- Configurable environments (Development, Test, Production)
- Decimal precision handling for monetary values
- Transaction history tracking
- Real-time balance calculation


## Architecture

The project follows a clean architecture pattern with the following components:

- **API Layer** (`/api`): HTTP handlers and routing
- **Ledger Layer** (`/ledger`): Core business logic and transaction processing
- **Models** (`/models`): Domain entities and data structures
- **Services** (`/services`): Supporting services like currency validation
- **Config** (`/config`): Environment-specific configurations


## Getting Started

### Prerequisites

- Go 1.18 or higher
- Git

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/ledgerproject.git
cd ledgerproject
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
go run main.go
```

The server will start on port 8080 by default.


## API Endpoints

### Create Account
```bash
POST /accounts
```
Creates a new account in the ledger system.

Example request:
```json
{
    "id": "1001",
    "name": "Cash Account",
    "type": "asset",
    "currency": "USD",
    "balance": {
        "amount": "0",
        "currency": "USD"
    }
}
```

### Record Transaction
```bash
POST /transactions
```
Records a new transaction between two accounts.

Example request:
```json
{
    "id": "tx001",
    "description": "Initial deposit",
    "debit_account": "1001",
    "credit_account": "2001",
    "money": {
        "amount": "1000.00",
        "currency": "USD"
    }
}
```

### Get Balance
```bash
GET /accounts/{accountId}/balance
```
Retrieves the current balance for an account.

### Get Transaction History
```bash
GET /accounts/{accountId}/history
```
Retrieves the transaction history for an account.


## Technical Details

### Concurrency Handling
- Uses `sync.RWMutex` for thread-safe operations
- Implements proper locking mechanisms for account creation and transaction processing
- Supports concurrent read operations for balance queries

### Currency Handling
- Uses `decimal.Decimal` for precise monetary calculations
- Validates currencies against ISO 4217 standards
- Prevents mixed-currency transactions

### Transaction Consistency
- Enforces double-entry accounting principles
- Validates account existence before transactions
- Ensures currency matching between accounts and transactions
- Maintains transaction history with timestamps

## Configuration

The system supports three environments:
- Development (`:8080`)
- Test (`:8081`)
- Production (`:80`)

Environment-specific configurations are managed through the `config` package.

The environment selection is done at startup time, and the appropriate configuration is provided to the application.

You can run the application with one of the pre-supported environment settings by specifying the `APP_ENV` variable as
below:

```bash
# For development
export APP_ENV=dev
# or
APP_ENV=dev go run main.go

# For testing
export APP_ENV=test
# or
APP_ENV=test go run main.go

# For production
export APP_ENV=prod
# or
APP_ENV=prod go run main.go
```

## Dependencies

- `github.com/gorilla/mux`: HTTP routing
- `github.com/shopspring/decimal`: Precise decimal calculations
- `go.uber.org/fx`: Dependency injection
- Standard Go libraries

## Error Handling

The system provides detailed error messages for:
- Invalid account creation attempts
- Non-existent accounts
- Currency mismatches
- Invalid transaction amounts
- Missing required fields

## Best Practices

When using this ledger system:

1. Always create accounts before attempting transactions
2. Ensure matching currencies for transactions
3. Use proper account types (asset, liability, equity, revenue, expense)
4. Monitor transaction history for audit purposes
5. Handle errors appropriately in your application

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request


## Complete Workflow Example

The following commands illustrate a standard ledger usage workflow using CURL while adhering to the fundamental
principles of double-entry accounting:

- Each transaction has equal and opposite effects on two accounts
- Asset and Expense accounts increase with debits
- Liability, Equity, and Revenue accounts increase with credits
- The sum of all debits equals the sum of all credits
- All transactions maintain currency consistency

### 1. Create Required Accounts

First, create accounts for each side of our transactions:

```bash
# Create Asset Account (Cash)
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{
    "id": "1001",
    "name": "Cash Account",
    "type": "asset",
    "currency": "USD",
    "balance": {
      "amount": "15000.00",
      "currency": "USD"
    }
  }'

# Create Liability Account (Loan)
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{
    "id": "2001",
    "name": "Loan Account",
    "type": "liability",
    "currency": "USD",
    "balance": {
      "amount": "0",
      "currency": "USD"
    }
  }'

# Create Revenue Account
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{
    "id": "4001",
    "name": "Sales Revenue",
    "type": "revenue",
    "currency": "USD",
    "balance": {
      "amount": "0",
      "currency": "USD"
    }
  }'

# Create Expense Account
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{
    "id": "5001",
    "name": "Operating Expenses",
    "type": "expense",
    "currency": "USD",
    "balance": {
      "amount": "2000.00",
      "currency": "USD"
    }
  }'
```

### 2. Record Some Transactions

```bash
# Record initial loan of $10,000
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "id": "tx001",
    "description": "Initial bank loan",
    "debit_account": "1001",
    "credit_account": "2001",
    "amount": {
      "amount": "10000.00",
      "currency": "USD"
    }
  }'

# Record revenue of $5,000
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "id": "tx002",
    "description": "Client payment for services",
    "debit_account": "1001",
    "credit_account": "4001",
    "amount": {
      "amount": "5000.00",
      "currency": "USD"
    }
  }'

# Record expense of $2,000
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "id": "tx003",
    "description": "Office rent payment",
    "debit_account": "5001",
    "credit_account": "1001",
    "amount": {
      "amount": "2000.00",
      "currency": "USD"
    }
  }'
```

### 3. Check Account Balances

```bash
# Check Cash Account Balance (Should be 13,000)
curl -X GET http://localhost:8080/accounts/1001/balance

# Check Loan Account Balance (Should be 10,000)
curl -X GET http://localhost:8080/accounts/2001/balance

# Check Revenue Account Balance (Should be 5,000)
curl -X GET http://localhost:8080/accounts/4001/balance

# Check Expense Account Balance (Should be 0.00)
curl -X GET http://localhost:8080/accounts/5001/balance
```

### 4. View Transaction History

```bash
# Get all transactions for Cash Account
curl -X GET http://localhost:8080/accounts/1001/history
```

Expected response:
```json
[
  {
    "id": "tx001",
    "datetime": "2025-02-14T21:52:33.5428Z",
    "description": "Initial bank loan",
    "debit_account": "1001",
    "credit_account": "2001",
    "amount": {
      "amount": "10000",
      "currency": "USD"
    }
  },
  {
    "id": "tx002",
    "datetime": "2025-02-14T21:52:38.176906Z",
    "description": "Client payment for services",
    "debit_account": "1001",
    "credit_account": "4001",
    "amount": {
      "amount": "5000",
      "currency": "USD"
    }
  },
  {
    "id": "tx003",
    "datetime": "2025-02-14T21:52:44.137485Z",
    "description": "Office rent payment",
    "debit_account": "5001",
    "credit_account": "1001",
    "amount": {
      "amount": "2000",
      "currency": "USD"
    }
  }
]
```

## Development: follow the below guidelines to introduce a new API.
### Example: Implementing an API to list all active accounts

1. First, add the interface method in `ledger/interfaces.go`:

```go
type LedgerService interface {
    // ... existing methods ...
    GetAllAccounts() []models.Account
}
```

2. Implement the method in `ledger/ledger.go`:

```go
func (l *ledger) GetAllAccounts() []models.Account {
    l.mu.RLock()
    defer l.mu.RUnlock()

    accounts := make([]models.Account, 0, len(l.accounts))
    for _, account := range l.accounts {
        accounts = append(accounts, *account)
    }
    return accounts
}
```

3. Add the handler in `api/handlers.go`:

```go
func (s *Server) ListAccountsHandler(w http.ResponseWriter, r *http.Request) {
    accounts := s.ledger.GetAllAccounts()
    
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(accounts); err != nil {
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
        return
    }
}
```

4. Register the new route in `api/server.go`:

```go
func (s *Server) setupRoutes() {
    // ... existing routes ...
    s.router.HandleFunc("/accounts", s.ListAccountsHandler).Methods("GET")
    // Note: The POST route for creating accounts uses the same path but different HTTP method
}
```

Now you can list all accounts by making a GET request to `/accounts`. The response will be a JSON array of accounts.

Example usage with curl:
```bash
curl -X GET http://localhost:8080/accounts
```

Example response:
```json
[
    {
        "id": "1001",
        "name": "Cash Account",
        "balance": {
            "amount": "1000.00",
            "currency": "USD"
        },
        "type": "asset",
        "currency": "USD"
    },
    {
        "id": "2001",
        "name": "Revenue Account",
        "balance": {
            "amount": "-1000.00",
            "currency": "USD"
        },
        "type": "revenue",
        "currency": "USD"
    }
]
```

## License

This project is licensed under the Apache version 2 - see the LICENSE file for details.