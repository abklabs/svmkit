package runner

import (
	_ "embed"
)

//go:embed assets/lib.bash
var LibBash string

//go:embed assets/run.sh
var RunScript string
