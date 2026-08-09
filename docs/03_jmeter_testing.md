# Chapter 3: Load Testing with JMeter

In Chapter 2, we built a connection to PostgreSQL. If you open your browser and hit refresh, it works perfectly. But writing code that works for one person is easy. Writing code that survives 2,000 people clicking at the exact same second is hard.

To see if our code survives, we use **Apache JMeter**. 

## What is Apache JMeter?
JMeter is an open-source software designed to load test functional behavior and measure performance. Instead of asking 2,000 of your friends to click a button at the same time, JMeter creates "Virtual Users" (Threads) to bombard your server with traffic automatically.

## Step 1: Creating a Test Plan
A Test Plan is a blueprint for your attack on the server. Here is how we set it up:

1.  Open JMeter GUI.
2.  Right-click "Test Plan" -> Add -> Threads (Users) -> **Thread Group**.

**What is a Thread Group?** 
A Thread Group defines your users. We configured ours like this:
*   **Number of Threads (users):** `2000`
*   **Ramp-up period (seconds):** `1` (This means all 2,000 users will join the attack within a single second).
*   **Loop Count:** `1` (Each user makes one request and stops).

## Step 2: Telling the Users What to Do
Right now, the 2,000 users are standing around doing nothing. We need to give them a target.

1.  Right-click the Thread Group -> Add -> Sampler -> **HTTP Request**.
2.  We configure it to point to our local Go server:
    *   **Protocol:** `http`
    *   **Server Name or IP:** `localhost`
    *   **Port Number:** `8080`
    *   **HTTP Method:** `GET`
    *   **Path:** `/analytics`

## Step 3: Recording the Carnage (Listeners)
When the attack happens, we need tools to record how fast the server responded, or if it crashed.

1.  Right-click the Thread Group -> Add -> Listener -> **Summary Report**.
2.  Right-click again -> Add -> Listener -> **Aggregate Report**.

## Step 4: The Baseline Run
We pressed the green "Start" button. 

Immediately, 2,000 requests slammed into our Go API. Our Go API efficiently accepted all 2,000 requests (Go is amazing at this). But then, our Go API sent 2,000 heavy `SELECT ... GROUP BY` SQL queries to PostgreSQL.

**PostgreSQL choked.** It had to calculate the sum and average of the transactions 2,000 times simultaneously. The database CPU hit 100%. Requests were placed in a queue, waiting their turn.

### The Results
*   **Throughput:** Only `41 requests per second` could be completed.
*   **Average Wait Time:** `1175 milliseconds` (Over 1 second).
*   **Maximum Wait Time:** `4922 milliseconds` (The unluckiest user waited almost 5 seconds for a simple web page to load).

If this were a real application, users would assume the app was broken and close it. In Chapter 4, we will learn how to fix this using Redis.
