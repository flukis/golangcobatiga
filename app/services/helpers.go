package services

type ServiceReturn struct {
	HttpStatusCode int
	Err            error
	Payload        any
}

type SignInPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
