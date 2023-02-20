# tf-doc-extractor
`tf-doc-extractor` is a command-line tool that generates Terraform examples from acceptance tests and import functions.

## Installation

To install `tf-doc-extractor`, use the following command:

```bash
go install github.com/FrangipaneTeam/tf-doc-extractor@latest
```

# Usage
## Generating import example
Add the go:generate directive above your Terraform import function.
* `filename` is the current file
* `example-dir` is the location of the Terraform example directory

```go
//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../../examples -resource
 func (r *orgUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
        resource.ImportStatePassthroughID(ctx, path.Root("user_name"), req, resp)
 }
```

Running go generate will create a file `import.sh` in `../../../examples/resources/cloudavenue_org_user` with this content :

```terraform
# use the user_name to import the resource
terraform import cloudavenue_org_user.example user_name
```

## Generating examples from acceptance tests
Add the go:generate directive above your Terraform example in the test file. For example, with a test file named `internal/tests/public_ip_datasource_test.go` :
```go
//go:generate go run github.com/FrangipaneTeam/tf-doc-extractor@latest -filename $GOFILE -example-dir ../../examples -test
 const testAccPublicIPDataSourceConfig = `
 data "cloudavenue_public_ip" "test" {}
 `
```
Running go generate will create a file `resource.tf` or `data-source.tf` in `../../examples/data-sources/cloudavenue_public_ip`.

# Contributing
Pull requests are welcome! If you find a bug or would like to request a new feature, please open an issue.

Before submitting a pull request, please ensure that your changes are properly tested and that the documentation has been updated.
