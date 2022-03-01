package web

type ErrorDTO struct {
	Error string `json:"errorMessage"`
}

type SessionIdentifierDTO struct {
	ID string `json:"sessionId"`
}
