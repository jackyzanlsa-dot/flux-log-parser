# flux-log-parser

High-performance zero-dependency CLI for parsing and streaming structured logs in real-time.

## Features

- Streaming parser for JSON and logfmt without loading full files into memory
- Filter engine by level/timestamp/key-value
- CI/CD friendly exit codes
- Zero dependencies single binary

## Installation

```bash
go install github.com/flux-cli-tools/flux-log-parser@latest
```

## Usage

Parse a log file with error level filter in JSON format:

```bash
flux-log-parser -file app.log -level error -format json
```

Stream logs from stdin with warn level filter:

```bash
cat app.log | flux-log-parser -level warn -format json
```

## Development

```bash
git clone https://github.com/flux-cli-tools/flux-log-parser.git
cd flux-log-parser
go build
```

## Docs

For detailed usage information, see [docs/USAGE.md](docs/USAGE.md).

## Contributing

Issues and pull requests are welcome.

## License

MIT