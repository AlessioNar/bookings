package forms

import (
  "net/http"
  "net/url"
)


// It creates a custom Form struct and it embeds a url.Values object
type Form struct {
  url.Values
  Errors errors
}

// New initializes a Form struct
func New(data url.Values) *Form {
  return &Form{
    data,
    errors(map[string][]string{}),
  }
}

// Has checks if form fields is in post and it is not empty
func (f *Form) Has(field string, r *http.Request) bool {
  x := r.Form.Get(field)
  if x == "" {
    return false
  }

  return true
}
