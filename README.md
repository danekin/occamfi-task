# OccamFi Task

Repository contains my implementation of OccamFi Interview Task

### Packages list:
- `aggregator` - contains simple object that calculates single ticker price from given array.
- `combiner` - provides main logic for combining ticker prices from multiple streams.
- `receiver` - takes responsibility for subscribing on price streams.
- `sender` - provides simple object for storing ticker prices in local file in csv-like format.

`cmd/main.go` - contains example of how to run an application with dummy mock subscribers.