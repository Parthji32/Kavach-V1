package services

import (
	"crypto/md5"
	"fmt"
	"strings"
)

// FingerprintGenerator generates device fingerprints
type FingerprintGenerator struct{}

// Generate creates a fingerprint from request metadata
func (fg *FingerprintGenerator) Generate(ip, userAgent, language, encoding string) string {
	data := fmt.Sprintf("%s|%s|%s|%s", ip, userAgent, language, encoding)
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// ParseUserAgent extracts device information from User-Agent
func (fg *FingerprintGenerator) ParseUserAgent(userAgent string) map[string]string {
	info := make(map[string]string)

	// Detect OS
	if strings.Contains(userAgent, "Windows") {
		info["os"] = "Windows"
	} else if strings.Contains(userAgent, "Mac") {
		info["os"] = "macOS"
	} else if strings.Contains(userAgent, "Linux") {
		info["os"] = "Linux"
	} else if strings.Contains(userAgent, "iPhone") {
		info["os"] = "iOS"
	} else if strings.Contains(userAgent, "Android") {
		info["os"] = "Android"
	} else {
		info["os"] = "Unknown"
	}

	// Detect Browser
	if strings.Contains(userAgent, "Chrome") && !strings.Contains(userAgent, "Chromium") {
		info["browser"] = "Chrome"
	} else if strings.Contains(userAgent, "Firefox") {
		info["browser"] = "Firefox"
	} else if strings.Contains(userAgent, "Safari") && !strings.Contains(userAgent, "Chrome") {
		info["browser"] = "Safari"
	} else if strings.Contains(userAgent, "Edge") {
		info["browser"] = "Edge"
	} else if strings.Contains(userAgent, "curl") {
		info["browser"] = "curl"
	} else if strings.Contains(userAgent, "wget") {
		info["browser"] = "wget"
	} else if strings.Contains(userAgent, "python") {
		info["browser"] = "Python"
	} else {
		info["browser"] = "Unknown"
	}

	// Detect Device Type
	if strings.Contains(userAgent, "Mobile") || strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "Android") {
		info["device_type"] = "mobile"
	} else if strings.Contains(userAgent, "Tablet") || strings.Contains(userAgent, "iPad") {
		info["device_type"] = "tablet"
	} else if strings.Contains(userAgent, "curl") || strings.Contains(userAgent, "wget") || strings.Contains(userAgent, "python") {
		info["device_type"] = "bot"
	} else {
		info["device_type"] = "desktop"
	}

	return info
}
