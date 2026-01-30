package uid

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/google/uuid"
)

// ID نوع داده‌ای ماست که جایگزین uuid.UUID می‌شود
type ID uuid.UUID

var Nil = ID(uuid.Nil)

// New یک ID تصادفی جدید می‌سازد
func New() ID {
	return ID(uuid.New())
}

// Parse رشته را به ID تبدیل می‌کند
func Parse(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return ID(id), err
}

// متد String برای تبدیل به رشته
func (id ID) String() string {
	return uuid.UUID(id).String()
}

// پیاده‌سازی رابط Scanner برای خواندن از دیتابیس (SQL)
func (id *ID) Scan(src interface{}) error {
	var u uuid.UUID
	if err := u.Scan(src); err != nil {
		return err
	}
	*id = ID(u)
	return nil
}

// پیاده‌سازی رابط Valuer برای نوشتن در دیتابیس (SQL)
func (id ID) Value() (driver.Value, error) {
	return uuid.UUID(id).Value()
}

// پیاده‌سازی رابط Marshaler برای JSON (تبدیل به رشته در خروجی API)
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(id).String())
}

// پیاده‌سازی رابط Unmarshaler برای JSON (خواندن رشته از ورودی API)
func (id *ID) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parsed, err := uuid.Parse(s)
	if err != nil {
		return err
	}
	*id = ID(parsed)
	return nil
}