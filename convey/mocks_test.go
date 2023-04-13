package convey

import (
	"io"

	"github.com/stretchr/testify/mock"
)

type mockTranslator struct {
	mock.Mock
}

func (m *mockTranslator) ReadFrom(source io.Reader) (C, error) {
	arguments := m.Called(source)
	return arguments.Get(0).(C), arguments.Error(1)
}

func (m *mockTranslator) WriteTo(destination io.Writer, source C) error {
	return m.Called(destination, source).Error(0)
}
