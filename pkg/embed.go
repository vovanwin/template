package pkg

import "embed"

//go:embed */*.swagger.json
var EmbedSwagger embed.FS
