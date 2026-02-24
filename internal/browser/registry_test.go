package browser

import "testing"

func TestGet(t *testing.T) {
	for _, name := range knownBrowsers {
		b, err := Get(name)
		if err != nil {
			t.Fatalf("Get(%q) returned error: %v", name, err)
		}
		if b.Name() != name {
			t.Errorf("Get(%q).Name() = %q", name, b.Name())
		}
	}
}

func TestGetUnknown(t *testing.T) {
	_, err := Get("netscape")
	if err == nil {
		t.Fatal("Get(\"netscape\") should return an error")
	}
}

func TestAll(t *testing.T) {
	browsers := All()
	if len(browsers) != len(knownBrowsers) {
		t.Fatalf("All() returned %d browsers, want %d", len(browsers), len(knownBrowsers))
	}
	names := map[string]bool{}
	for _, b := range browsers {
		names[b.Name()] = true
	}
	for _, want := range knownBrowsers {
		if !names[want] {
			t.Errorf("All() missing %q", want)
		}
	}
}

func TestNames(t *testing.T) {
	names := Names()
	if len(names) != len(knownBrowsers) {
		t.Fatalf("Names() returned %d names, want %d", len(names), len(knownBrowsers))
	}
	// Verify it returns a copy
	names[0] = "modified"
	if Names()[0] == "modified" {
		t.Error("Names() should return a copy, not the original slice")
	}
}
