package entity

// Response is a struct with health response definition
type Response struct {
	Name    string `json:"name"`
	Healthy bool   `json:"healthy"`
	Time    string `json:"time"`
	Message string `json:"message,omitempty"`
}
