package api


//Login represents auth data
type Login struct {
	Email    string `json:"email" valid:"email~Invalid email"`
	Password string `json:"password" valid:"required~The password is required"`
}

//AuthenticatedUser represents an authed user
type AuthenticatedUser struct {
	ID    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
}
