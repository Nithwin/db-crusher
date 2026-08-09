# Performance Comparison (Baseline vs Redis Caching)

## Key Metrics Comparison

| Metric | Baseline (Without Caching) | With Redis Caching | Improvement |
| :--- | :--- | :--- | :--- |
| **Average Response Time** | 1175 ms | 0 ms | ~100% reduction |
| **Median Response Time** | 890 ms | 1 ms | ~99.9% reduction |
| **99% Line (P99)** | 4235 ms | 3 ms | ~99.9% reduction |
| **Max Response Time** | 4922 ms | 88 ms | ~98.2% reduction |
| **Throughput** | 41.0/sec | 101.3/sec | ~147% increase (2.47x) |
| **Error Rate** | 0.00% | 0.00% | No change |

## Analysis

Implementing Redis caching resulted in a massive performance improvement across all metrics. 
*   **Response Times:** The average and median response times dropped to near-zero levels (from ~1 second down to 0-1 ms). This demonstrates the immense benefit of serving repeated requests directly from in-memory cache instead of querying the database.
*   **Throughput:** The system's throughput more than doubled, jumping from 41.0 requests per second to 101.3 requests per second. The application can now handle significantly more concurrent traffic.
*   **Consistency:** The 99th percentile response time went from over 4 seconds down to just 3 milliseconds, highlighting a much more stable and predictable user experience, even for the slowest requests.
