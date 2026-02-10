# kubectl-plugin-arcane

A kubectl plugin written in Go.

## Installation

```bash
go build -o kubectl-arcane main.go
# Move to a directory in your PATH
sudo mv kubectl-arcane /usr/local/bin/
```

## Usage

```bash
kubectl arcane
```

## Development

### Prerequisites
- Go 1.21 or higher

### Build
```bash
go build -o kubectl-arcane main.go
```

### Run
```bash
./kubectl-arcane
```
