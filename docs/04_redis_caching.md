# Chapter 4: The Power of Redis Caching

In Chapter 3, our Postgres database choked because reading from a hard drive and recalculating math 2,000 times a second is incredibly slow. To fix this, we need a way to answer the user *without* talking to the database.

Enter **Redis**. 

## What is Redis?
Redis is an in-memory data structure store. Unlike PostgreSQL, which writes data to a slow hard drive to keep it safe permanently, Redis stores data in **RAM** (memory). RAM is temporary (if the power goes out, the data vanishes), but it is *lightning fast*. 

We use Redis as a **Cache**. A cache is like a shortcut. If we already calculated the analytics 1 second ago, why calculate them again? Just save the answer in Redis and give it to the next 1,999 users instantly!

## Step 1: Connecting to Redis
Let's create a new file `internal/cache/redis.go`. We use the `go-redis` library.

```go
package cache

import (
    "context"
    "github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
    // 1. Point the client to your local Redis server
    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379", 
    })

    // 2. Ping it to make sure it's alive
    err := client.Ping(context.Background()).Err()
    if err != nil {
        panic("Redis is offline!")
    }
    
    return client
}
```

## Step 2: The "Cache-Aside" Pattern
Now we update our API handler. We are going to implement a famous architectural pattern called "Cache-Aside". 

Here is how we write it:

```go
func GetAnalytics(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // --- PART 1: THE CACHE CHECK ---
    // We ask Redis: "Do you have the string stored under the key 'analytics'?"
    cachedData, err := cacheClient.Get(ctx, "analytics").Result()
    
    if err == nil {
        // CACHE HIT! Redis has it!
        // We just throw the raw string straight back to the user.
        w.Write([]byte(cachedData))
        return // We STOP the function here. Postgres is never touched!
    }

    // --- PART 2: THE FALLBACK (CACHE MISS) ---
    // If we reach this line, Redis didn't have the data. 
    // We MUST go the slow route and ask Postgres.
    analyticsData := database.GetAnalytics(dbPool)

    // We convert the Postgres data into a JSON string
    jsonData, _ := json.Marshal(analyticsData)


    // --- PART 3: SAVING FOR THE FUTURE ---
    // We tell Redis to save this JSON string under the key "analytics".
    // We tell it to self-destruct (expire) in 10 seconds.
    cacheClient.Set(ctx, "analytics", string(jsonData), 10*time.Second)

    // Finally, give the data to the user who asked for it.
    w.Write(jsonData)
}
```

## The New JMeter Results
We ran the exact same 2,000-user JMeter test. Here is what happened:

1.  **User 1** asked for the data. Redis was empty (Cache Miss). User 1 waited 88ms for Postgres to calculate it. The API saved the answer to Redis.
2.  **Users 2 through 2000** asked for the data. Redis had the answer (Cache Hit)! 

The API served 1,999 users directly out of RAM. The wait time for all of them was **0 milliseconds**.

Our throughput skyrocketed from 41 requests per second to **101.3 requests per second**. We successfully crushed the bottleneck!
