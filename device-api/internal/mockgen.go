package internal

//go:generate mockgen -destination=./mocks/repo_mock.go -package=mocks gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/app/repo EventRepo
//go:generate mockgen -destination=./mocks/sender_mock.go -package=mocks gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/app/sender EventSender
