package data

type ImplementationFields struct {
	BuilderType, ImplType, OutputType string
}

type InterfaceFields struct {
	OutputType  string
	BuilderType string
	CLIFunc     string
	CLITakesArg bool
}
