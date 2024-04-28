// Code generated by cmd/gen-types; DO NOT EDIT.
package types

import (
	"database/sql/driver"
	"errors"

	"github.com/google/uuid"
)

func Parse[T UserID | RequestID | TenantID](id string) (T, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return T(uuid.Nil), err
	}
	return T(uid), nil
}

func MustParse[T UserID | RequestID | TenantID](id string) T {
	uid, err := uuid.Parse(id)
	if err != nil {
		panic(err)
	}
	return T(uid)
}

type UserID uuid.UUID

var UserIDNil = UserID(uuid.Nil)

func NewUserID() UserID {
	return UserID(uuid.New())
}

func (c UserID) MarshalText() (text []byte, err error) {
	return uuid.UUID(c).MarshalText()
}

func (c *UserID) UnmarshalText(text []byte) error {
	return (*uuid.UUID)(c).UnmarshalText(text)
}

func (c UserID) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c *UserID) Scan(src any) error {
	return (*uuid.UUID)(c).Scan(src)
}

func (c UserID) Validate() error {
	if c.IsZero() {
		return errors.New("zero UserID")
	}
	return nil
}

func (c UserID) Matches(x interface{}) bool {
	return c == x
}

// String describes what the matcher matches.
func (c UserID) String() string {
	return uuid.UUID(c).String()
}

func (c UserID) IsZero() bool {
	return c == UserIDNil
}

type RequestID uuid.UUID

var RequestIDNil = RequestID(uuid.Nil)

func NewRequestID() RequestID {
	return RequestID(uuid.New())
}

func (c RequestID) MarshalText() (text []byte, err error) {
	return uuid.UUID(c).MarshalText()
}

func (c *RequestID) UnmarshalText(text []byte) error {
	return (*uuid.UUID)(c).UnmarshalText(text)
}

func (c RequestID) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c *RequestID) Scan(src any) error {
	return (*uuid.UUID)(c).Scan(src)
}

func (c RequestID) Validate() error {
	if c.IsZero() {
		return errors.New("zero RequestID")
	}
	return nil
}

func (c RequestID) Matches(x interface{}) bool {
	return c == x
}

// String describes what the matcher matches.
func (c RequestID) String() string {
	return uuid.UUID(c).String()
}

func (c RequestID) IsZero() bool {
	return c == RequestIDNil
}

type TenantID uuid.UUID

var TenantIDNil = TenantID(uuid.Nil)

func NewTenantID() TenantID {
	return TenantID(uuid.New())
}

func (c TenantID) MarshalText() (text []byte, err error) {
	return uuid.UUID(c).MarshalText()
}

func (c *TenantID) UnmarshalText(text []byte) error {
	return (*uuid.UUID)(c).UnmarshalText(text)
}

func (c TenantID) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c *TenantID) Scan(src any) error {
	return (*uuid.UUID)(c).Scan(src)
}

func (c TenantID) Validate() error {
	if c.IsZero() {
		return errors.New("zero TenantID")
	}
	return nil
}

func (c TenantID) Matches(x interface{}) bool {
	return c == x
}

// String describes what the matcher matches.
func (c TenantID) String() string {
	return uuid.UUID(c).String()
}

func (c TenantID) IsZero() bool {
	return c == TenantIDNil
}
