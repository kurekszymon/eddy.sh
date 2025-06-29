package exit_codes

const (
	SUCCESS                          = 0
	SOMETHING_WENT_WRONG             = 2
	NO_CONFIG                        = 3
	WRONG_CONFIG                     = 4
	BREW_SPECIFIED_BUT_NOT_INSTALLED = 5
	NO_GIT                           = 6
	SSH_KEYS_DENIED                  = 7
	CUSTOM_SCRIPT_EXIT               = 8
	// cli
	UNKNOWN_COMMAND                = 11
	CLI_INSTALL_TOOL_NOT_SPECIFIED = 12
	TOOL_NOT_INSTALLED             = 13
	UNKNOWN_TOOL                   = 14
)
