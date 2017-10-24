type Logger struct {
	*MLogger
	LogFile string
	LogFlags []string
	Stdout bool
	Stderr bool
}

type Loggers map[string]*Logger
