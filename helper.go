package Shipyard

import "embed"

//go:embed all:web/build
var WebContent embed.FS
