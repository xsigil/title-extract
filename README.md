# title-extract

Ultra-fast, zero-allocation stream-oriented HTML `<title>` extractor written in Go.

## Features

- **Zero-Allocation Streaming**: Tokenizer-based single-pass parsing via `golang.org/x/net/html` without building a DOM tree.
- **Early Exit Optimization**: Scanning aborts immediately when `<title>` is resolved or upon reaching `</head>` / `<body>`.
- **POSIX Standard Compliant**: Reads seamlessly from `stdin` or multiple positional file arguments.
- **Sanitized Output**: Decodes HTML entities and folds redundant whitespace.

## Build & Install

```sh
# Build binary into build/title-extract
make

# Install binary and man page into ~/.local
make install

```

## Usage

```sh
# Via pipe (streaming)
curl -sL [https://example.com](https://example.com) | title-extract

# Via arguments
title-extract page1.html page2.html

# Help
title-extract --help

```

## Manual

```sh
man title-extract

```

## License

MIT
