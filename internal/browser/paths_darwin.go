//go:build darwin

package browser

func init() {
	knownBrowsers = append([]string{"safari"}, knownBrowsers...)
	constructors["safari"] = func() Browser { return NewSafari("") }
}
