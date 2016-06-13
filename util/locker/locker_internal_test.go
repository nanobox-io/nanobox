package locker

import "testing"

// TestLocalLock ...
func TestLocalLock(t *testing.T) {
	err := LocalLock()
	if err != nil {
		t.Errorf("unable to aquire lock %+v", err)
	}
	if lln == nil {
		t.Errorf("lock was aqquired but the listener is not populated")
	}
}

// TestLocalTryLock ...
func TestLocalTryLock(t *testing.T) {
	locked, err := LocalTryLock()
	if err != nil {
		t.Errorf("lock check failed %+v", err)
	}
	if !locked {
		t.Errorf("lock was not aqquired")
	}
	if lln == nil {
		t.Errorf("lock was aqquired but the listener is not populated")
	}
}

// TestLocalUnlock ...
func TestLocalUnlock(t *testing.T) {
	err := LocalUnlock()
	if err != nil {
		t.Errorf("unlock failed %+v", err)
	}
	if lln != nil {
		t.Errorf("unlock succseeded but the listener is still present")
	}
}

// TestLocalStackLocking ...
func TestLocalStackLocking(t *testing.T) {
	for i := 0; i < 10; i++ {
		if lCount != i {
			t.Errorf("global count not equil to lock calls(%d)", gCount)
		}
		err := LocalLock()
		if err != nil {
			t.Errorf("errored on multiple locks")
		}
	}
	for i := 10; i > 0; i-- {
		if lCount != i {
			t.Errorf("global count not equil to lock calls(%d)", gCount)
		}
		err := LocalUnlock()
		if err != nil {
			t.Errorf("errored on multiple locks")
		}
	}
}

// TestGlobalLock ...
func TestGlobalLock(t *testing.T) {
	err := GlobalLock()
	if err != nil {
		t.Errorf("unable to aquire lock %+v", err)
	}
	if gln == nil {
		t.Errorf("lock was aqquired but the listener is not populated")
	}
}

// TestGlobalTryLock ...
func TestGlobalTryLock(t *testing.T) {
	locked, err := GlobalTryLock()
	if err != nil {
		t.Errorf("lock check failed %+v", err)
	}
	if !locked {
		t.Errorf("lock was not aqquired")
	}
	if gln == nil {
		t.Errorf("lock was aqquired but the listener is not populated")
	}
}

// TestGlobalUnlock ...
func TestGlobalUnlock(t *testing.T) {
	err := GlobalUnlock()
	if err != nil {
		t.Errorf("unlock failed %+v", err)
	}
	if gln != nil {
		t.Errorf("unlock succseeded but the listener is still present")
	}
}

// TestGlobalStackLocking ...
func TestGlobalStackLocking(t *testing.T) {
	for i := 0; i < 10; i++ {
		if gCount != i {
			t.Errorf("global count not equil to lock calls(%d)", gCount)
		}
		err := GlobalLock()
		if err != nil {
			t.Errorf("errored on multiple locks")
		}
	}
	for i := 10; i > 0; i-- {
		if gCount != i {
			t.Errorf("global count not equil to lock calls(%d)", gCount)
		}
		err := GlobalUnlock()
		if err != nil {
			t.Errorf("errored on multiple locks")
		}
	}
}
