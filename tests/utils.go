package test

type LoggerMock struct {
}

func (l *LoggerMock) Info(args ...interface{}) {
}

func (l *LoggerMock) Debug(args ...interface{}) {
}

func (l *LoggerMock) Error(args ...interface{}) {
}
