package example

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/FrangipaneTeam/terraform-templates/pkg/file"

	"github.com/FrangipaneTeam/tf-doc-extractor/internal/logger"
)

func genExampleFromResource(str string) (string, string) {
	importNameRe := regexp.MustCompile(`resource\.ImportStatePassthroughID\(ctx, path\.Root\("(\S+)"\), req, resp\)`)
	tfNameRe := regexp.MustCompile(`resp\.TypeName = req\.ProviderTypeName \+ "_" \+ categoryName \+ "_" \+ "(\S+)"`)

	importName := ""
	tfName := ""

	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		line := scanner.Text()

		if importNameRe.MatchString(line) {
			logger.Logger.Debug().Msgf("found import line : %s", line)
			importName = importNameRe.FindStringSubmatch(line)[1]
			break
		}

		if tfNameRe.MatchString(line) {
			tfName = tfNameRe.FindStringSubmatch(line)[1]
			logger.Logger.Info().Msgf("found tf name: %s", tfName)
			continue
		}
	}

	return importName, tfName
}

func getCategoryFromBase(resourceFile string) string {

	categoryNameRe := regexp.MustCompile(`categoryName = "(\S+)"`)
	categoryName := ""

	x := strings.Split(resourceFile, "/")
	x[len(x)-1] = "base.go"
	fileCategory := strings.Join(x, "/")

	f, err := file.ToString(fileCategory)
	if err != nil {
		return ""
	}

	scannerCategory := bufio.NewScanner(strings.NewReader(f))
	for scannerCategory.Scan() {
		line := scannerCategory.Text()

		if categoryNameRe.MatchString(line) {
			logger.Logger.Info().Msgf("found category name: %s", categoryName)
			categoryName = categoryNameRe.FindStringSubmatch(line)[1]
			break
		}
	}

	return categoryName
}

func CreateImportExampleFile(fileName, exampleDir string) error {
	f, err := file.ToString(fileName)
	if err != nil {
		return err
	}

	importName, tfName := genExampleFromResource(f)
	logger.Logger.Info().Msgf("importName: %s", importName)
	logger.Logger.Info().Msgf("tfName: %s", tfName)

	categoryName := getCategoryFromBase(fileName)
	logger.Logger.Info().Msgf("categoryName: %s", categoryName)

	if importName == "" || tfName == "" || categoryName == "" {
		return fmt.Errorf("failed to generate import example for %s", fileName)
	}

	exampleDirPath := exampleDir + "/resources/cloudavenue_" + categoryName + "_" + tfName
	errMkdir := os.MkdirAll(exampleDirPath, 0o755)
	if errMkdir != nil {
		return errMkdir
	}

	doc := fmt.Sprintf("# use the %s to import the resource\n", importName)
	doc += fmt.Sprintf("terraform import cloudavenue_%s_%s.example %s", categoryName, tfName, importName)

	errWrite := os.WriteFile(exampleDirPath+"/import.sh", []byte(doc), 0o644)
	if errWrite != nil {
		return errWrite
	}
	return nil
}
