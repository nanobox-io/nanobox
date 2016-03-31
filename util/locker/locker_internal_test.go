package locker

import "testing"

func TestLock(t *testing.T) {
	err := Lock()
	if err != nil {
		t.Errorf("unable to aquire lock %+v", err)
	}
	if ln == nil {
		t.Errorf("lock was aqquired but the listener is not populated")
	}
}

func TestTryLock(t *testing.T) {
	locked, err := TryLock()
	if err != nil {
		t.Errorf("lock check failed %+v", err)
	}
	if !locked {
		t.Errorf("lock was not aqquired")
	}
	if ln == nil {
		t.Errorf("lock was aqquired but the listener is not populated")
	}
}

func TestUnlock(t *testing.T) {
	err := Unlock()
	if err != nil {
		t.Errorf("unlock failed %+v", err)
	}
	if ln != nil {
		t.Errorf("unlock succseeded but the listener is still present")
	}
}
