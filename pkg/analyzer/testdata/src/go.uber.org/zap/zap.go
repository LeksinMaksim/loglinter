package zap

type Logger struct{}

func NewProduction() (*Logger, error) {
	return &Logger{}, nil
}

func (l *Logger) Info(msg string, fields ...any)  {}
func (l *Logger) Error(msg string, fields ...any) {}
func (l *Logger) Warn(msg string, fields ...any)  {}
func (l *Logger) Debug(msg string, fields ...any) {}
func (l *Logger) Fatal(msg string, fields ...any) {}

type SugaredLogger struct{}

func (l *Logger) Sugar() *SugaredLogger { return &SugaredLogger{} }

func (s *SugaredLogger) Info(args ...any)  {}
func (s *SugaredLogger) Error(args ...any) {}
func (s *SugaredLogger) Warn(args ...any)  {}
func (s *SugaredLogger) Debug(args ...any) {}
func (s *SugaredLogger) Fatal(args ...any) {}
