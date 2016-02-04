package engine_test

import (
	"fmt"
	"testing"

	engineutil "github.com/nanobox-io/nanobox/util/engine"
)

// TestRemountLocal
func TestRemountLocal(t *testing.T) {
}

// TestMountLocal
func TestMountLocal(t *testing.T) {
}

// TestCreate
func TestCreate(t *testing.T) {
}

// TestGet
func TestGet(t *testing.T) {
}

// TestParseArchive
func TestParseArchive(t *testing.T) {
	archive := "engine-name"

	user, engine := engineutil.ParseArchive(archive)

	if user != "" {
		t.Error(fmt.Sprintf("Expected nothing got '%s'!", user))
	}

	if engine != "engine-name" {
		t.Error(fmt.Sprintf("Expected 'engine-name' got '%s'!", engine))
	}
}

// TestParseArchiveWVersion
func TestParseArchiveWVersion(t *testing.T) {
	archive := "engine-name=0.0.1"

	user, engine := engineutil.ParseArchive(archive)

	if user != "" {
		t.Error(fmt.Sprintf("Expected 'user' got '%s'!", user))
	}

	if engine != "engine-name=0.0.1" {
		t.Error(fmt.Sprintf("Expected 'engine-name=0.0.1' got '%s'!", engine))
	}
}

// TestParseArchiveWUser
func TestParseArchiveWUser(t *testing.T) {
	archive := "user/engine-name"

	user, engine := engineutil.ParseArchive(archive)

	if user != "user" {
		t.Error(fmt.Sprintf("Expected 'user' got '%s'!", user))
	}

	if engine != "engine-name" {
		t.Error(fmt.Sprintf("Expected 'engine-name' got '%s'!", engine))
	}
}

// TestParseArchiveWUserAndVersion
func TestParseArchiveWUserAndVersion(t *testing.T) {
	archive := "user/engine-name=0.0.1"

	user, engine := engineutil.ParseArchive(archive)

	if user != "user" {
		t.Error(fmt.Sprintf("Expected 'user' got '%s'!", user))
	}

	if engine != "engine-name=0.0.1" {
		t.Error(fmt.Sprintf("Expected 'engine-name=0.0.1' got '%s'!", engine))
	}
}

// TestParseEngine
func TestParseEngine(t *testing.T) {
	engine := "engine-name"

	name, version := engineutil.ParseEngine(engine)

	if name != "engine-name" {
		t.Error(fmt.Sprintf("Expected 'engine-name' got %s!", name))
	}

	if version != "" {
		t.Error(fmt.Sprintf("Expected nothing got '%s'!", version))
	}
}

// TestParseEngineWVersion
func TestParseEngineWVersion(t *testing.T) {
	engine := "engine-name=0.0.1"

	name, version := engineutil.ParseEngine(engine)

	if name != "engine-name" {
		t.Error(fmt.Sprintf("Expected 'engine-name' got '%s'!", name))
	}

	if version != "0.0.1" {
		t.Error(fmt.Sprintf("Expected '0.0.1' got '%s'!", version))
	}
}
