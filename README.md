# URL Shortener

A simple URL shortening service built with Go.

## Features

- Shorten long URLs
- Redirect to original URLs
- Simple web interface
- In-memory storage

## Usage

```bash
go run main.go
```

Then visit http://localhost:8080

## API

- `GET /` - Home page with form
- `POST /shorten` - Shorten a URL
- `GET /r/{code}` - Redirect to original URL

## Example

1. Go to http://localhost:8080
2. Enter a URL like "https://www.google.com"
3. Get a short URL like "http://localhost:8080/r/aBc123"
4. The short URL redirects to the original