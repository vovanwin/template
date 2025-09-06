package templgen

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "templ",
	Short: "Generate templ components",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := exec.Command("templ", "generate")
		out, err := c.CombinedOutput()
		if err != nil {
			return fmt.Errorf("templ generate failed: %w, output: %s", err, string(out))
		}
		fmt.Print(string(out))
		return nil
	},
}
