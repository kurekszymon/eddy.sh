package types

type LogType string

const (
	LogDebug   LogType = "debug"
	LogWarning LogType = "warning"
	LogError   LogType = "error"
	LogInfo    LogType = "info"
)
