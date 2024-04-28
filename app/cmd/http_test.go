package cmd

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"testing"
)

func TestValidateApp(t *testing.T) {
	err := fx.ValidateApp(inject())
	require.NoError(t, err)
}
