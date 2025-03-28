# GFinder: Advanced GitHub Code Search Utility

## Overview

GFinder is a powerful command-line tool for searching and extracting code snippets from GitHub repositories using flexible search parameters. It allows devebug bounty hunters to perform advanced code searches with multiple filtering options.

## Features

- Search GitHub code repositories using custom queries
- Filter results with regular expressions
- Extract specific content types:
  - Raw code snippets
  - Unique URLs
  - Domain names
- Customizable search pagination
- Silent mode for clean output
- Optional GitHub API authentication

## Installation

### Prerequisites

- Go (version 1.16 or higher)
- GitHub API access token (optional, but recommended)

### Install via Go

```bash
go install github.com/gilsgil/gfinder@latest
```

### Manual Installation

```bash
git clone https://github.com/gilsgil/gfinder.git
cd gfinder
go build
```

## Usage

### Basic Search

```bash
# Basic search for code containing a specific term
gfinder -q "example search term" -r "regex_pattern"
```

### Advanced Search Modes

1. Default Mode (Code Snippets):
```bash
gfinder -q "mercadolivre" -r "sensitive_pattern"
```

2. URL Extraction Mode:
```bash
gfinder -q "mercadolivre" -m urls -r "https://.*example\.com"
```

3. Domain Extraction Mode:
```bash
gfinder -q "mercadolivre" -m domains -r "example\.com"
```

### Parameters

- `-q`: Search query for GitHub API
- `-r`: Regular expression for filtering results
- `-m`: Extraction mode (`urls` or `domains`)
- `-d`: Delay between requests (default: 2 seconds)
- `-s`: Silent mode (only unique results)

### Authentication

Set GitHub API token to increase rate limits:
```bash
export GITHUB_KEY=your_github_token
```

## Examples

```bash
# Find code snippets related to payment processing
gfinder -q "payment gateway" -r "credit_card"

# Extract unique URLs from repositories
gfinder -q "api endpoint" -m urls -r "https://secure\..*"

# Find domains in code repositories
gfinder -q "cloud service" -m domains -r "aws\.com"
```

## Limitations

- Maximum of 1000 results (10 pages of 100 items)
- Requires internet connection
- GitHub API rate limits apply

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Contact

Gil's GitHub - [@gilsgil](https://github.com/gilsgil)

Project Link: [https://github.com/gilsgil/gfinder](https://github.com/gilsgil/gfinder)