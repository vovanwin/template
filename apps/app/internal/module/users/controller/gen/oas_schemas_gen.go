// Code generated by ogen, DO NOT EDIT.

package usersGenv1

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *ErrorStatusCode) Error() string {
	return fmt.Sprintf("code %d: %+v", s.StatusCode, s.Response)
}

// Ref: #/components/schemas/AuthToken
type AuthToken struct {
	// Токен для авторизации.
	Access string `json:"access"`
	// Токен для получения нового access токена.
	Refresh string `json:"refresh"`
}

// GetAccess returns the value of Access.
func (s *AuthToken) GetAccess() string {
	return s.Access
}

// GetRefresh returns the value of Refresh.
func (s *AuthToken) GetRefresh() string {
	return s.Refresh
}

// SetAccess sets the value of Access.
func (s *AuthToken) SetAccess(val string) {
	s.Access = val
}

// SetRefresh sets the value of Refresh.
func (s *AuthToken) SetRefresh(val string) {
	s.Refresh = val
}

type BearerAuth struct {
	Token string
}

// GetToken returns the value of Token.
func (s *BearerAuth) GetToken() string {
	return s.Token
}

// SetToken sets the value of Token.
func (s *BearerAuth) SetToken(val string) {
	s.Token = val
}

// Represents error object.
// Ref: #/components/schemas/Error
type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// GetCode returns the value of Code.
func (s *Error) GetCode() int64 {
	return s.Code
}

// GetMessage returns the value of Message.
func (s *Error) GetMessage() string {
	return s.Message
}

// SetCode sets the value of Code.
func (s *Error) SetCode(val int64) {
	s.Code = val
}

// SetMessage sets the value of Message.
func (s *Error) SetMessage(val string) {
	s.Message = val
}

// ErrorStatusCode wraps Error with StatusCode.
type ErrorStatusCode struct {
	StatusCode int
	Response   Error
}

// GetStatusCode returns the value of StatusCode.
func (s *ErrorStatusCode) GetStatusCode() int {
	return s.StatusCode
}

// GetResponse returns the value of Response.
func (s *ErrorStatusCode) GetResponse() Error {
	return s.Response
}

// SetStatusCode sets the value of StatusCode.
func (s *ErrorStatusCode) SetStatusCode(val int) {
	s.StatusCode = val
}

// SetResponse sets the value of Response.
func (s *ErrorStatusCode) SetResponse(val Error) {
	s.Response = val
}

// Ref: #/components/schemas/LoginRequest
type LoginRequest struct {
	// Логин пользователя. Может быть как email так и логином.
	Username string `json:"username"`
	// Пароль.
	Password string `json:"password"`
}

// GetUsername returns the value of Username.
func (s *LoginRequest) GetUsername() string {
	return s.Username
}

// GetPassword returns the value of Password.
func (s *LoginRequest) GetPassword() string {
	return s.Password
}

// SetUsername sets the value of Username.
func (s *LoginRequest) SetUsername(val string) {
	s.Username = val
}

// SetPassword sets the value of Password.
func (s *LoginRequest) SetPassword(val string) {
	s.Password = val
}

// NewOptUUID returns new OptUUID with value set to v.
func NewOptUUID(v uuid.UUID) OptUUID {
	return OptUUID{
		Value: v,
		Set:   true,
	}
}

// OptUUID is optional uuid.UUID.
type OptUUID struct {
	Value uuid.UUID
	Set   bool
}

// IsSet returns true if OptUUID was set.
func (o OptUUID) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptUUID) Reset() {
	var v uuid.UUID
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptUUID) SetTo(v uuid.UUID) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptUUID) Get() (v uuid.UUID, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptUUID) Or(d uuid.UUID) uuid.UUID {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// Ref: #/components/schemas/UserMe
type UserMe struct {
	// Токен для авторизации.
	ID uuid.UUID `json:"id"`
	// Email пользователя, а также его логин. Может не быть
	// почтовым адресом.
	Email string `json:"email"`
	// Роль текущего пользователя.
	Role string `json:"role"`
	// Тенант текущего пользователя.
	Tenant string `json:"tenant"`
	// Время создания пользователя.
	CreatedAt time.Time `json:"created_at"`
	// Тут хранятся все настройки пользователя для
	// фронтенда, фильтры, таймзона и тд.
	Settings string `json:"settings"`
	// Разделы меню доступные пользователю (сейчас
	// захардкожено).
	Components []string `json:"components"`
}

// GetID returns the value of ID.
func (s *UserMe) GetID() uuid.UUID {
	return s.ID
}

// GetEmail returns the value of Email.
func (s *UserMe) GetEmail() string {
	return s.Email
}

// GetRole returns the value of Role.
func (s *UserMe) GetRole() string {
	return s.Role
}

// GetTenant returns the value of Tenant.
func (s *UserMe) GetTenant() string {
	return s.Tenant
}

// GetCreatedAt returns the value of CreatedAt.
func (s *UserMe) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// GetSettings returns the value of Settings.
func (s *UserMe) GetSettings() string {
	return s.Settings
}

// GetComponents returns the value of Components.
func (s *UserMe) GetComponents() []string {
	return s.Components
}

// SetID sets the value of ID.
func (s *UserMe) SetID(val uuid.UUID) {
	s.ID = val
}

// SetEmail sets the value of Email.
func (s *UserMe) SetEmail(val string) {
	s.Email = val
}

// SetRole sets the value of Role.
func (s *UserMe) SetRole(val string) {
	s.Role = val
}

// SetTenant sets the value of Tenant.
func (s *UserMe) SetTenant(val string) {
	s.Tenant = val
}

// SetCreatedAt sets the value of CreatedAt.
func (s *UserMe) SetCreatedAt(val time.Time) {
	s.CreatedAt = val
}

// SetSettings sets the value of Settings.
func (s *UserMe) SetSettings(val string) {
	s.Settings = val
}

// SetComponents sets the value of Components.
func (s *UserMe) SetComponents(val []string) {
	s.Components = val
}
