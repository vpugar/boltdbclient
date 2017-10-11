# Simple BoltDB Client

<p align="left">
    <a href="https://travis-ci.org/vpugar/boltdbclient"><img src="https://travis-ci.org/vpugar/boltdbclient.svg?branch=master" alt="Build Status"></a>
    <a href="https://coveralls.io/github/vpugar/boltdbclient?branch=master"><img src="https://coveralls.io/repos/vpugar/boltdbclient/badge.svg?branch=master&service=github" alt="Coverage Status"></a>
    <a href="https://goreportcard.com/report/github.com/vpugar/boltdbclient"><img src="https://goreportcard.com/badge/github.com/vpugar/boltdbclient" alt="Go Report Card"></a>
</p>


Simple client with utility methods for
- opening and closing DB - Open, Close
- creation of initial entities (for example buckets) - InitEntity
- read transaction with callback - ReadTransaction
- write transaction with callback - WriteTransaction
- delete entry from bucket - DeleteWithTransaction
- find nested bucked according to path - FindBucket