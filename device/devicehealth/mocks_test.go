package devicehealth

import (
	"github.com/stretchr/testify/mock"
	"github.com/jithin-kg/webpa-common/health"
)

type mockDispatcher struct {
	mock.Mock
}

func (m *mockDispatcher) SendEvent(hf health.HealthFunc) {
	m.Called(hf)
}
