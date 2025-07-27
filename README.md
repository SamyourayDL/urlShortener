# UrlShortener service

Simple service to store shorten urls and redirect, when needed. 

# TechStack 
- Go
  - slog  
  - chi
- SQlite

## 📘 API

Docs:
- [Swagger-file](https://samyouraydl.github.io/urlShortener/)

## Methods
- `POST /url/{alias}` — add url to storage with {alias}
- `GET /{alias}` — redirect on a website under {alias}
- `DELETE /url/{alias}` — delete url with given {alias}
