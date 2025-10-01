package logs

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func GetFunctionName() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "unknown_function"
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown_function"
	}
	name := fn.Name()
	return name[strings.LastIndex(name, ".")+1:]
}

var logger zerolog.Logger

func init() {
	cw := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		PartsOrder: []string{zerolog.LevelFieldName, zerolog.TimestampFieldName, zerolog.MessageFieldName},
	}

	cw.FormatLevel = func(i interface{}) string {
		var l string
		if ll, ok := i.(string); ok {
			l = strings.ToUpper(ll)
		} else {
			l = "????"
		}
		switch l {
		case "INFO":
			return fmt.Sprintf("\x1b[32m| %s |\x1b[0m", l)
		case "WARN":
			return fmt.Sprintf("\x1b[33m| %s |\x1b[0m", l)
		case "ERROR":
			return fmt.Sprintf("\x1b[31m| %s |\x1b[0m", l)
		default:
			return fmt.Sprintf("| %s |", l)
		}
	}

	logger = zerolog.New(cw).With().Timestamp().Str("service", "activie-core").Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func Error(method string, message string, err error) {
	logger.Error().
		Str("function", GetFunctionName()).
		Str("method", method).
		Err(err).
		Msg(message)
}

func Info(message string, extra ...any) {
	event := logger.Info().
		Str("function", GetFunctionName())

	for i := 0; i < len(extra); i += 2 {
		if i+1 < len(extra) {
			key, ok := extra[i].(string)
			if !ok {
				continue
			}
			event = event.Interface(key, extra[i+1])
		}
	}

	event.Msg(message)
}

func Debug(method string, message string, key string) {
	logger.Debug().
		Str("function", GetFunctionName()).
		Str("method", method).
		Str("key", key).
		Msg(message)
}

func Warn(method string, message string, extra any) {
	logger.Warn().
		Str("function", GetFunctionName()).
		Str("method", method).
		Interface("extra", extra).
		Msg(message)
}
