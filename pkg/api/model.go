package api

// This interface is used to communicate with a specific backend API.

// ApiModel defines handler object for different backend models.
// This should be abstract class.
type ApiModel struct {
	Client ApiClient
	Name   string
}

// NewAPI initializes a specific API.
func NewAPI(model string) ApiModel {
	return ApiModel{
		Client: NewClient(model),
		Name:   model,
	}
}
