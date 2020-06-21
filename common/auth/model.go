package auth

type (
	// Authorization ...
	Authorization struct {
		AppName string `json:"appname" bson:"appname,omitempty"`
		Role    string `json:"role" bson:"role,omitempty"`
	}

	// Authentication ...
	Authentication struct {
		ID             string          `json:"username" bson:"username,omitempty"`
		Authorizations []Authorization `json:"authorizations" bson:"authorizations,omitempty"`
	}
)
