package ua

import "strings"

// Device describes a client in human-readable form, e.g. "Chrome on macOS".
type Device struct {
	Browser string
	OS      string
}

func (d Device) String() string {
	switch {
	case d.Browser != "" && d.OS != "":
		return d.Browser + " on " + d.OS
	case d.Browser != "":
		return d.Browser
	case d.OS != "":
		return d.OS
	default:
		return "Unknown device"
	}
}

// Parse extracts a rough browser + OS label from a User-Agent string.
func Parse(userAgent string) Device {
	s := userAgent
	if strings.TrimSpace(s) == "" {
		return Device{}
	}
	lower := strings.ToLower(s)

	var d Device
	switch {
	case strings.Contains(lower, "edg/") || strings.Contains(lower, "edge/"):
		d.Browser = "Edge"
	case strings.Contains(lower, "opr/") || strings.Contains(lower, "opera"):
		d.Browser = "Opera"
	case strings.Contains(lower, "firefox/"):
		d.Browser = "Firefox"
	case strings.Contains(lower, "chrome/") || strings.Contains(lower, "crios/"):
		d.Browser = "Chrome"
	case strings.Contains(lower, "safari/") && strings.Contains(lower, "version/"):
		d.Browser = "Safari"
	case strings.Contains(lower, "curl/"):
		d.Browser = "curl"
	case strings.Contains(lower, "postman"):
		d.Browser = "Postman"
	case strings.Contains(lower, "python"):
		d.Browser = "Python"
	case strings.Contains(lower, "go-http-client"):
		d.Browser = "Go client"
	}

	switch {
	case strings.Contains(lower, "iphone") || strings.Contains(lower, "ipad"):
		d.OS = "iOS"
	case strings.Contains(lower, "mac os x") || strings.Contains(lower, "macintosh"):
		d.OS = "macOS"
	case strings.Contains(lower, "android"):
		d.OS = "Android"
	case strings.Contains(lower, "windows"):
		d.OS = "Windows"
	case strings.Contains(lower, "linux") || strings.Contains(lower, "x11"):
		d.OS = "Linux"
	}
	return d
}
