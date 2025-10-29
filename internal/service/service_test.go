package service_test

import (
	"errors"
	"testing"

	svc "github.com/chernrom/masker/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ---- Моки ----

type MockProducer struct{ mock.Mock }

func (m *MockProducer) Produce() ([]string, error) {
	args := m.Called()
	if v := args.Get(0); v != nil {
		return v.([]string), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockPresenter struct{ mock.Mock }

func (m *MockPresenter) Present(out []string) error {
	args := m.Called(out)
	return args.Error(0)
}

// ---- Тесты ----

func TestRun_HappyPath(t *testing.T) {
	prod := new(MockProducer)
	pres := new(MockPresenter)

	in := []string{
		"see http://abc page",
		"no links here",
		"http://only",
		"https://secure should stay",
		"xhttp://not-a-start",
		"http://a http://bb",
	}
	expected := []string{
		"see http://*** page",
		"no links here",
		"http://****",
		"https://secure should stay",
		"xhttp://***********",
		"http://* http://**",
	}

	prod.On("Produce").Return(in, nil).Once()
	pres.On("Present", expected).Return(nil).Once()

	s := svc.NewService(prod, pres)
	err := s.Run()
	require.NoError(t, err)

	prod.AssertExpectations(t)
	pres.AssertExpectations(t)
}

func TestRun_ProduceError(t *testing.T) {
	prod := new(MockProducer)
	pres := new(MockPresenter)

	prod.On("Produce").Return(nil, errors.New("read-fail")).Once()

	s := svc.NewService(prod, pres)
	err := s.Run()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "produce:")

	pres.AssertNotCalled(t, "Present", mock.Anything)
	prod.AssertExpectations(t)
}

func TestRun_PresentError(t *testing.T) {
	prod := new(MockProducer)
	pres := new(MockPresenter)

	in := []string{"prefix http://abc\nnext"}
	out := []string{"prefix http://***\nnext"}

	prod.On("Produce").Return(in, nil).Once()
	pres.On("Present", out).Return(errors.New("sink")).Once()

	s := svc.NewService(prod, pres)
	err := s.Run()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "present:")

	prod.AssertExpectations(t)
	pres.AssertExpectations(t)
}

func TestRun_EmptyInput(t *testing.T) {
	prod := new(MockProducer)
	pres := new(MockPresenter)

	prod.On("Produce").Return([]string{}, nil).Once()
	pres.On("Present", []string{}).Return(nil).Once()

	s := svc.NewService(prod, pres)
	err := s.Run()
	assert.NoError(t, err)

	prod.AssertExpectations(t)
	pres.AssertExpectations(t)
}

func TestNewService_PanicsOnNilDeps(t *testing.T) {
	assert.Panics(t, func() { svc.NewService(nil, new(MockPresenter)) })
	assert.Panics(t, func() { svc.NewService(new(MockProducer), nil) })
}
