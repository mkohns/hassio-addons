package main

// get the swagger here: https://bbernhard.github.io/signal-cli-rest-api/src/docs/swagger.json
// convert the swagger 2.0 to OAS here: https://converter.swagger.io/#/Converter/convertByContent
// convert the json to yaml here: https://editor.swagger.io/

//go:generate oapi-codegen -generate "types, client" -package signal -o signal/signal_client_gen.go oas_signal.yaml
