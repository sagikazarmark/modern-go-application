package templates

import "embed"

//go:embed landing.html
var files embed.FS

// Files returns a filesystem with static files.
func Files() embed.FS {
	return files
}
