package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Запускает cmd команду файл приведен для теста",
	Long:  `Запуск всех воркеров для работы (пока чисто для теста)`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Это тестовая команда")
	},
}
