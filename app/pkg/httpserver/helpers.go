package httpserver

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
)

// PrintAllRegisteredRoutes печатает все зарегистрированные маршруты с маршрутизатора Chi.
// exceptions - исключает роуты.
func PrintAllRegisteredRoutes(r *chi.Mux, exceptions ...string) {
	exceptions = append(exceptions, "/swagger", "/metrics", "/debug")

	walkFunc := func(method string, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		for _, val := range exceptions {
			if strings.HasPrefix(path, val) {
				return nil
			}
		}

		switch method {
		case "GET":
			fmt.Printf("%s", color.GreenString(fmt.Sprintf("%-8s", method)))
		case "POST", "PUT", "PATCH":
			fmt.Printf("%s", color.YellowString(fmt.Sprintf("%-8s", method)))
		case "DELETE":
			fmt.Printf("%s", color.RedString(fmt.Sprintf("%-8s", method)))
		default:
			fmt.Printf("%s", color.WhiteString(fmt.Sprintf("%-8s", method)))
		}

		fmt.Printf("%s", strPad(path, 25, "-", "RIGHT"))
		fmt.Printf("%s\n", strPad(getHandler(getModuleName(), handler), 90, "-", "LEFT"))

		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		log.Printf("Error walking routes: %v", err)
	}
}

// strPad возвращает входную строку, заполненную с левой, правой или обеих сторон с использованием pad Type, с указанной длиной заполнения padLength.
func strPad(input string, padLength int, padString string, padType string) string {
	var output string
	inputLength := len(input)
	padStringLength := len(padString)

	if inputLength >= padLength {
		return input
	}

	repeat := int(math.Ceil(float64(1) + (float64(padLength-padStringLength))/float64(padStringLength)))

	switch padType {
	case "RIGHT":
		output = input + strings.Repeat(padString, repeat)
		output = output[:padLength]
	case "LEFT":
		output = strings.Repeat(padString, repeat) + input
		output = output[len(output)-padLength:]
	case "BOTH":
		length := (padLength - inputLength) / 2
		repeat = int(math.Ceil(float64(length) / float64(padStringLength)))
		output = strings.Repeat(padString, repeat)[:length] + input + strings.Repeat(padString, repeat)[:length]
	}

	return output
}

// getHandler возвращает имя обработчика для заданного проекта.
func getHandler(projectName string, handler http.Handler) string {
	funcName := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	base := filepath.Base(funcName)

	nameSplit := strings.Split(funcName, "")
	names := nameSplit[len(projectName):]
	path := strings.Join(names, "")

	pathSplit := strings.Split(path, "/")
	path = strings.Join(pathSplit[:len(pathSplit)-1], "/")

	sFull := strings.Split(base, ".")
	s := sFull[len(sFull)-1:]

	s = strings.Split(s[0], "")
	if len(s) <= 4 && len(sFull) >= 3 {
		s = sFull[len(sFull)-3 : len(sFull)-2]
		return "@" + color.BlackString(strings.Join(s, ""))
	}
	s = s[:len(s)-3]
	funcName = strings.Join(s, "")

	return path + "@" + color.BlackString(funcName)
}

// getModuleName возвращает имя модуля
func getModuleName() string {
	moduleName := "app"
	return moduleName
}
