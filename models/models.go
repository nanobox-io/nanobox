// Package models ...
package models

// BehaviorPresent ...
func (p Plan) BehaviorPresent(b string) bool {
	for _, behavior := range p.Behaviors {
		if behavior == b {
			return true
		}
	}

	return false
}
