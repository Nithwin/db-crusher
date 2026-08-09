# DB-Crusher

DB-Crusher is a high-performance backend project built to explore database optimization, load testing, and caching strategies. The primary goal of this project was to implement a robust REST API in Go, connect it to PostgreSQL, and then crush it with traffic using Apache JMeter to see where it breaks. 

Once we found the limits, we introduced **Redis Caching** to understand its real-world impact.

## The Impact of Redis Caching

To learn Redis properly, we didn't just read about it; we tested it under heavy load. We ran a simulated load test of 2000 concurrent requests against our analytics endpoint.

Here is the dramatic difference Redis made when we cached the database results:

| Metric | Baseline (PostgreSQL Only) | With Redis Caching | Improvement |
| :--- | :--- | :--- | :--- |
| **Average Response Time** | 1175 ms | 0 ms | ~100% reduction! |
| **Median Response Time** | 890 ms | 1 ms | ~99.9% reduction |
| **99% Line (P99)** | 4235 ms | 3 ms | ~99.9% reduction |
| **Max Response Time** | 4922 ms | 88 ms | ~98.2% reduction |
| **Throughput** | 41.0 req/sec | 101.3 req/sec | ~147% increase (2.47x) |

### Why This Matters

Before Redis, hitting the database for every single request caused a massive bottleneck. The database had to compute and fetch the analytics data every time, causing our 99th percentile users to wait over 4 seconds! 

By implementing a Redis cache:
1. We intercept the request before it hits the database.
2. If the data is in Redis, we serve it instantly from memory.
3. Our application's throughput more than doubled, and wait times dropped to practically zero.

This proves how essential in-memory caching is for scaling backend systems.

## Documentation

* [Step-by-Step Learning Guide](learning_journey.md): Read this for a simple, human-readable guide on how this project was built from scratch.
* [Load Test Details](load-test/comparison.md): For full JMeter report metrics.
