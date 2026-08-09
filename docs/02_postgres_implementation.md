# Lesson 2: PostgreSQL Implementation in Go

An API needs data. We chose PostgreSQL for our database and used the `pgx` driver because of its incredible performance in Go.

## Connecting to the Database

Connecting to a database shouldn't happen inside every request. We need to create a **connection pool**. A pool keeps several connections open, so when a user requests data, a connection is already ready to go.

Here's how we did it in `internal/database/postgres.go`:

```go
func NewDB() (*Database, error) {
    // 1. Build the connection string from environment variables
    connectionString := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s",
        user, password, host, port, name,
    )
    
    // 2. Open a connection pool using the pgx driver
    db, err := sql.Open("pgx", connectionString)
    
    // 3. Ping the database to ensure it's actually reachable
    if err := db.Ping(); err != nil {
        db.Close()
        return nil, err
    }
    
    // Return our custom Database struct
    return &Database{DB: db}, nil
}
```

## Executing a Query (Fetching Data)

Let's see how we query the database to get analytics data. This is a heavy query that counts and averages transactions.

```go
func (d *Database) GetAnalytics() ([]AnalyticsResult, error) {
    // 1. Execute the SQL query
    rows, err := d.DB.Query(`
        SELECT status, COUNT(*), SUM(amount), AVG(amount)
        FROM transactions
        GROUP BY status
    `)
    defer rows.Close() // ALWAYS close rows to prevent memory leaks

    var results []AnalyticsResult

    // 2. Loop through every row returned by the database
    for rows.Next() {
        var res AnalyticsResult

        // 3. Scan (copy) the data from the database row into our Go struct
        rows.Scan(&res.Status, &res.Count, &res.TotalAmount, &res.AverageAmount)
        
        // 4. Add it to our results slice
        results = append(results, res)
    }

    return results, nil
}
```

### Explaining the Flow

```mermaid
flowchart TD
    A[API Request: GET /analytics] --> B[Handler Calls GetAnalytics()]
    B --> C[DB.Query executes SQL]
    C --> D[(PostgreSQL Database)]
    D --> E[Returns Rows]
    E --> F[Go loops through Rows.Next]
    F --> G[JSON Encoded & Sent to User]
```

*   **`d.DB.Query`**: Used when we expect multiple rows back.
*   **`rows.Scan`**: Maps the columns from the SQL database directly into the fields of our Go variables.

In the next lesson, we will push this database to its absolute limits using JMeter.
