avocado-chunk
=============

This is the chunkserver application for avocado. It handles submitting randomness and maintaing the encrypted chunks of media.

To launch:
```
go run server.go <CONTRACT LOCATION> <ACCOUNT HASH>
```

Example:
```
go run server.go 0x8f0483125fcb9aaaefa9209d8e9d7b9c8b9fb90f 0xf17f52151ebef6c7334fad080c5704d77216b732
```
