package web

type ErrorDTO struct {
	Error string `json:"errorMessage"`
}

type SessionIdentifierDTO struct {
	ID string `json:"sessionId"`
}

type LoadPackDTO struct {
	Pack string `json:"pack"`
}

type PackDefDTO struct {
	Name  string `json:"id"`
	Title string `json:"title"`
}

func NewPackDefDTO(name string, title string) *PackDefDTO {
	return &PackDefDTO{
		Name:  name,
		Title: title,
	}
}

type EvalDTO struct {
	Expr string `json:"expression"`
	Pack string `json:"pack"`
}

type EvalResultDTO struct {
	*EvalDTO
	Result       string `json:"result,omitempty"`
	CompileError string `json:"compile-error,omitempty"`
	RuntimeError string `json:"runtime-error,omitempty"`
}
