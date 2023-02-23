package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

// SetupZeroLog setup zerolog.
func SetupZeroLog(version string, debug bool) zerolog.Logger { // Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	output.FormatLevel = func(i interface{}) string {
		return fmt.Sprintf("tf-doc-extractor %s -  %-6s:", version, strings.ToUpper(i.(string)))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	output.FormatTimestamp = func(i interface{}) string {
		return ""
	}
	return zerolog.New(output).With().Logger()
}
