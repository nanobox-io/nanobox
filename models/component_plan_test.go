package models

import (
  "testing"
)

func TestComponentPlanBehaviors(t *testing.T) {
  p1 := ComponentPlan{
    Behaviors: []string{"backupable", "migratable"},
  }
  
  p2 := ComponentPlan{}
  
  if !p1.BehaviorPresent("backupable") {
    t.Errorf("behavior doesn't exist and should")
  }
  
  if p2.BehaviorPresent("backupable") {
    t.Errorf("behavior exists and should not")
  }
}
