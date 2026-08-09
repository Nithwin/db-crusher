# Chapter 2: Mastering PostgreSQL in Go

In Chapter 1, we built an API that returned a simple text string. Real applications need real data. In this chapter, we will learn how to connect our Go application to a PostgreSQL database.

## What is PostgreSQL?
PostgreSQL (or just Postgres) is a highly advanced, open-source Relational Database. "Relational" means it stores data in tables (like Excel spreadsheets) with strict rows and columns, and these tables can relate to each other (e.g., a User table relates to a Transactions table).

## Step 1: Setting up the Database Code
To keep our code clean, we won't put database logic in `main.go`. We will create a new file.

In your terminal, run:
```bash
mkdir -p internal/database
touch internal/database/postgres.go
```

Open `postgres.go`. We need to import Go's built-in `database/sql` library. We also need a "Driver". The Go standard library knows *how* to talk to databases generally, but it needs a specific translator for Postgres. We use `pgx`.

```go
package database

import (
    "database/sql" 
    
    // The underscore means we just want to load the driver in the background, 
    // we don't want to call its functions directly.
    _ "github.com/jackc/pgx/v5/stdlib" 
)
```

## Step 2: The Connection Pool
If a million users visit your site, you cannot open a million individual database connections. Your database will crash. Instead, we create a **Connection Pool**. 

A connection pool opens a set number of connections (e.g., 50) and keeps them open. When a user needs data, they borrow a connection, use it, and put it back in the pool for the next user.

Let's write the connection function:
```go
func ConnectDB() *sql.DB {
    // 1. Define the connection string (usually hidden in a .env file for security!)
    // Format: postgres://username:password@host:port/database_name
    connString := "postgres://postgres:password@localhost:5432/db_crusher"
    
    // 2. Open the pool using our "pgx" driver
    db, err := sql.Open("pgx", connString)
    if err != nil {
        panic("Failed to read config!")
    }
    
    // 3. PING the database. sql.Open doesn't actually connect, it just reads the config.
    // db.Ping() forces it to actually reach out over the network to PostgreSQL.
    if err := db.Ping(); err != nil {
        panic("Database is offline!")
    }
    
    return db
}
```

## Step 3: Fetching Data
Let's pretend we have a table called `transactions`. We want to write a Go function that runs a SQL query to count them and get the average amount.

First, define the Go `struct` (blueprint) to hold the result:
```go
type Analytics struct {
    Status   string
    Count    int
    TotalAvg float64
}
```

Now, write the function to execute the query:
```go
func GetAnalytics(db *sql.DB) []Analytics {
    // 1. Execute the query. We use db.Query because we expect multiple rows back.
    rows, err := db.Query("SELECT status, COUNT(*), AVG(amount) FROM transactions GROUP BY status")
    if err != nil {
        return nil
    }
    
    // 2. CRITICAL: Always defer closing the rows! If you forget this, 
    // your connection pool will run out of connections and your app will freeze.
    defer rows.Close()

    var results []Analytics

    // 3. Loop over the rows the database sent back
    for rows.Next() {
        var a Analytics
        
        // 4. "Scan" copies the data from the PostgreSQL row into our Go struct.
        // The order must match the SQL SELECT statement!
        rows.Scan(&a.Status, &a.Count, &a.TotalAvg)
        
        // 5. Add the row to our list
        results = append(results, a)
    }

    return results
}
```

## Running the Code
If you call `GetAnalytics(db)` from your `main.go` handler, Go will securely connect to Postgres, run the SQL, and return real data to your API users. 

But what happens when thousands of people do this at once? Let's find out in Chapter 3.
