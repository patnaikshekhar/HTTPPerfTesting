# HTTP Performance Testing Tool

A simple http performance testing tool written in go. It will spin up a number
of goroutines and make http get calls to the url passed as an argument.

## Usage
```sh
./perftest --concurrency <Number of Concurrent Requests> --number <Total Number of Requests> <url>
```

## Example
```sh
./perftest --concurrency 100 --number 10000 http://localhost:8080
```