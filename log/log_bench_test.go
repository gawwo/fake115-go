package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"testing"
)

type Syncer struct {
	err    error
	called bool
}

func (s *Syncer) SetError(err error) {
	s.err = err
}

// Sync records that it was called, then returns the user-supplied error (if
// any).
func (s *Syncer) Sync() error {
	s.called = true
	return s.err
}

// Called reports whether the Sync method was called.
func (s *Syncer) Called() bool {
	return s.called
}

type Discarder struct{ Syncer }

func (d *Discarder) Write(b []byte) (int, error) {
	return ioutil.Discard.Write(b)
}

func withBenchedLogger(b *testing.B, f func(logger *zap.Logger)) {
	logger := InitLogger("test", false)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f(logger)
		}
	})
}

func withBenchedFileLogger(b *testing.B, f func(logger *zap.Logger)) {
	logger := FileLogger()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f(logger)
		}
	})
}

func withBenchedPureLogger(b *testing.B, f func(*zap.Logger)) {
	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionConfig().EncoderConfig),
			&Discarder{},
			zap.DebugLevel,
		))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			f(logger)
		}
	})
}

func BenchmarkInitLogger(b *testing.B) {
	withBenchedLogger(b, func(logger *zap.Logger) {
		logger.Info("info")
	})
}

func BenchmarkFileLogger(b *testing.B) {
	withBenchedFileLogger(b, func(logger *zap.Logger) {
		logger.Info("info")
	})
}

func BenchmarkInitLogger10Field(b *testing.B) {
	withBenchedLogger(b, func(logger *zap.Logger) {

		logger.Info("info",
			zap.Int("one", 1),
			zap.Int("two", 2),
			zap.Int("three", 3),
			zap.Int("four", 4),
			zap.Int("five", 5),
			zap.Int("six", 6),
			zap.Int("seven", 7),
			zap.Int("eight", 8),
			zap.Int("nine", 9),
			zap.Int("ten", 10))
	})
}

func BenchmarkFileLogger10Field(b *testing.B) {
	withBenchedFileLogger(b, func(logger *zap.Logger) {

		logger.Info("info",
			zap.Int("one", 1),
			zap.Int("two", 2),
			zap.Int("three", 3),
			zap.Int("four", 4),
			zap.Int("five", 5),
			zap.Int("six", 6),
			zap.Int("seven", 7),
			zap.Int("eight", 8),
			zap.Int("nine", 9),
			zap.Int("ten", 10))
	})
}

func BenchmarkPureLogger10Field(b *testing.B) {
	withBenchedPureLogger(b, func(logger *zap.Logger) {

		logger.Info("info",
			zap.Int("one", 1),
			zap.Int("two", 2),
			zap.Int("three", 3),
			zap.Int("four", 4),
			zap.Int("five", 5),
			zap.Int("six", 6),
			zap.Int("seven", 7),
			zap.Int("eight", 8),
			zap.Int("nine", 9),
			zap.Int("ten", 10))
	})
}
