package main

import (
	_ "embed"
	"flag"
	"strings"

	"github.com/FrangipaneTeam/terraform-templates/pkg/file"
	latest "github.com/tcnksm/go-latest"

	"github.com/FrangipaneTeam/tf-doc-extractor/internal/example"
	"github.com/FrangipaneTeam/tf-doc-extractor/internal/logger"
)

//go:embed version.txt
var version string

func main() {
	fileName := flag.String("filename", "", "filename")
	exampleDir := flag.String("example-dir", "", "example directory")
	fromTest := flag.Bool("test", false, "from test")
	fromResource := flag.Bool("resource", false, "from resource")
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()

	version = strings.Trim(version, "\n")

	logger.Logger = logger.SetupZeroLog(version, *debug)

	logger.Logger.Info().Msgf("tf-doc-extractor version %s", version)

	if *fileName == "" {
		logger.Logger.Fatal().Msg("filename is required")
	}

	if !file.IsFileExists(*fileName) {
		logger.Logger.Fatal().Msgf("file %s not found", *fileName)
	}

	if !*fromTest && !*fromResource {
		logger.Logger.Fatal().Msg("test or resource is required")
	}

	if *fromTest && *fromResource {
		logger.Logger.Fatal().Msg("test and resource are exclusive")
	}

	// check version
	githubTag := &latest.GithubTag{
		Owner:      "FrangipaneTeam",
		Repository: "tf-doc-extractor",
	}

	res, err := latest.Check(githubTag, version)
	if err == nil {
		if res.Outdated {
			logger.Logger.Warn().Msgf("new version availaible : %s", res.Current)
		}
	} else {
		logger.Logger.Warn().Err(err).Msg("failed to check version")
	}

	logger.Logger.Info().Msgf("using file %s", *fileName)
	if *fromTest {
		errCreateExample := example.CreateExampleFile(*fileName, *exampleDir)
		if errCreateExample != nil {
			logger.Logger.Fatal().Err(errCreateExample).Msg("failed to create example file")
		}
	} else if *fromResource {
		errCreateImport := example.CreateImportExampleFile(*fileName, *exampleDir)
		if errCreateImport != nil {
			logger.Logger.Fatal().Err(errCreateImport).Msg("failed to create example file")
		}
	}
}
