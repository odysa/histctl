//go:build darwin

package browser

var knownBrowsers = []string{"safari", "chrome", "edge", "firefox"}

func init() {
	constructors["safari"] = func() Browser { return NewSafari("") }
}
