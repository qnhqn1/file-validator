

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	pgxpool "github.com/jackc/pgx/v5/pgxpool"
)


type MockStorageInterface struct {
	mock.Mock
}

type MockStorageInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockStorageInterface) EXPECT() *MockStorageInterface_Expecter {
	return &MockStorageInterface_Expecter{mock: &_m.Mock}
}


func (_m *MockStorageInterface) Close() {
	_m.Called()
}


type MockStorageInterface_Close_Call struct {
	*mock.Call
}


func (_e *MockStorageInterface_Expecter) Close() *MockStorageInterface_Close_Call {
	return &MockStorageInterface_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *MockStorageInterface_Close_Call) Run(run func()) *MockStorageInterface_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockStorageInterface_Close_Call) Return() *MockStorageInterface_Close_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockStorageInterface_Close_Call) RunAndReturn(run func()) *MockStorageInterface_Close_Call {
	_c.Run(run)
	return _c
}


func (_m *MockStorageInterface) InsertEvent(ctx context.Context, key string, payload []byte) error {
	ret := _m.Called(ctx, key, payload)

	if len(ret) == 0 {
		panic("no return value specified for InsertEvent")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []byte) error); ok {
		r0 = rf(ctx, key, payload)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}


type MockStorageInterface_InsertEvent_Call struct {
	*mock.Call
}





func (_e *MockStorageInterface_Expecter) InsertEvent(ctx interface{}, key interface{}, payload interface{}) *MockStorageInterface_InsertEvent_Call {
	return &MockStorageInterface_InsertEvent_Call{Call: _e.mock.On("InsertEvent", ctx, key, payload)}
}

func (_c *MockStorageInterface_InsertEvent_Call) Run(run func(ctx context.Context, key string, payload []byte)) *MockStorageInterface_InsertEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].([]byte))
	})
	return _c
}

func (_c *MockStorageInterface_InsertEvent_Call) Return(_a0 error) *MockStorageInterface_InsertEvent_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockStorageInterface_InsertEvent_Call) RunAndReturn(run func(context.Context, string, []byte) error) *MockStorageInterface_InsertEvent_Call {
	_c.Call.Return(run)
	return _c
}


func (_m *MockStorageInterface) PrimaryPool() *pgxpool.Pool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for PrimaryPool")
	}

	var r0 *pgxpool.Pool
	if rf, ok := ret.Get(0).(func() *pgxpool.Pool); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pgxpool.Pool)
		}
	}

	return r0
}


type MockStorageInterface_PrimaryPool_Call struct {
	*mock.Call
}


func (_e *MockStorageInterface_Expecter) PrimaryPool() *MockStorageInterface_PrimaryPool_Call {
	return &MockStorageInterface_PrimaryPool_Call{Call: _e.mock.On("PrimaryPool")}
}

func (_c *MockStorageInterface_PrimaryPool_Call) Run(run func()) *MockStorageInterface_PrimaryPool_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockStorageInterface_PrimaryPool_Call) Return(_a0 *pgxpool.Pool) *MockStorageInterface_PrimaryPool_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockStorageInterface_PrimaryPool_Call) RunAndReturn(run func() *pgxpool.Pool) *MockStorageInterface_PrimaryPool_Call {
	_c.Call.Return(run)
	return _c
}



func NewMockStorageInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockStorageInterface {
	mock := &MockStorageInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}


