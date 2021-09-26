# Go Kitchen!

## Getting started
NOTE: All commands are considered to be run from project root

### Clone the repo
```shell
$ git clone https://github.com/FoxFurry/GoKitchen.git
```

### Install Go 1.17 (or at least 1.15)
[Go Install Guide](https://golang.org/doc/install)

### Install the dependencies
```shell
$ go mod download
```
## Start the kitchen

### Simple start with default config path
```shell
go run main.go
```

### Build the project
```shell
go build -o <binary_name>
```

## Coverage

| Package | Coverage |
| ----------- | ----------- |
| No tests available right now | :( |


## Config file

### Content

#### kitchen_host
Address of host:port which will run the kitchen.<br>

**Example:**<br>
```json
"localhost:8080"
```

#### dining_host
Address of dining hall server.

**Example:**<br>
```json
"localhost:8081"
```

#### log_level
Specifies level of messages which will be logged.<br>
There are 4 different log levels available:

- 0 - panic
- 1 - show critical errors
- 2 - show critical and non-critical errors
- 3 - show everything

**Example:**<br>
```json
2
```