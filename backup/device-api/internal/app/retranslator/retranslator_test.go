package retranslator

import (
	"errors"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/app/repo"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/app/sender"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/mocks"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/model"
)

func TestStart(t *testing.T) {

	ctrl := gomock.NewController(t)
	repoMock := mocks.NewMockEventRepo(ctrl)
	senderMock := mocks.NewMockEventSender(ctrl)

	repoMock.EXPECT().Lock(gomock.Any(), gomock.Any()).AnyTimes()

	cfg := getConfig(repoMock, senderMock)

	retranslatorStartAndClose(cfg)
}

func TestRemoveAfterSendSuccess(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := mocks.NewMockEventRepo(ctrl)
	senderMock := mocks.NewMockEventSender(ctrl)

	repoMock.EXPECT().Lock(gomock.Any(), gomock.Any()).AnyTimes()

	cfg := getConfig(repoMock, senderMock)

	events := getEvents()

	repoMock.EXPECT().Lock(gomock.Any(), uint64(10)).Return(events, nil).MinTimes(1).MaxTimes(1)
	senderMock.EXPECT().Send(&events[0]).Return(nil).MinTimes(1).MaxTimes(1)
	repoMock.EXPECT().Remove(gomock.Any(), []uint64{events[0].ID}).Return(nil).MinTimes(1).MaxTimes(1)

	retranslatorStartAndClose(cfg)
}

func TestUnlockAfterSendError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repoMock := mocks.NewMockEventRepo(ctrl)
	senderMock := mocks.NewMockEventSender(ctrl)

	repoMock.EXPECT().Lock(gomock.Any(), gomock.Any()).AnyTimes()

	cfg := getConfig(repoMock, senderMock)

	events := getEvents()

	repoMock.EXPECT().Lock(gomock.Any(), uint64(10)).Return(events, nil).MinTimes(1).MaxTimes(1)
	senderMock.EXPECT().Send(&events[0]).Return(errors.New("event sending failed")).MinTimes(1).MaxTimes(1)
	repoMock.EXPECT().Unlock(gomock.Any(), []uint64{events[0].ID}).Return(nil).MinTimes(1).MaxTimes(1)

	retranslatorStartAndClose(cfg)
}

func getConfig(repo repo.EventRepo, sender sender.EventSender) Config {
	return Config{
		ChannelSize:    512,
		ConsumerCount:  2,
		ConsumeSize:    10,
		ConsumeTimeout: 10 * time.Second,
		ProducerCount:  2,
		WorkerCount:    2,
		Repo:           repo,
		Sender:         sender,
	}
}

func getEvents() []model.DeviceEvent {
	t := time.Now()
	return []model.DeviceEvent{
		{
			ID:     1,
			Type:   model.Created,
			Status: model.Processed,
			Device: &model.Device{
				ID:        1,
				Platform:  "Android",
				UserID:    123456,
				EnteredAt: &t,
			},
		},
	}
}

func retranslatorStartAndClose(cfg Config) {
	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	time.Sleep(time.Second)
	retranslator.Close()
}
