package browser

import "fmt"

// knownBrowsers is defined in registry_darwin.go and registry_other.go

var constructors = map[string]func() Browser{
	"chrome":  func() Browser { return NewChrome("") },
	"edge":    func() Browser { return NewEdge("") },
	"firefox": func() Browser { return NewFirefox("") },
}

// Get returns a Browser by name.
func Get(name string) (Browser, error) {
	if ctor, ok := constructors[name]; ok {
		return ctor(), nil
	}
	return nil, fmt.Errorf("unknown browser %q; supported: %v", name, knownBrowsers)
}

// All returns Browser instances for every known browser.
func All() []Browser {
	out := make([]Browser, 0, len(knownBrowsers))
	for _, name := range knownBrowsers {
		out = append(out, constructors[name]())
	}
	return out
}

// Available returns browsers that are actually installed (DB file exists).
func Available() []Browser {
	var out []Browser
	for _, b := range All() {
		if _, err := b.DBPath(); err == nil {
			out = append(out, b)
		}
	}
	return out
}

// Names returns the list of supported browser names.
func Names() []string {
	return append([]string{}, knownBrowsers...)
}
