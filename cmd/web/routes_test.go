package main

import (
  "github.com/go-chi/chi"
  "github.com/tsawler/bookings-app/internal/config"
  "testing"
  "fmt"

)

func TestRoutes(t *testing.T) {
  var app config.AppConfig

  mux := routes(&app)

  switch v:= mux.(type) {
  case *chi.Mux:
    // do nothing, test passed
  default:
    t.Error(fmt.Sprintf("type is not *chi.Mux, type is %T", v))
  }

}
