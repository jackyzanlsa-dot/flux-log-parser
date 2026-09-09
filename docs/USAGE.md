# Usage Guide

This document provides detailed usage information for `flux-log-parser`.

## Flag Reference

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-file` | string | stdin | Path to log file. If not specified, reads from standard input. |
| `-level` | string | `info` | Minimum log level to output. Valid values: `debug`, `info`, `warn`, `error`, `fatal`. |
| `-format` | string | `json` | Output format. Valid values: `json`, `logfmt`. |

## Level Filtering

The tool filters logs based on severity levels with the following hierarchy (lowest to highest):

1. `debug` (0) - Shows all logs
2. `info` (1) - Shows info, warn, error, fatal
3. `warn` (2) - Shows warn, error, fatal
4. `error` (3) - Shows error, fatal
5. `fatal` (4) - Shows only fatal

Level detection is automatic:
- Lines containing "ERROR" → error level
- Lines containing "WARN" → warn level
- Lines containing "DEBUG" → debug level
- All other lines → info level

## Examples

### Basic Usage

Read from a file and output JSON:
```bash
flux-log-parser -file app.log -format json
```

Read from stdin (piped input):
```bash
cat app.log | flux-log-parser -format json
```

### Filtering by Level

Show only errors and fatal logs:
```bash
flux-log-parser -file app.log -level error -format json
```

Show warnings and above:
```bash
flux-log-parser -file app.log -level warn -format json
```

Show all logs including debug:
```bash
flux-log-parser -file app.log -level debug -format json
```

### Output Formats

#### JSON Output

```bash
flux-log-parser -file app.log -level error -format json
```

Output:
```json
{"timestamp":"2025-01-15T10:30:00Z","level":"error","message":"Database connection failed","fields":{"host":"db.example.com","port":5432}}
```

#### Logfmt Output

```bash
flux-log-parser -file app.log -level error -format logfmt
```

Output:
```
time="2025-01-15T10:30:00Z" level="error" msg="Database connection failed" host="db.example.com" port="5432"
```

### CI/CD Pipeline Examples

#### GitHub Actions

```yaml
- name: Parse logs for errors
  run: |
    cat build.log | flux-log-parser -level error -format json > errors.json
    if [ -s errors.json ]; then
      echo "Errors found in build log"
      exit 1
    fi
```

#### GitLab CI

```yaml
check_logs:
  script:
    - flux-log-parser -file deploy.log -level error -format json | tee errors.json
    - test ! -s errors.json || exit 1
```

#### CircleCI

```yaml
- run:
    name: Check for critical errors
    command: |
      flux-log-parser -file application.log -level fatal -format json > fatal_errors.json
      # Exit with error if any fatal errors found
      [ $(wc -l < fatal_errors.json) -eq 0 ]
```

## Exit Code Behavior

The tool follows standard Unix exit code conventions for CI/CD integration:

| Exit Code | Meaning |
|-----------|---------|
| `0` | Success - logs processed without errors |
| `1` | Error - invalid arguments, file not found, or parse errors |

### CI/CD Best Practices

1. **Fail on Errors**: Use the exit code to fail pipeline steps when parsing fails:
   ```bash
   flux-log-parser -file app.log -level error || exit 1
   ```

2. **Capture Errors for Review**: Redirect filtered errors to a file for later review:
   ```bash
   flux-log-parser -file build.log -level error -format json > errors.json
   ```

3. **Conditional Failure**: Only fail if errors exceed a threshold:
   ```bash
   flux-log-parser -file app.log -level error | wc -l | grep -q '^0$' || exit 1
   ```

4. **Combine with Other Tools**: Pipe output to tools like `jq` for further processing:
   ```bash
   flux-log-parser -file app.log -format json | jq 'select(.level == "error")'
   ```

## Input Format Support

### JSON Logs

The parser handles standard JSON log lines:
```json
{"timestamp":"2025-01-15T10:30:00Z","level":"error","message":"Connection timeout","host":"api.example.com"}
```

### Plain Text Logs

Plain text logs are parsed with level detection based on keywords:
```
2025-01-15 10:30:00 ERROR Database connection failed
2025-01-15 10:30:01 WARN High memory usage detected
2025-01-15 10:30:02 DEBUG Processing request id=12345
2025-01-15 10:30:03 Request completed successfully
```
