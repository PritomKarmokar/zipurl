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
