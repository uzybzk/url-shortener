package main

import (
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "net/url"
    "time"
)

var storage = NewStorage()

func main() {
    rand.Seed(time.Now().UnixNano())
    
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/shorten", shortenHandler)
    http.HandleFunc("/r/", redirectHandler)
    
    fmt.Println("URL Shortener starting on :8080")
    fmt.Println("Visit http://localhost:8080 to use the service")
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    
    html := `
    <!DOCTYPE html>
    <html>
    <head><title>URL Shortener</title></head>
    <body>
        <h1>URL Shortener</h1>
        <form action="/shorten" method="POST">
            <input type="url" name="url" placeholder="Enter URL to shorten" required>
            <button type="submit">Shorten</button>
        </form>
    </body>
    </html>
    `
    
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprint(w, html)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    longURL := r.FormValue("url")
    if longURL == "" {
        http.Error(w, "URL is required", http.StatusBadRequest)
        return
    }
    
    // Validate URL
    _, err := url.ParseRequestURI(longURL)
    if err != nil {
        http.Error(w, "Invalid URL", http.StatusBadRequest)
        return
    }
    
    shortCode := generateShortCode()
    storage.Set(shortCode, longURL)
    
    shortURL := fmt.Sprintf("http://localhost:8080/r/%s", shortCode)
    
    html := fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head><title>URL Shortened</title></head>
    <body>
        <h1>URL Shortened Successfully</h1>
        <p>Original URL: %s</p>
        <p>Short URL: <a href="%s">%s</a></p>
        <a href="/">Shorten another URL</a>
    </body>
    </html>
    `, longURL, shortURL, shortURL)
    
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprint(w, html)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
    shortCode := r.URL.Path[len("/r/"):]
    
    longURL, exists := storage.Get(shortCode)
    if !exists {
        http.NotFound(w, r)
        return
    }
    
    http.Redirect(w, r, longURL, http.StatusFound)
}

func generateShortCode() string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    const length = 6
    
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
}