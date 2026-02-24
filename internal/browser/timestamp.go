package browser

import "time"

const (
	webkitEpochOffset = 978307200   // seconds between Unix epoch and 2001-01-01
	chromeEpochOffset = 11644473600 // seconds between 1601-01-01 and Unix epoch
)

// WebKitToTime converts a Safari/WebKit timestamp (seconds since 2001-01-01) to time.Time.
func WebKitToTime(webkit float64) time.Time {
	unix := webkit + float64(webkitEpochOffset)
	sec := int64(unix)
	nsec := int64((unix - float64(sec)) * 1e9)
	return time.Unix(sec, nsec)
}

// ChromeToTime converts a Chrome/Edge timestamp (microseconds since 1601-01-01) to time.Time.
func ChromeToTime(ts int64) time.Time {
	unixMicro := ts - chromeEpochOffset*1_000_000
	return time.UnixMicro(unixMicro)
}

// FirefoxToTime converts a Firefox timestamp (microseconds since Unix epoch) to time.Time.
func FirefoxToTime(ts int64) time.Time {
	return time.UnixMicro(ts)
}
