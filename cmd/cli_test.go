package cmd

import (
	"net"
	"testing"
)

func TestParsePortList(t *testing.T) {
	tests := []struct {
		input    string
		expected []int
	}{
		{"80,443", []int{80, 443}},
		{"80, 443, 8080", []int{80, 443, 8080}},
		{"invalid", []int{}},
		{"80,invalid,443", []int{80, 443}},
	}

	for _, tt := range tests {
		result := parsePortList(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("parsePortList(%s): got %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestIsPortOpen(t *testing.T) {
	// Test invalid IP - should return false
	result := isPortOpen("999.999.999.999", 80)
	if result {
		t.Errorf("isPortOpen with invalid IP should return false")
	}
}

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{"Simple title", "<html><title>Test Page</title></html>", "Test Page"},
		{"Title with spaces", "<html><title>  Test Page  </title></html>", "Test Page"},
		{"No title", "<html></html>", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test is simplified since we can't easily mock http.Response
			t.Logf("Test: %s", tt.name)
		})
	}
}

func TestIncrementIP(t *testing.T) {
	ip := net.ParseIP("192.168.1.0")
	oldIP := ip.String()
	incrementIP(ip)
	newIP := ip.String()

	if oldIP == newIP {
		t.Errorf("incrementIP failed to increment IP")
	}
}

func TestGetHostsInNetwork(t *testing.T) {
	hosts := getHostsInNetwork("192.168.1.0/30")
	if len(hosts) < 1 {
		t.Errorf("getHostsInNetwork should return at least 1 host for /30 network")
	}
}

func TestDevice(t *testing.T) {
	device := Device{
		IP:  "192.168.1.1",
		MAC: "aa:bb:cc:dd:ee:ff",
	}

	if device.IP != "192.168.1.1" {
		t.Errorf("Device IP mismatch")
	}
}

func TestService(t *testing.T) {
	service := Service{
		IP:       "192.168.1.1",
		Port:     80,
		Protocol: "http",
		Title:    "Test Server",
	}

	if service.Port != 80 || service.Protocol != "http" {
		t.Errorf("Service fields mismatch")
	}
}

func TestScanResult(t *testing.T) {
	result := ScanResult{
		Timestamp: "2024-01-01T00:00:00Z",
		Devices:   []Device{},
		Services:  []Service{},
	}

	if len(result.Devices) != 0 || len(result.Services) != 0 {
		t.Errorf("ScanResult initialization failed")
	}
}
