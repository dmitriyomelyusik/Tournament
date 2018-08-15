// Code generated by mockery v1.0.0. DO NOT EDIT.
package handlers

import entity "github.com/dmitriyomelyusik/Tournament/entity"
import mock "github.com/stretchr/testify/mock"

// mockCtlr is an autogenerated mock type for the ctlr type
type mockCtlr struct {
	mock.Mock
}

// AnnounceTournament provides a mock function with given fields: id, deposit
func (_m *mockCtlr) AnnounceTournament(id string, deposit int) error {
	ret := _m.Called(id, deposit)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int) error); ok {
		r0 = rf(id, deposit)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Balance provides a mock function with given fields: id
func (_m *mockCtlr) Balance(id string) (entity.Player, error) {
	ret := _m.Called(id)

	var r0 entity.Player
	if rf, ok := ret.Get(0).(func(string) entity.Player); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(entity.Player)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Fund provides a mock function with given fields: id, points
func (_m *mockCtlr) Fund(id string, points int) (entity.Player, error) {
	ret := _m.Called(id, points)

	var r0 entity.Player
	if rf, ok := ret.Get(0).(func(string, int) entity.Player); ok {
		r0 = rf(id, points)
	} else {
		r0 = ret.Get(0).(entity.Player)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, int) error); ok {
		r1 = rf(id, points)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// JoinTournament provides a mock function with given fields: tourID, playerID
func (_m *mockCtlr) JoinTournament(tourID string, playerID string) error {
	ret := _m.Called(tourID, playerID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(tourID, playerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Results provides a mock function with given fields: tourID
func (_m *mockCtlr) Results(tourID string) (entity.Winners, error) {
	ret := _m.Called(tourID)

	var r0 entity.Winners
	if rf, ok := ret.Get(0).(func(string) entity.Winners); ok {
		r0 = rf(tourID)
	} else {
		r0 = ret.Get(0).(entity.Winners)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tourID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Take provides a mock function with given fields: id, points
func (_m *mockCtlr) Take(id string, points int) error {
	ret := _m.Called(id, points)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int) error); ok {
		r0 = rf(id, points)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}