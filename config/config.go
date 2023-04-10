package config

var (
	OutputFormat    string // json or yaml
	AppOutputFormat string // should match your app's output format
	AltFieldFormat  string // used for alternative field formatting: snakeCase, kebabCase, camelCase, pascalCase, upperCase, lowerCase
)
