package response

import "net/http"

var GenericSuccess200 = NewCodeObject(
	http.StatusOK,
	"GENERIC_SUC200",
	map[string]string{
		"en": "Request processed successfully",
		"bn": "অনুরোধ সফলভাবে প্রক্রিয়াকৃত হয়েছে",
	},
	nil,
)

var UserSignUpSuccess = NewCodeObject(
	http.StatusCreated,
	"ZIP_URL_USUS",
	map[string]string{
		"en": "User Sign Up Successful",
		"bn": "ইউজার সাইন আপ সফল হয়েছে",
	},
	nil,
)
