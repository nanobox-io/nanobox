package evar

import (
	"testing"
)

// TestEvarAdd tests the adding of environment variables from the cli.
func TestEvarAdd(t *testing.T) {
	vars := []string{`key="this
is
a
multiline
"`, "key2=val", "key3=\"this\nis\na\nmultiline too\""}
	evars := parseEvars(vars)

	if len(evars) != 3 {
		t.Fatalf("Failed to parse all evars - %d - %q", len(evars), evars)
	}

	if evars["KEY"] != "this\nis\na\nmultiline\n" {
		t.Fatalf("multiline var failed - %q", evars["KEY"])
	}

	if evars["KEY2"] != "val" {
		t.Fatalf("Commas, spaces, = var failed - %q", evars["KEY2"])
	}

	if evars["KEY3"] != "this\nis\na\nmultiline too" {
		t.Fatalf("Single quote, semicolon var failed - %q", evars["KEY3"])
	}
}

// TestEvarLoad tests the loading of environment variables from a file.
func TestEvarLoad(t *testing.T) {
	vars, _ := loadVars([]string{""}, testGetter{})
	evars := parseEvars(vars)

	if len(evars) != 13 {
		t.Fatalf("Failed to parse all evars - %d - %q", len(evars), evars)
	}

	if evars["KEY4"] != "\nanother\nmultiline\n" {
		t.Fatalf("multiline var failed - %q", evars["KEY4"])
	}

	if evars["KEY5"] != "yes, even spaces and = are allowed as values" {
		t.Fatalf("Commas, spaces, = var failed - %q", evars["KEY5"])
	}

	if evars["KEY9"] != "you're \"welcome ;)" {
		t.Fatalf("Single quote, semicolon var failed - %q", evars["KEY9"])
	}

	if evars["KEY_11"] != "x\ny\nz" {
		t.Fatalf("Underscored key failed - %q", evars["KEY_11"])
	}

	if evars["KEY_12"] != "this\none\nhas\nan=in\nthe\nmiddle" {
		t.Fatalf("Multiline var with equal sign in value failed to parse - %q should be \"this\none\nhas\nan=in\nthe\nmiddle\"", evars["KEY_12"])
	}

	if evars["KEY_13"] != "previous value with equal sign isn't greedy and leaves me alone" {
		t.Fatalf("Variable(s) after one with equal sign in middle failed to parse - %q should be \"previous value with equal sign isn't greedy and leaves me alone\"", evars["KEY_13"])
	}
}

type testGetter struct{}

func (f testGetter) getContents(filename string) (string, error) {
	return testContents, nil
}

var testContents = `# comment
key1=val
key2="val"
key3="this
is
a
multiline
value"
key4="
another
multiline
"
key5="yes, even spaces and = are allowed as values"
export key6=gasp
export key7="how is this guy doing these awesome things"

# comment
# more comment

key8="yep, even whitespace is _allowed (gets stripped)"
export key9="you're \"welcome ;)"
key10="x"
key_11="x
y
z"
key_12="this
one
has
an=in
the
middle"
key_13="previous value with equal sign isn't greedy and leaves me alone"
`
