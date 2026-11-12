//go:build (darwin && !ios) || (linux && !android) || windows

package res

import (
	"golang.org/x/mobile/asset"
)

// Open opens a named asset.
func Open(name string) (File, error) {
	return asset.Open(name)
}
