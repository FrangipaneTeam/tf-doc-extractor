package example

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/FrangipaneTeam/terraform-templates/pkg/file"

	"github.com/FrangipaneTeam/tf-doc-extractor/internal/format"
	"github.com/FrangipaneTeam/tf-doc-extractor/internal/logger"
)

func genExampleFromTest(str, tfType string) (string, string) {
	startDoc := regexp.MustCompile("const testAcc.*`")
	endDoc := regexp.MustCompile("^`$")
	tfNameRe := regexp.MustCompile(`^(resource|data)\s+"(\S+)"\s+.*`)
	definition := regexp.MustCompile(`\S+\s+=\s+cloudavenue`)

	doc := ""
	startFound := false
	endFound := false
	tfName := ""

	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		line := scanner.Text()

		// check for doc start
		if startDoc.MatchString(line) {
			logger.Logger.Info().Msgf("found start doc: %s", line)
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
			if tfType == foundTfType {
				tfName = tfNameRe.FindStringSubmatch(line)[2]
				logger.Logger.Info().Msgf("found tf name: %s", tfName)
			}
		}

		// check for ref in definition
		if definition.MatchString(line) {
			logger.Logger.Debug().Msgf("found definition: %s", line)
			doc += "\t" + line + "\n"
			continue
		}
		// check for doc end
		if endDoc.MatchString(line) {
			logger.Logger.Info().Msgf("found end doc: %s", line)
			break
		}
		if startFound {
			doc += line + "\n"
		}
	}

	f, err := hclwrite.ParseConfig([]byte(doc), "", hcl.Pos{Line: 1, Column: 1})
	if err.HasErrors() {
		logger.Logger.Error().Msg("error parsing doc with hclsyntax")
		return "", ""
	}

	format.Body(f.Body(), nil)

	return string(f.Bytes()), tfName
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

	f, err := file.ToString(fileName)
	if err != nil {
		return err
	}

	doc, tfName := genExampleFromTest(f, tfTest)
	if doc == "" {
		return errors.New("doc is empty")
	}

	logger.Logger.Info().Msgf("doc: %s", doc)

	doc = strings.TrimSpace(doc)

	exampleDirPath := exampleDir + "/" + tfType + "s/" + tfName
	if err := os.MkdirAll(exampleDirPath, 0o755); err != nil {
		return err
	}

	return os.WriteFile(exampleDirPath+"/"+tfType+".tf", []byte(doc), 0o644) //nolint:gosec
}
