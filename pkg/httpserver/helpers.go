package httpserver

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
	"golang.org/x/mod/modfile"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// PrintAllRegisteredRoutes печатает все зарегистрированные маршруты с маршрутизатора Chi.
// exceptions - исключает роуты
func printAllRegisteredRoutes(r *chi.Mux, exceptions ...string) {
	exceptions = append(exceptions, "/swagger")
	exceptions = append(exceptions, "/metrics")

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

		//fmt.Printf("%-25s %60s\n", path, getHandler(getModName(), handler))
		fmt.Printf("%s", strPad(path, 25, "-", "RIGHT"))
		fmt.Printf("%s\n", strPad(getHandler(getModName(), handler), 90, "-", "LEFT"))

		return nil
	}
	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Print(err)
	}

}

// StrPad возвращает входную строку, заполненную с левой, правой или обеих сторон с использованием pad Type, с указанной длиной заполнения padLength.
//
// Example:
// input := "Codes";
// StrPad(input, 10, " ", "RIGHT")        // produces "Codes     "
// StrPad(input, 10, "-=", "LEFT")        // produces "=-=-=Codes"
// StrPad(input, 10, "_", "BOTH")         // produces "__Codes___"
// StrPad(input, 6, "___", "RIGHT")       // produces "Codes_"
// StrPad(input, 3, "*", "RIGHT")         // produces "Codes"
// taken from // https://gist.github.com/asessa/3aaec43d93044fc42b7c6d5f728cb039
func strPad(input string, padLength int, padString string, padType string) string {
	var output string

	inputLength := len(input)
	padStringLength := len(padString)

	if inputLength >= padLength {
		return input
	}

	repeat := math.Ceil(float64(1) + (float64(padLength-padStringLength))/float64(padStringLength))

	switch padType {
	case "RIGHT":
		output = input + strings.Repeat(padString, int(repeat))
		output = output[:padLength]
	case "LEFT":
		output = strings.Repeat(padString, int(repeat)) + input
		output = output[len(output)-padLength:]
	case "BOTH":
		length := (float64(padLength - inputLength)) / float64(2)
		repeat = math.Ceil(length / float64(padStringLength))
		output = strings.Repeat(padString, int(repeat))[:int(math.Floor(float64(length)))] + input + strings.Repeat(padString, int(repeat))[:int(math.Ceil(float64(length)))]
	}

	return output
}

func getHandler(projectName string, handler http.Handler) (funcName string) {
	// https://github.com/go-chi/chi/issues/424
	funcName = runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
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

// для получения имени модуля.  Адаптировано из https://stackoverflow.com/a/63393712/1033134
func getModName() string {
	goModBytes, err := os.ReadFile("go.mod")
	if err != nil {
		os.Exit(0)
	}
	return modfile.ModulePath(goModBytes)
}
