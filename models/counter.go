package models

// the Counter model provides a generic dataset to increment and decrement
// counters across the entire application with varying levels of granularity
type Counter struct {
  ID      string
  Count   int
}
