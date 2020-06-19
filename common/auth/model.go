package auth

type (
	// Authorization ...
	Authorization struct {
		AppName string `json:"appname" bson:"appname"`
		Role    string `json:"role" bson:"role"`
	}

	// Authentication ...
	Authentication struct {
		ID             string          `json:"username"`
		Authorizations []Authorization `json:"authorizations"`
	}
)
