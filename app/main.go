package main

import (
	"template/app/cmd"
)

//	@title			API Service
//	@version		1.0
//	@description	API service Backend.
//	@termsOfService	http://swagger.io/terms/

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

// @contact.name   API Serivce
// @contact.url    https://vovanwin.ru
// @contact.email  iot@megafon.ru

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @Schemes		http https
// @host		localhost:8080
// @BasePath	/api/
func main() {
	cmd.Execute()
}
