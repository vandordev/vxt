package bind

import "github.com/vandordev/vxt/source"

type embeddedAssets struct {
	Main source.Source
	Uses map[string]source.Source
}
