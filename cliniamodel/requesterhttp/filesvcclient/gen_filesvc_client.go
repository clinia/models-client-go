//go:build generate
// +build generate

package filesvcclient

//go:generate go tool oapi-codegen --config ../../../oapi/filesvc/config.yaml https://raw.githubusercontent.com/clinia/model-foundry/refs/tags/v0.1.0/packages/baseten/filesvc/filesvc/api/openapi.yaml?token=$GH_URL_TOKEN
