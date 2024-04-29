package framework

const (
	Local = "local"
	Dev   = "dev"
	Prod  = "prod"

	// File uploaded file from file upload middleware
	File = "@uploaded_file"

	// Rate Limit
	RateLimit = "RateLimit"

	// Для авторизации пользователя
	UserId = "user_id"
	Claims = "claims"
	// Заголовок с токеном
	Authorization = "Authorization"
)
