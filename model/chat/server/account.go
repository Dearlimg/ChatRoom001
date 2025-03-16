package server

type UpdateEmail struct {
	EnToken string `json:"en_token,omitempty"`
	Email   string `json:"email,omitempty"`
}

type UpdateAccount struct {
	EnToken   string `json:"en_token,omitempty"`
	Name      string `json:"name,omitempty"`
	Gender    string `json:"gender,omitempty"`
	Signature string `json:"signature,omitempty"`
}
