package globals

import (
	"os"
	"strings"
)

const CONFIG_FILE = "config.yaml"
const CONFIG_URL = "https://raw.githubusercontent.com/kurekszymon/eddy.sh/refs/heads/main/config.yaml"

var DebugEnabled = strings.EqualFold(os.Getenv("EDDY_DEBUG"), "1")
