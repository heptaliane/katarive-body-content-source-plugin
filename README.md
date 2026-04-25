# katarive-body-content-source-plugin

A [katarive](https://github.com/heptaliane/katarive-server) source plugin that extracts the main body content from web pages.

## Features

- Extracts the page title and plain text content.
- Automatically filters out non-content elements such as:
  - `<script>` and `<noscript>`
  - `<style>`
  - `<iframe>`
  - `<header>` and `<footer>`
- Supports any `http://` or `https://` URLs.

## Installation

To build the plugin from source, you need [Go](https://go.dev/) installed.

```bash
go build -o katarive-body-content-source-plugin .
```

## Usage

This is a `katarive` plugin. Once built, you can configure `katarive` to use this binary as a source plugin.

For more information on how to use plugins with `katarive`, please refer to the [katarive-server documentation](https://github.com/heptaliane/katarive-server).

## Development

### Running Tests

```bash
go test -v ./...
```

## License

[MIT](LICENSE)
