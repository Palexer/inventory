# Inventory

Inventory is a selfhosted web app written in Go to keep track of items in stock.
It's meant for use in a small scope and might not be well suited for bigger databases.

Inventory works by using a CSV file for storing the data and an HTML table as a web interface.
This has the advantage of being able to view the table from whatever device you want.
You can sort the table by every column, so for example by count, date or alphabetically.

## Installation

Prebuilt binaries can be found under the releases section for Linux/amd64 and
Linux/arm64. 

If you want to host the application on another operating system you can simply clone this repository
and run in the src/ folder to get your binary.

```go build -o inventory```

## Usage

To start the server, just execute the binary. The following flags are available:

|Flag|Functionality|
|---|---|
|-port|specify a port to run the server on|
|-key|specify a custom key to delete the table (default: Inventory)|
|-heading|specify a custom table heading|
|-h/-help|get help information|

To get the currently used version type 

```
inventory version
```

## License

BSD-3-Clause License (see LICENSE file for details)

