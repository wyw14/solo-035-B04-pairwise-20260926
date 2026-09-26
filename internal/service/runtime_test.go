package service

import (
	"testing"
	"time"
)

func TestBootstrapBuildsOperationalSnapshot(t *testing.T) {
	runtime, err := NewRuntime(t.TempDir(), func() time.Time { return time.Unix(100, 0).UTC() })
	if err != nil {
		t.Fatal(err)
	}
	if err := Bootstrap(runtime); err != nil {
		t.Fatal(err)
	}
	snapshot := runtime.Snapshot()
	if len(snapshot.Operations) != 1 {
		t.Fatalf("operations = %d", len(snapshot.Operations))
	}
	if len(snapshot.Equipment) < 3 {
		t.Fatalf("equipment = %d", len(snapshot.Equipment))
	}
}
