package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// DecodeJSON представляет собой универсальный декодер со встроенной системой безопасности. Скопировано с
// https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("тело содержит плохо сформированный JSON (символ %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("тело содержит плохо сформированный JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("данные не соответствует структуре в которую записываются %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("данные не соответствуют структуре в которую записываются (символ %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("тело не должно быть пустым")

		case strings.HasPrefix(err.Error(), "json: неизвестное поле "):
			fieldName := strings.TrimPrefix(err.Error(), "json: неизвестное поле ")
			return fmt.Errorf("тело содержит неизвестный ключ %s", fieldName)

		case err.Error() == "http: тело запроса слишком велико":
			return fmt.Errorf("тело не должно быть больше, чем %d ,байт", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("тело должно содержать только одно значение JSON")
	}

	return nil
}
