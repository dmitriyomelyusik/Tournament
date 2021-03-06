// Code generated by mockery v1.0.0. DO NOT EDIT.
package controller

import entity "github.com/dmitriyomelyusik/Tournament/entity"
import mock "github.com/stretchr/testify/mock"

// MockDatabase is an autogenerated mock type for the Database type
type MockDatabase struct {
	mock.Mock
}

// CloseTournament provides a mock function with given fields: id
func (_m *MockDatabase) CloseTournament(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreatePlayer provides a mock function with given fields: id, points
func (_m *MockDatabase) CreatePlayer(id string, points int) (entity.Player, error) {
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

// CreateTournament provides a mock function with given fields: id, deposit
func (_m *MockDatabase) CreateTournament(id string, deposit int) error {
	ret := _m.Called(id, deposit)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int) error); ok {
		r0 = rf(id, deposit)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetParticipants provides a mock function with given fields: id
func (_m *MockDatabase) GetParticipants(id string) ([]string, error) {
	ret := _m.Called(id)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPlayer provides a mock function with given fields: id
func (_m *MockDatabase) GetPlayer(id string) (entity.Player, error) {
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

// GetTournamentState provides a mock function with given fields: id
func (_m *MockDatabase) GetTournamentState(id string) (bool, error) {
	ret := _m.Called(id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWinner provides a mock function with given fields: id
func (_m *MockDatabase) GetWinner(id string) (entity.Winners, error) {
	ret := _m.Called(id)

	var r0 entity.Winners
	if rf, ok := ret.Get(0).(func(string) entity.Winners); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(entity.Winners)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetTournamentWinner provides a mock function with given fields: id, winner
func (_m *MockDatabase) SetTournamentWinner(id string, winner entity.Winner) error {
	ret := _m.Called(id, winner)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, entity.Winner) error); ok {
		r0 = rf(id, winner)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePlayer provides a mock function with given fields: id, dif
func (_m *MockDatabase) UpdatePlayer(id string, dif int) error {
	ret := _m.Called(id, dif)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, int) error); ok {
		r0 = rf(id, dif)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTourAndPlayer provides a mock function with given fields: tourID, playerID
func (_m *MockDatabase) UpdateTourAndPlayer(tourID string, playerID string) error {
	ret := _m.Called(tourID, playerID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(tourID, playerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
