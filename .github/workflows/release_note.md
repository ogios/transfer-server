# transfer-server

Data transfer tool written in Go

## Feat

- Metadata for easier manage and smaller network consumption.
- storing text using a 16KB page-based storage scheme.
- separate delete and clear delete to reduce runtime resource utilization thus faster response.
- UDP server for message update subscription, making metadata sync more efficiently
- third-party proxy server for ip addresses communicate

## TODO

- p2p
