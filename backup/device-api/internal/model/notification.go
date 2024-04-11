package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type Lang uint8

const (
	LangEnglish Lang = iota
	LangRussian
	LangEspanol
	LangItalian
)

type NotificationStatus uint8

const (
	NotificationCreated NotificationStatus = iota
	NotificationProcessed
	NotificationDelivered
)

type DayTime uint8

const (
	DayTimeMorning DayTime = iota
	DayTimeDay
	DayTimeEvening
	DayTimeNight
)

//
//
//Language
//
//| Поле         | Описание    |
//|--------------|-------------|
//| LANG_ENGLISH | Английский  |
//| LANG_RUSSIAN | Русский     |
//| LANG_ESPANOL | Испанский   |
//| LANG_ITALIAN | Итальянский |
//
//Status
//
//| Поле               | Описание                        |
//|--------------------|---------------------------------|
//| STATUS_CREATED     | Уведомление создано             |
//| STATUS_IN_PROGRESS | Уведомление в процессе отправки |
//| STATUS_DELIVERED   | Уведомление доставлено          |

// UserNotificationEvent represents a user Event.
type UserNotificationEvent struct {
	DeviceId       uint64 `db:"device_id"         json:"device_id,omitempty"`
	NotificationId uint64 `db:"id"         json:"notification_id,omitempty"`
	Text           string `db:"message"   json:"text,omitempty"`
}

func parseJSONToModel(src interface{}, dest interface{}) error {
	var data []byte

	if b, ok := src.([]byte); ok {
		data = b
	} else if s, ok := src.(string); ok {
		data = []byte(s)
	} else if src == nil {
		return nil
	}

	return json.Unmarshal(data, dest)
}

func (r *UserNotificationEvent) Scan(src interface{}) error {
	return parseJSONToModel(src, r)
}

// NotificationEvent represents a notification event.
type NotificationEvent struct {
	ID                    uint64                 `db:"id"`
	DeviceID              uint64                 `db:"device_id"`
	Lang                  Lang                   `db:"lang"`
	Message               string                 `db:"message"`
	Status                NotificationStatus     `db:"status"`
	UserNotificationEvent *UserNotificationEvent `db:"payload"`
	CreatedAt             time.Time              `db:"created_at"`
	UpdatedAt             time.Time              `db:"updated_at"`
}

func (e *NotificationEvent) BuildUserNotification() {
	e.UserNotificationEvent = e.toUserNotificationEvent()
}
func (e *NotificationEvent) toUserNotificationEvent() *UserNotificationEvent {
	dayTime := getDayTimesOfDay(e.CreatedAt)
	greetMessage := getMessageByDayTimeAndLeng(dayTime, e.Lang)
	return &UserNotificationEvent{
		DeviceId:       e.DeviceID,
		NotificationId: e.ID,
		Text:           fmt.Sprintf("%s %s", greetMessage, e.Message),
	}
}

func getMessageByDayTimeAndLeng(dayTime DayTime, lang Lang) string {
	switch dayTime {
	case DayTimeNight:
		switch lang {
		case LangEnglish:
			return "Good night"
		case LangRussian:
			return "Доброй ночи"
		case LangEspanol:
			return "Buenas noches"
		case LangItalian:
			return "Buona notte"
		}
	case DayTimeDay:
		switch lang {
		case LangEnglish:
			return "Good afternoon"
		case LangRussian:
			return "Добрый день"
		case LangEspanol:
			return "Buenas tardes"
		case LangItalian:
			return "Buon pomeriggio"
		}
	case DayTimeEvening:
		switch lang {
		case LangEnglish:
			return "Good afternoon"
		case LangRussian:
			return "Добрый вечер"
		case LangEspanol:
			return "Buenas noches"
		case LangItalian:
			return "Buona serata"
		}
	case DayTimeMorning:
		switch lang {
		case LangEnglish:
			return "Good morning"
		case LangRussian:
			return "Доброе утро"
		case LangEspanol:
			return "Buenos dias"
		case LangItalian:
			return "Buon giorno"
		}
	}
	return ""
}

func getDayTimesOfDay(at time.Time) DayTime {
	if at.Hour() >= 6 && at.Hour() <= 10 {
		return DayTimeMorning
	} else if at.Hour() >= 11 && at.Hour() <= 14 {
		return DayTimeDay
	} else if at.Hour() >= 15 && at.Hour() <= 20 {
		return DayTimeEvening
	} else {
		return DayTimeNight
	}
}
