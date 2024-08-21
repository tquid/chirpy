package main

import (
  "net/http"
)

func main () {
  sm := http.NewServeMux()
  s := http.Server{
    Handler: sm,
    Addr: ":8080",
  }
  s.ListenAndServe()
}
