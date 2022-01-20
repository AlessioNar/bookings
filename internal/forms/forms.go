package forms

import (
  "net/http"
  "net/url"
  "strings"
  "fmt"
)


// It creates a custom Form struct and it embeds a url.Values object
type Form struct {
  url.Values
  Errors errors
}
// Valid returns true if there are no errors, otherwise it is false
func (f *Form) Valid() bool {
  return len(f.Errors) == 0
}

// New initializes a Form struct
func New(data url.Values) *Form {
  return &Form{
    data,
    errors(map[string][]string{}),
  }
}

// Checks for required fields
func (f *Form) Required(fields ...string) {
  for _, field := range fields {
    value := f.Get(field)
    if strings.TrimSpace(value) == "" {
      f.Errors.Add(field, "This field cannot be blank")
    }
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

// MinLength checks for string minimum length
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
  x := r.Form.Get(field)
  if len(x) < length {
    f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
    return false
  }
  return true
}
