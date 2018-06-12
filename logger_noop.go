package amqpx

type noopLogger struct{}

func (l noopLogger) Debug(args ...interface{}) {}
func (l noopLogger) Info(args ...interface{})  {}
func (l noopLogger) Warn(args ...interface{})  {}
func (l noopLogger) Error(args ...interface{}) {}
func (l noopLogger) Fatal(args ...interface{}) {}
