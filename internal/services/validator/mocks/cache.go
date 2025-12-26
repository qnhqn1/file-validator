package mocks

import (
	context "context"
	time "time"

	mock "github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

type MockCache_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCache) EXPECT() *MockCache_Expecter {
	return &MockCache_Expecter{mock: &_m.Mock}
}

func (_m *MockCache) Close() {
	_m.Called()
}

type MockCache_Close_Call struct {
	*mock.Call
}

func (_e *MockCache_Expecter) Close() *MockCache_Close_Call {
	return &MockCache_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *MockCache_Close_Call) Run(run func()) *MockCache_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCache_Close_Call) Return() *MockCache_Close_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockCache_Close_Call) RunAndReturn(run func()) *MockCache_Close_Call {
	_c.Run(run)
	return _c
}

func (_m *MockCache) Get(ctx context.Context, key string) ([]byte, error) {
	ret := _m.Called(ctx, key)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]byte, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []byte); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type MockCache_Get_Call struct {
	*mock.Call
}

func (_e *MockCache_Expecter) Get(ctx interface{}, key interface{}) *MockCache_Get_Call {
	return &MockCache_Get_Call{Call: _e.mock.On("Get", ctx, key)}
}

func (_c *MockCache_Get_Call) Run(run func(ctx context.Context, key string)) *MockCache_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockCache_Get_Call) Return(_a0 []byte, _a1 error) *MockCache_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCache_Get_Call) RunAndReturn(run func(context.Context, string) ([]byte, error)) *MockCache_Get_Call {
	_c.Call.Return(run)
	return _c
}

func (_m *MockCache) Set(ctx context.Context, key string, val []byte, ttl time.Duration) error {
	ret := _m.Called(ctx, key, val, ttl)

	if len(ret) == 0 {
		panic("no return value specified for Set")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, []byte, time.Duration) error); ok {
		r0 = rf(ctx, key, val, ttl)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type MockCache_Set_Call struct {
	*mock.Call
}

func (_e *MockCache_Expecter) Set(ctx interface{}, key interface{}, val interface{}, ttl interface{}) *MockCache_Set_Call {
	return &MockCache_Set_Call{Call: _e.mock.On("Set", ctx, key, val, ttl)}
}

func (_c *MockCache_Set_Call) Run(run func(ctx context.Context, key string, val []byte, ttl time.Duration)) *MockCache_Set_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].([]byte), args[3].(time.Duration))
	})
	return _c
}

func (_c *MockCache_Set_Call) Return(_a0 error) *MockCache_Set_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCache_Set_Call) RunAndReturn(run func(context.Context, string, []byte, time.Duration) error) *MockCache_Set_Call {
	_c.Call.Return(run)
	return _c
}

func NewMockCache(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCache {
	mock := &MockCache{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
