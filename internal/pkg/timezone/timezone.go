package timezone

import "time"

// UserLocation возвращает таймзону пользователя.
// TODO: в будущем можно брать из профиля пользователя.
var UserLocation = func() *time.Location {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return time.FixedZone("MSK", 3*60*60)
	}
	return loc
}()

// ToUser конвертирует UTC-время в локальное время пользователя.
func ToUser(t time.Time) time.Time {
	return t.In(UserLocation)
}

// FromUser парсит строку datetime-local (без зоны) как локальное время пользователя.
func FromUser(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, UserLocation)
}

// FormatUser форматирует время в локальной зоне пользователя.
func FormatUser(t time.Time, layout string) string {
	return t.In(UserLocation).Format(layout)
}
