package main

import (
    "encoding/json"
    "fmt"
    "os"
    "sync"
)

type Storage struct {
    urls map[string]string
    mu   sync.RWMutex
}

func NewStorage() *Storage {
    return &Storage{
        urls: make(map[string]string),
    }
}

func (s *Storage) Set(shortCode, longURL string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.urls[shortCode] = longURL
}

func (s *Storage) Get(shortCode string) (string, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    url, exists := s.urls[shortCode]
    return url, exists
}

func (s *Storage) SaveToFile(filename string) error {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    encoder := json.NewEncoder(file)
    return encoder.Encode(s.urls)
}

func (s *Storage) LoadFromFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    decoder := json.NewDecoder(file)
    
    s.mu.Lock()
    defer s.mu.Unlock()
    
    return decoder.Decode(&s.urls)
}

func (s *Storage) Stats() map[string]interface{} {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    return map[string]interface{}{
        "total_urls": len(s.urls),
    }
}