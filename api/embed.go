package api

import "embed"

//go:embed */*.proto
var EmbedProto embed.FS
