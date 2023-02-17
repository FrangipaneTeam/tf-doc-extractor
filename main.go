package main

import (
	"flag"

	"github.com/FrangipaneTeam/terraform-templates/pkg/file"
	"github.com/rs/zerolog/log"

	"github.com/FrangipaneTeam/tf-doc-extractor/internal/example"
)

func main() {
	fileName := flag.String("filename", "", "filename")
	exampleDir := flag.String("example-dir", "", "example directory")
	fromTest := flag.Bool("test", false, "from test")
	fromResource := flag.Bool("resource", false, "from resource")
	flag.Parse()
	// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

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
