# envdiff

Compare `.env` files across environments and highlight missing or mismatched keys.

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
envdiff [flags] <file1> <file2>
```

**Example:**

```bash
envdiff .env.development .env.production
```

**Sample output:**

```
MISSING in .env.production:
  - DATABASE_URL
  - REDIS_HOST

MISMATCHED values:
  - LOG_LEVEL: "debug" vs "info"
  - PORT: "3000" vs "8080"
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--keys-only` | Only compare key names, ignore values |
| `--quiet` | Exit with non-zero code if differences found (useful in CI) |
| `--format json` | Output results as JSON |

---

## Why envdiff?

Keeping environment configs in sync across staging, production, and local setups is error-prone. `envdiff` gives you a fast, scriptable way to catch drift before it causes issues in production.

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)