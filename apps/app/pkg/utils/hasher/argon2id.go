// Пакет argon2id предоставляет удобную оболочку для реализации Go golang.org/x/crypto/argon2
// , упрощающую безопасное хэширование и проверку паролей с использованием Argon2.
//
// Он обеспечивает использование варианта алгоритма Argon2id и криптографически защищенных
// случайных солей.
package hasher

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"runtime"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	// ErrInvalidHash возвращается в ComparePasswordAndHash если предоставленный
	// хэш не соответствует ожидаемому формату.
	ErrInvalidHash = errors.New("argon2id: хэш имеет неправильный формат")

	// ErrIncompatibleVariant is returned by ComparePasswordAndHash if the
	// provided hash was created using a unsupported variant of Argon2.
	// Currently only argon2id is supported by this package.
	ErrIncompatibleVariant = errors.New("argon2id: несовместимый вариант argon2")

	// ErrIncompatibleVersion is returned by ComparePasswordAndHash if the
	// provided hash was created using a different version of Argon2.
	ErrIncompatibleVersion = errors.New("argon2id: несовместимая версия argon2")
)

// DefaultParams предоставляет несколько разумных параметров по умолчанию для хэширования паролей.
//
// Следует рекомендациям, приведенным в Argon2 RFC:
// "Вариант Argon2id с t=1 и максимальной доступной памятью рекомендуется в качестве
// настройки по умолчанию для всех сред. Этот параметр защищен от атак по сторонним каналам
// // и максимизирует затраты на использование специализированного оборудования bruteforce.""
//
// Параметры по умолчанию, как правило, следует использовать только для целей разработки/тестирования
// . Пользовательские параметры следует устанавливать для производственных приложений в зависимости от
// доступная память/ресурсы процессора и бизнес-требования.
var DefaultParams = &Params{
	Memory:      64 * 1024,
	Iterations:  1,
	Parallelism: uint8(runtime.NumCPU()),
	SaltLength:  16,
	KeyLength:   32,
}

// Params описывает входные параметры, используемые алгоритмом Argon2id. То
// Параметры памяти и итераций управляют вычислительными затратами на хэширование
// пароля. Чем выше эти показатели, тем больше стоимость генерации
// хэша и тем дольше время выполнения. Из этого также следует, что чем больше стоимость
// будет для любого злоумышленника, пытающегося угадать пароль. Если код выполняется
// на машине с несколькими ядрами, то вы можете уменьшить время выполнения без
// снижения стоимости за счет увеличения параметра параллелизма. Это управляет временем выполнения
// количество потоков, на которые распределена работа. Важное примечание: Изменение
// значения параметра Parallelism изменяет вывод хэша.
//
// Для получения рекомендаций и общего описания процесса выбора соответствующих параметров смотрите
// https://tools.ietf.org/html/draft-irtf-cfrg-argon2-04#section-4
type Params struct {
	// The amount of memory used by the algorithm (in kibibytes).
	Memory uint32

	// The number of iterations over the memory.
	Iterations uint32

	// The number of threads (or lanes) used by the algorithm.
	// Recommended value is between 1 and runtime.NumCPU().
	Parallelism uint8

	// Length of the random salt. 16 bytes is recommended for password hashing.
	SaltLength uint32

	// Length of the generated key. 16 bytes or more is recommended.
	KeyLength uint32
}

// // Create Hash возвращает хэш Argon2id обычного текстового пароля с использованием
// предоставленных параметров алгоритма. Возвращаемый хэш соответствует формату, используемому
// реализацией Argon2 reference C, и содержит Argon2id d в кодировке base64.
// производный ключ с префиксом соли и параметров. Это выглядит следующим образом:
//
//	$argon2id$v=19$m=65536,t=3,p=2$c29tZXNhbHQ$RdescudvJCsgt3ub+b+dWRWJTmaaJObG
func CreateHash(password string, params *Params) (hash string, err error) {
	salt, err := generateRandomBytes(params.SaltLength)
	if err != nil {
		return "", err
	}

	key := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)

	hash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, params.Memory, params.Iterations, params.Parallelism, b64Salt, b64Key)
	return hash, nil
}

// Сравнение пароля и хэша выполняет сравнение в режиме постоянного времени между
// паролем в виде обычного текста и хэшем Argon2id, используя параметры и соль
// , содержащиеся в хэше. Возвращает значение true, если они совпадают, в противном случае возвращает
// false.
func ComparePasswordAndHash(password, hash string) (match bool, err error) {
	match, _, err = CheckHash(password, hash)
	return match, err
}

// CheckHash похож на ComparePasswordAndHash, за исключением того, что он также возвращает параметры, с которыми был создан хэш
// Это может быть полезно, если вы хотите со временем обновлять параметры хэша .
func CheckHash(password, hash string) (match bool, params *Params, err error) {
	params, salt, key, err := DecodeHash(hash)
	if err != nil {
		return false, nil, err
	}

	otherKey := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	keyLen := int32(len(key))
	otherKeyLen := int32(len(otherKey))

	if subtle.ConstantTimeEq(keyLen, otherKeyLen) == 0 {
		return false, params, nil
	}
	if subtle.ConstantTimeCompare(key, otherKey) == 1 {
		return true, params, nil
	}
	return false, params, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// DecodeHash ожидает хэш, созданный из этого пакета, и анализирует его, чтобы вернуть параметры, использованные для
// его создания, а также соль и ключ (хэш пароля).
func DecodeHash(hash string) (params *Params, salt, key []byte, err error) {
	vals := strings.Split(hash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	if vals[1] != "argon2id" {
		return nil, nil, nil, ErrIncompatibleVariant
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	params = &Params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.SaltLength = uint32(len(salt))

	key, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	params.KeyLength = uint32(len(key))

	return params, salt, key, nil
}
