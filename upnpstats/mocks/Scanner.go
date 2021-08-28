// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	upnpstats "github.com/clambin/upnp-exporter/upnpstats"
	mock "github.com/stretchr/testify/mock"
)

// Scanner is an autogenerated mock type for the Scanner type
type Scanner struct {
	mock.Mock
}

// ReportNetworkStats provides a mock function with given fields:
func (_m *Scanner) ReportNetworkStats() ([]upnpstats.Stats, error) {
	ret := _m.Called()

	var r0 []upnpstats.Stats
	if rf, ok := ret.Get(0).(func() []upnpstats.Stats); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]upnpstats.Stats)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Routers provides a mock function with given fields:
func (_m *Scanner) Routers() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}