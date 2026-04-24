# envdiff

> CLI tool to diff and reconcile environment variable files across staging and production configs.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff && go build -o envdiff .
```

---

## Usage

```bash
# Compare two .env files
envdiff staging.env production.env

# Output missing keys and value differences
envdiff --missing-only staging.env production.env

# Reconcile: generate a merged output file
envdiff --reconcile staging.env production.env -o reconciled.env
```

**Example output:**

```
~ DB_HOST        staging=db.staging.internal  prod=db.prod.internal
+ SENTRY_DSN     only in staging
- NEW_RELIC_KEY  only in production
```

Keys prefixed with `~` indicate value mismatches, `+` indicates keys missing from production, and `-` indicates keys missing from staging.

---

## Flags

| Flag | Description |
|------|-------------|
| `--missing-only` | Show only keys absent from one file |
| `--reconcile` | Merge files, preferring production values |
| `-o, --output` | Write result to a file instead of stdout |
| `--ignore` | Comma-separated list of keys to skip |

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)