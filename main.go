package main

import (
	_ "embed"
	"flag"
	"strings"

	"github.com/FrangipaneTeam/terraform-templates/pkg/file"
	"github.com/rs/zerolog/log"

	"github.com/tcnksm/go-latest"

	"github.com/FrangipaneTeam/tf-doc-extractor/internal/example"
)

//go:embed version.txt
var version string

func main() {
	fileName := flag.String("filename", "", "filename")
	exampleDir := flag.String("example-dir", "", "example directory")
	fromTest := flag.Bool("test", false, "from test")
	fromResource := flag.Bool("resource", false, "from resource")
	flag.Parse()
	// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	version = strings.Trim(version, "\n")
	log.Info().Msgf("tf-doc-extractor version %s", version)

	if *fileName == "" {
		log.Fatal().Msg("filename is required")
	}

	if !file.IsFileExists(*fileName) {
		log.Fatal().Msgf("file %s not found", *fileName)
	}

	if !*fromTest && !*fromResource {
		log.Fatal().Msg("test or resource is required")
	}

	if *fromTest && *fromResource {
		log.Fatal().Msg("test and resource are exclusive")
	}

	// check version
	githubTag := &latest.GithubTag{
		Owner:      "FrangipaneTeam",
		Repository: "tf-doc-extractor",
	}

	res, err := latest.Check(githubTag, version)
	if err == nil {
		if res.Outdated {
			log.Warn().Msgf("new version availaible : %s", res.Current)
		}
	} else {
		log.Warn().Err(err).Msg("failed to check version")
	}

	log.Info().Msgf("using file %s", *fileName)
	if *fromTest {
		err := example.CreateExampleFile(*fileName, *exampleDir)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create example file")
		}
	} else if *fromResource {
		err := example.CreateImportExampleFile(*fileName, *exampleDir)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create example file")
		}
	}
}
