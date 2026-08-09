# Lesson 4: Step-by-Step Redis Caching Implementation

To prevent our PostgreSQL database from crashing under heavy load, we introduce **Redis**. Redis is an in-memory data store. Reading data from RAM is lightning fast compared to recalculating SQL queries on a hard drive.

## Step 1: Connecting to Redis

First, we need to import the official `go-redis` library and create a connection client. Just like PostgreSQL, we only want to do this once when the server starts.

```go
package cache

import (
    "context"
    "github.com/redis/go-redis/v9"
)

func ConnectRedis(ctx context.Context) *redis.Client {
    // 1. Create a new client pointing to the default local Redis port
    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379", 
    })

    // 2. Ping it to ensure it is actually running
    err := client.Ping(ctx).Err()
    if err != nil {
        panic(err) // Stop the app if Redis isn't running
    }
    
    // 3. Return the ready-to-use client
    return client
}
```

## Step 2: The Cache-Aside Logic (The Handler)

Now we need to update our `GetAnalytics` API handler. Instead of going straight to the database, we are going to implement the **Cache-Aside Pattern**.

Let's build the logic step-by-step.

### Part A: Check the Cache First

When a user asks for analytics, we check Redis first:

```go
func (h *UserHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // 1. Ask Redis: "Do you have data saved under the key 'analytics'?"
    cachedData, err := h.Cache.Get(ctx, "analytics").Result()
    
    if err == nil {
        // 2. CACHE HIT! Redis found the data.
        // We write the cached JSON string directly back to the user.
        w.Write([]byte(cachedData))
        
        // 3. RETURN! We stop the function right here. The database is NEVER touched.
        return 
    }
    
    // If we get down here, it means we had a CACHE MISS. Redis did not have the data.
}
```

### Part B: The Database Fallback

If Redis doesn't have the data (a Cache Miss), we must fetch it from Postgres, convert it to JSON, and send it to the user.

```go
func (h *UserHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
    // ... (Cache check from Part A happened here) ...

    // 4. Ask Postgres for the data
    analytics, _ := h.DB.GetAnalytics()

    // 5. Convert the Go struct data into a JSON string
    jsonData, _ := json.Marshal(analytics)

    // 6. Send the JSON data to the user
    w.Write(jsonData)
}
```

### Part C: Saving the Data in Redis

Right now, every request is a Cache Miss. We need to save the data in Redis so the *next* request is a Cache Hit!

```go
func (h *UserHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
    // ... (Part A & B) ...
    analytics, _ := h.DB.GetAnalytics()
    jsonData, _ := json.Marshal(analytics)

    // 7. Save the data in Redis under the key "analytics"
    // We set an expiration time of 10 seconds. 
    // This ensures data stays fresh; after 10 seconds, the cache deletes itself!
    h.Cache.Set(ctx, "analytics", string(jsonData), 10*time.Second)

    w.Write(jsonData)
}
```

## The Final Flow and Results

By combining these steps, this is what happens when 2000 users hit the API at the exact same second:

```mermaid
flowchart TD
    User1[User 1] --> API[Go API /analytics]
    User2[User 2 to 2000] --> API
    
    API --> CheckCache{Check Redis}
    
    CheckCache -- "Cache Miss (User 1)" --> DB[PostgreSQL]
    DB --> Calc[Calculate Data (Takes 1 second)]
    Calc --> Save[Save to Redis for 10s]
    Save --> Return1[Return to User 1]
    
    CheckCache -- "Cache Hit (Users 2-2000)" --> Inst[Return Instantly from RAM!]
    Inst --> Return2[Return to Users 2-2000]
```

### The New Load Test Results
When we re-ran our JMeter test, the results were staggering:
*   **Average Wait Time:** `0 ms` (Down from 1175ms).
*   **Max Wait Time:** `88 ms` (This was just User 1. Users 2-2000 got 0ms wait times).
*   **Throughput:** `101.3 req/sec` (More than double our original capacity).

By intercepting requests before they hit the database, we solved our performance bottleneck!
