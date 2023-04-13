package device

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/jithin-kg/webpa-common/logging"
	"github.com/xmidt-org/wrp-go/v3"
)

func TestDevice(t *testing.T) {
	var (
		assert              = assert.New(t)
		require             = require.New(t)
		expectedConnectedAt = time.Now().UTC()
		expectedUpTime      = 15 * time.Hour

		testData = []struct {
			expectedID        ID
			expectedQueueSize int
		}{
			{
				ID("ID 1"),
				50,
			},
			{
				ID("ID 2"),
				27,
			},
			{
				ID("ID 3"),
				137,
			},
			{
				ID("ID 4"),
				2,
			},
		}
	)

	for _, record := range testData {
		t.Logf("%v", record)

		var (
			ctx, cancel = context.WithCancel(context.Background())
			testMessage = new(wrp.Message)
			device      = newDevice(deviceOptions{
				ID:          record.expectedID,
				QueueSize:   record.expectedQueueSize,
				ConnectedAt: expectedConnectedAt,
				Logger:      logging.NewTestLogger(nil, t),
			})
		)

		require.NotNil(device)
		assert.NotEmpty(device.sessionID)
		device.statistics = NewStatistics(func() time.Time { return expectedConnectedAt.Add(expectedUpTime) }, expectedConnectedAt)

		assert.Equal(string(record.expectedID), device.String())
		actualConnectedAt := device.Statistics().ConnectedAt()
		assert.Equal(expectedConnectedAt, actualConnectedAt)

		assert.Equal(record.expectedID, device.ID())
		assert.False(device.Closed())

		assert.Equal(record.expectedID, device.ID())
		assert.Equal(actualConnectedAt, device.Statistics().ConnectedAt())
		assert.False(device.Closed())

		data, err := device.MarshalJSON()
		require.NotEmpty(data)
		require.NoError(err)

		assert.JSONEq(
			fmt.Sprintf(
				`{"id": "%s", "pending": 0, "statistics": {"duplications": 0, "bytesSent": 0, "messagesSent": 0, "bytesReceived": 0, "messagesReceived": 0, "connectedAt": "%s", "upTime": "%s"}}`,
				record.expectedID,
				expectedConnectedAt.UTC().Format(time.RFC3339Nano),
				expectedUpTime,
			),
			string(data),
		)

		for repeat := 0; repeat < record.expectedQueueSize; repeat++ {
			go func() {
				request := (&Request{Message: testMessage}).WithContext(ctx)
				device.Send(request)
			}()
		}

		cancel()

		assert.False(device.Closed())
		device.requestClose(CloseReason{Text: "test"})
		assert.True(device.Closed())
		device.requestClose(CloseReason{Text: "test"})
		assert.True(device.Closed())

		response, err := device.Send(&Request{Message: testMessage})
		assert.Nil(response)
		assert.Error(err)
	}
}

func TestDeviceSessionID(t *testing.T) {
	assert := assert.New(t)

	connectOptions := deviceOptions{
		ID:          "1",
		QueueSize:   10,
		ConnectedAt: time.Now(),
		Logger:      logging.NewTestLogger(nil, t),
	}
	sessionOne := newDevice(connectOptions)
	sessionTwo := newDevice(connectOptions)
	assert.Equal(sessionOne.ID(), sessionTwo.ID())
	assert.NotEqual(sessionOne.SessionID(), sessionTwo.SessionID())
}
