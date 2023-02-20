package example

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/FrangipaneTeam/terraform-templates/pkg/file"
	"github.com/rs/zerolog/log"
)

func genExampleFromTest(str, tfType string) (string, string) {
	startDoc := regexp.MustCompile("const testAcc.*`")
	endDoc := regexp.MustCompile("^`$")
	tfNameRe := regexp.MustCompile(`^(resource|data)\s+"(\S+)"\s+.*`)
	definition := regexp.MustCompile(`(\S+\s+=\s+)cloudavenue`)

	doc := ""
	startFound := false
	endFound := false
	badTfType := false
	tfName := ""

	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		line := scanner.Text()

		// check for doc start
		if startDoc.MatchString(line) {
			log.Info().Msgf("found start doc: %s", line)
			startFound = true
			if endFound {
				doc += "\n"
				endFound = false
			}
			continue
		}

		// check for tf name
		if tfNameRe.MatchString(line) {
			foundTfType := tfNameRe.FindStringSubmatch(line)[1]
			if tfType != foundTfType {
				log.Info().Msgf("tf type %s not match %s", foundTfType, tfType)
				badTfType = true
				continue
			} else {
				badTfType = false
				tfName = tfNameRe.FindStringSubmatch(line)[2]
				log.Info().Msgf("found tf name: %s", tfName)
			}
		}

		// check for ref in definition
		if definition.MatchString(line) {
			log.Info().Msgf("found definition: %s", line)
			doc += "\t" + definition.FindStringSubmatch(line)[1] + "\"your_value\"\n"
			continue
		}
		// check for doc end
		if endDoc.MatchString(line) {
			log.Info().Msgf("found end doc: %s", line)
			startFound = false
			endFound = true
			continue
		}
		if startFound && !badTfType {
			doc += line + "\n"
		}
	}

	return doc, tfName
}

func CreateExampleFile(fileName, exampleDir string) error {
	// check if filename contain datasource or resource
	if !strings.Contains(fileName, "datasource") && !strings.Contains(fileName, "resource") {
		return errors.New("filename must contain datasource or resource")
	}

	tfType := "resource"
	tfTest := "resource"
	if strings.Contains(fileName, "datasource") {
		tfType = "data-source"
		tfTest = "data"
	}

	f, err := file.FileToString(fileName)
	if err != nil {
		return err
	}

	doc, tfName := genExampleFromTest(f, tfTest)
	if doc == "" {
		return errors.New("doc is empty")
	}

	log.Info().Msgf("doc: %s", doc)

	doc = strings.TrimSpace(doc)

	exampleDirPath := exampleDir + "/" + tfType + "s/" + tfName
	errMkdir := os.MkdirAll(exampleDirPath, 0o755)
	if errMkdir != nil {
		return errMkdir
	}

	errWrite := os.WriteFile(exampleDirPath+"/"+tfType+".tf", []byte(doc), 0o644)
	if errWrite != nil {
		return errWrite
	}
	return nil
}
