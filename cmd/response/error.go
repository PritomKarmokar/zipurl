package response

import "net/http"

// InvalidCredentials indicates the username or password is incorrect.
var InvalidCredentials = NewCodeObject(
	http.StatusUnauthorized,
	"ZIPURL_INV401",
	map[string]string{
		"en": "Invalid username or password",
		"bn": "ভুল ইউজারনেম বা পাসওয়ার্ড",
	},
	nil,
)
