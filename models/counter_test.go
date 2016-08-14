package models

import (
  "testing"
)

func TestCounterSave(t *testing.T) {
  // clear the counters table when we're finished
  defer truncate("counters")
  
  counter := Counter{
    ID:         "123",
    Count:      1,
  }
  
  err := counter.Save()
  if err != nil {
    t.Error(err)
  }
  
  // fetch the counter
  counter2 := Counter{}
  
  if err = get("counters", counter.ID, &counter2); err != nil {
    t.Errorf("failed to fetch counter: %s", err.Error())
  }
  
  if counter2.ID != "123" {
    t.Errorf("counter doesn't match")
  }
}

func TestCounterDelete(t *testing.T) {
  // clear the counters table when we're finished
  defer truncate("counters")
  
  counter := Counter{
    ID:      "123",
    Count:   1,
  }
  
  if err := counter.Save(); err != nil {
    t.Error(err)
  }
  
  if err := counter.Delete(); err != nil {
    t.Error(err)
  }
  
  // make sure the counter is gone
  keys, err := keys("counters")
  if err != nil {
    t.Error(err)
  }
  
  if len(keys) > 0 {
    t.Errorf("counter was not deleted")
  }
}

func TestFindCounterByID(t *testing.T) {
  // clear the counters table when we're finished
  defer truncate("counters")
  
  counter := Counter{
    ID:     "123",
    Count:  1,
  }
  
  if err := counter.Save(); err != nil {
    t.Error(err)
  }
  
  counter2, err := FindCounterByID("123")
  if err != nil {
    t.Error(err)
  }
  
  if counter2.ID != "123" {
    t.Errorf("did not load the correct counter")
  }
}

func TestAllCounters(t *testing.T) {
  // clear the counters table when we're finished
  defer truncate("counters")
  
  counter1 := Counter{ID: "1"}
  counter2 := Counter{ID: "2"}
  counter3 := Counter{ID: "3"}
  
  if err := counter1.Save(); err != nil {
    t.Error(err)
  }
  if err := counter2.Save(); err != nil {
    t.Error(err)
  }
  
  if err := counter3.Save(); err != nil {
    t.Error(err)
  }
  
  counters, err := AllCounters()
  if err != nil {
    t.Error(err)
  }
  
  if len(counters) != 3 {
    t.Errorf("did not load all counters")
  }
}
