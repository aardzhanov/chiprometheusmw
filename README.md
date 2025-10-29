# Prometheus middleware library for Chi http framework

```go
import "github.com/aardzhanov/prometheusmw"
...
r := chi.NewRouter()
metrics := prometheusmw.NewMiddleware("backend")
r.Use(metrics.Handler)
```

With custom bucket for histogram (0,1, 0,5, 1 an5 5 seconds)
```go
import "github.com/aardzhanov/prometheusmw"
...
r := chi.NewRouter()
metrics := prometheusmw.NewMiddleware("backend", 0.1, 0.5, 1, 5)
r.Use(metrics.Handler)
```
