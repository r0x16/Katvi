package model

type CommandOutput struct {
	Stdout         string
	FilteredOutput [][][]string
	ExitCode       int
}
