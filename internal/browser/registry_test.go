package browser

import "testing"

func TestGet(t *testing.T) {
	for _, name := range []string{"safari", "chrome", "edge", "firefox"} {
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
	if len(browsers) != 4 {
		t.Fatalf("All() returned %d browsers, want 4", len(browsers))
	}
	names := map[string]bool{}
	for _, b := range browsers {
		names[b.Name()] = true
	}
	for _, want := range []string{"safari", "chrome", "edge", "firefox"} {
		if !names[want] {
			t.Errorf("All() missing %q", want)
		}
	}
}

func TestNames(t *testing.T) {
	names := Names()
	if len(names) != 4 {
		t.Fatalf("Names() returned %d names, want 4", len(names))
	}
	// Verify it returns a copy
	names[0] = "modified"
	if Names()[0] == "modified" {
		t.Error("Names() should return a copy, not the original slice")
	}
}
