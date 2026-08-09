# DB-Crusher

DB-Crusher is a high-performance backend project built to explore database optimization, load testing, and caching strategies. The primary goal of this project was to implement a robust REST API in Go, connect it to PostgreSQL, and then crush it with traffic using Apache JMeter to see where it breaks. 

Once we found the limits, we introduced **Redis Caching** to understand its real-world impact.

## How to Run This Project

To run this project locally, you will need Go, PostgreSQL, and Redis installed on your machine.

### 1. Prerequisites
*   [Go (1.22+)](https://go.dev/dl/)
*   [PostgreSQL](https://www.postgresql.org/download/)
*   [Redis](https://redis.io/download)

### 2. Setup Environment Variables
Create a `.env` file in the root directory and add your database credentials:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=db_crusher
```

### 3. Run the Server
Open your terminal, navigate to the project directory, and run:
```bash
go mod tidy
go run main.go
```
You should see:
```text
Redis connected successfully
Your server is running in http://localhost:8080
```
You can now test the API by visiting `http://localhost:8080/health` in your browser!

---

## The Course: Backend Mastery Journey

I have written a comprehensive, book-style course explaining exactly how this project was built. It teaches the fundamental concepts of backend engineering from scratch.

Read the chapters here:
1. [Chapter 1: Teaching Go REST APIs](docs/01_go_rest_api.md)
2. [Chapter 2: Mastering PostgreSQL in Go](docs/02_postgres_implementation.md)
3. [Chapter 3: Load Testing with JMeter](docs/03_jmeter_testing.md)
4. [Chapter 4: The Power of Redis Caching](docs/04_redis_caching.md)
5. [Chapter 5: Full Architecture & Summary](docs/05_summary_architecture.md)

## The Impact of Redis Caching

Here is the dramatic difference Redis made when we cached the database results for 2000 concurrent requests:

| Metric | Baseline (PostgreSQL Only) | With Redis Caching | Improvement |
| :--- | :--- | :--- | :--- |
| **Average Response Time** | 1175 ms | 0 ms | ~100% reduction! |
| **Median Response Time** | 890 ms | 1 ms | ~99.9% reduction |
| **99% Line (P99)** | 4235 ms | 3 ms | ~99.9% reduction |
| **Max Response Time** | 4922 ms | 88 ms | ~98.2% reduction |
| **Throughput** | 41.0 req/sec | 101.3 req/sec | ~147% increase (2.47x) |

Check out the [Load Test Comparison Details](load-test/comparison.md) for full JMeter report metrics.
