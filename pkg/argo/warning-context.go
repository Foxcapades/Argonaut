package argo

type WarningContext struct {
	warnings []string
}

func (w *WarningContext) AppendWarning(warning string) {
	w.warnings = append(w.warnings, warning)
}

func (w *WarningContext) GetWarnings() []string {
	return w.warnings
}
