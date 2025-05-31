package timeprovider

import "time"

var Current TimeProvider = RealTimeProvider{}

// TimeProvider is an interface for providing time.
type TimeProvider interface {
	Now() time.Time
}

// RealTimeProvider implements TimeProvider using the real time.
type RealTimeProvider struct{}

func (RealTimeProvider) Now() time.Time {
	return time.Now()
}

// mockTimeProvider implements timeprovider.Provider for testing
type MockTimeProvider struct {
	FixedTime time.Time
}

func (m *MockTimeProvider) Now() time.Time {
	return m.FixedTime
}
