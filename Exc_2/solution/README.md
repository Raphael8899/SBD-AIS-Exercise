# Big Data Exercise 2 - REST API Solution

## What This Is

This is my solution for Exercise 2 in the Software Architecture for Big Data course. The assignment was to complete a REST API server written in Golang that manages a drink ordering system.

## What I Did

I had to fill in all the TODO comments in the skeleton code to make a working REST server. Here's what I implemented:

### 1. Data Models
- **Drink model**: Added Name, Price, and Description fields
- **Order model**: Added Amount and CreatedAt fields  
- Made sure all JSON tags use snake_case like the instructions said

### 2. Database Handler
- Set up test data with 3 drinks (Beer, Spritzer, Coffee) matching the prices from the slides
- Implemented the `GetTotalledOrders()` function that counts how many of each drink was ordered
- Completed the `AddOrder()` function to save new orders

### 3. REST API Endpoints
I implemented all 4 required endpoints:

- **GET /api/menu** - Returns the list of available drinks
- **GET /api/order/all** - Shows all orders that have been placed
- **GET /api/order/totalled** - Returns a summary showing total orders per drink
- **POST /api/order** - Lets you place a new order

### 4. OpenAPI Documentation
- Ran the `build-openapi-docs.sh` script after completing the code
- This generated the Swagger docs that you can view at `/openapi/index.html`

## How to Run It

```bash
# First, download dependencies
go mod tidy

# Build the application
go build -o ordersystem

# Start the server
./ordersystem
```

The server runs on port 3000. You can access:
- Frontend dashboard: http://localhost:3000/
- API documentation: http://localhost:3000/openapi/index.html

## Testing

I tested all the endpoints using curl:

```bash
# Get the menu
curl http://localhost:3000/api/menu

# Get all orders
curl http://localhost:3000/api/order/all

# Get order totals
curl http://localhost:3000/api/order/totalled

# Place a new order
curl -X POST http://localhost:3000/api/order \
  -H "Content-Type: application/json" \
  -d '{"drink_id": 1, "amount": 2}'
```

Everything works as expected!

## Project Structure

```
solution/
├── docs/               # OpenAPI documentation (auto-generated)
├── frontend/           # HTML dashboard
├── model/              # Data models (Drink, Order)
├── repository/         # Database handler
├── rest/               # API endpoint handlers
├── scripts/            # Build scripts
├── main.go             # Application entry point
└── go.mod              # Go dependencies
```

## Notes

- I used the exact drink data from the presentation slides (Beer $2.00, Spritzer $1.40, Coffee $0.00)
- The CreatedAt timestamp gets set automatically when you post an order
- Error handling is included for invalid JSON in POST requests
- All the code compiles without warnings

## What I Learned

This exercise helped me understand:
- How to structure a REST API in Go
- Working with the Chi router framework
- Generating OpenAPI documentation from code comments
- In-memory data storage patterns
- Proper error handling in HTTP handlers

---

For more details, see SOLUTION.md

