<p align="center">
  <img src="https://gitlab.com/jjlabs-dev/frumpy/-/raw/main/branding/frumpy.png" alt="Frumpy Logo" width="300" height="300">
</p>

# FRUMPy — [F]ind [R]epos, [U]ploads, [M]irrors, [P]ackages, [P]ayloads[y]

Fetches package and repository info from a Nexus Repository and prints it as a simple CLI table.

## Features
- List packages stored in a Nexus repository
- Minimal, dependency-light CLI designed for quick inspection.

## Requirements
- Go 1.25+
- Make
- Network access to your Nexus repository
- Valid credentials if your Nexus instance requires authentication

## Installation
Build from source:
```bash
git clone <repo-url>
cd frumpy
make [mac-|imac-]install
```

## Usage
Basic command to list packages from a repository named `quay`:
```bash
fr quay
```

Example output columns:
- Name
- Format
- Version
- DownloadUrl
- Uploader
- SHA1
- LastModified
- LastDownloaded
- FileSize

## Configuration
Set Nexus base URL and credentials via environment variables or config file (check code for exact names). Common env vars:
- `FRUMPY_URL`
- `FRUMPY_USERNAME`
- `FRUMPY_PASSWORD` or `FRUMPY_TOKEN`

## Examples
List all entries from repository `quay`:
```bash
fr quay
```

Supports regex filtering
```bash
fr quay prometheus
```


## License
Copyright (c) 2023 jjlabs.dev

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

