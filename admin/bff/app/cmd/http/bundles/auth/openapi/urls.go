package openapi

// Внимание!
// Url должны совпадать с генерируемыми на основе `spec/openapi.yml`!

const (
	UrlSignInWelcome      = "/auth/sign-in/welcome"
	UrlSignInRequest      = "/auth/sign-in/request"
	UrlSignInRequestRetry = "/auth/sign-in/request-retry"
	UrlSignInConfirm      = "/auth/sign-in/confirm"
	UrlSignOut            = "/auth/sign-out"
)
