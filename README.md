# DB-Crusher

DB-Crusher is a high-performance backend project built to explore database optimization, load testing, and caching strategies. The primary goal of this project was to implement a robust REST API in Go, connect it to PostgreSQL, and then crush it with traffic using Apache JMeter to see where it breaks. 

Once we found the limits, we introduced **Redis Caching** to understand its real-world impact.

## The Course: Building DB-Crusher

I have written a step-by-step course explaining exactly how this project was built, complete with code snippets, architecture diagrams, and plain-english explanations of how the functions work.

Read the course here:
1. [Lesson 1: Creating a Go REST API](docs/01_go_rest_api.md)
2. [Lesson 2: PostgreSQL Implementation](docs/02_postgres_implementation.md)
3. [Lesson 3: Load Testing with JMeter](docs/03_jmeter_testing.md)
4. [Lesson 4: Supercharging Performance with Redis](docs/04_redis_caching.md)

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
