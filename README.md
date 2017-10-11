# Simple BoltDB Client

Simple client with utility methods for
- opening and closing DB - Open, Close
- creation of initial entities (for example buckets) - InitEntity
- read transaction with callback - ReadTransaction
- write transaction with callback - WriteTransaction
- delete entry from bucket - DeleteWithTransaction
- find nested bucked according to path - FindBucket