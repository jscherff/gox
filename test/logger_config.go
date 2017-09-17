type Logger struct {
	*MLogger
	LogFile string
	LoggerFlags int
	Stdout bool
	Stderr bool
}

type Loggers map[string]*Logger
