package store

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/versioned-migration   --target ./gen --template ./templates ./schema
