# Lesson 3: Step-by-Step Load Testing with JMeter

Writing code that works for one user is easy. Writing code that works for 2,000 users simultaneously is hard. To test our API, we used Apache JMeter to simulate heavy traffic. 

Here is exactly how we set up the test, step-by-step.

## Step 1: Creating a Thread Group

In JMeter, a "Thread" represents a virtual user. A "Thread Group" is a collection of these users. 

1.  **Open JMeter**, right-click the Test Plan -> Add -> Threads (Users) -> **Thread Group**.
2.  We configured the Thread Group with the following settings:
    *   **Number of Threads (users):** `2000` (This means 2000 users will hit the API).
    *   **Ramp-up period (seconds):** `1` (This means all 2000 users will start within 1 second. This creates an immediate, massive spike in traffic).
    *   **Loop Count:** `1` (Each user makes exactly one request).

## Step 2: Adding an HTTP Request

Now that we have our virtual users, we need to tell them exactly what to do.

1.  Right-click the Thread Group -> Add -> Sampler -> **HTTP Request**.
2.  We configured the request to hit our Go API's heavy database endpoint:
    *   **Protocol:** `http`
    *   **Server Name or IP:** `localhost`
    *   **Port Number:** `8080`
    *   **HTTP Method:** `GET`
    *   **Path:** `/analytics`

## Step 3: Adding Listeners (Reporting)

We need to capture the results to see if the API survives. 

1.  Right-click the Thread Group -> Add -> Listener -> **Summary Report**. (This gives us averages, minimums, and maximums).
2.  Right-click again -> Add -> Listener -> **Aggregate Report**. (This gives us percentiles, like the 99th percentile, which tells us how bad the experience is for the slowest users).

## Step 4: Running the Test (The Baseline)

We clicked the green "Start" button in JMeter. All 2000 requests hit the Go server, which immediately forwarded all 2000 requests to the PostgreSQL database. 

Because the `/analytics` endpoint requires the database to calculate `SUM` and `COUNT` for every request, the database CPU maxed out.

### The Results:
*   **Average Wait Time:** `1175 ms` (Over a full second).
*   **Maximum Wait Time:** `4922 ms` (Almost 5 seconds for the unluckiest users!).
*   **Throughput:** `41.0/sec` (The server could only handle 41 requests per second before choking).

### The Bottleneck Diagram

```mermaid
sequenceDiagram
    participant 2000 Virtual Users
    participant Go Web Server
    participant PostgreSQL Database
    
    2000 Virtual Users->>Go Web Server: Simultaneous GET /analytics
    Note over Go Web Server: Go easily accepts all 2000 connections
    Go Web Server->>PostgreSQL Database: Send 2000 heavy SQL Queries
    Note over PostgreSQL Database: WARNING: CPU 100%. Processing queues up.
    PostgreSQL Database-->>Go Web Server: Slowly returns data one by one
    Go Web Server-->>2000 Virtual Users: Users wait 1 to 5 seconds
```

This test proved that hitting a relational database for heavy, repetitive analytics data is a terrible idea for scaling. In the next lesson, we fix this by intercepting requests with a Redis Cache.
