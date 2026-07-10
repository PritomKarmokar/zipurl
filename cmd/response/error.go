package response

import "net/http"

// InvalidCredentials indicates the username or password is incorrect.
var (
	InvalidUrlsProvided = NewCodeObject(
		http.StatusBadRequest,
		"ZIPURL_Invalid_Url_Provided",
		map[string]string{
			"en": "Invalid short url provided",
			"bn": "ভুল শর্ট ইউআরএল প্রদান করা হয়েছে",
		},
		nil,
	)

	TechnicalError = NewCodeObject(
		http.StatusInternalServerError,
		"ZIPURL_Technical_Error",
		map[string]string{
			"en": "Due to the unexpected Issues We couldn't process your request, please try again later",
			"bn": "কিছু ভুল হয়েছে, অনুগ্রহ করে পরে আবার চেষ্টা করুন",
		},
		nil,
	)

	DataValidationErr = NewCodeObject(
		http.StatusBadRequest,
		"ZIPURL_Data_Validation_Err",
		map[string]string{
			"en": "Invalid request data",
			"bn": "অবৈধ অনুরোধ ডেটা",
		},
		nil,
	)

	UserSignUpFailed = NewCodeObject(
		http.StatusInternalServerError,
		"ZIPURL_User_SignUp_Failed",
		map[string]string{
			"en": "User Sign Up Failed",
			"bn": "ইউজার সাইন আপ ব্যর্থ হয়েছে",
		},
		nil,
	)

	ShortURLCreationFailed = NewCodeObject(
		http.StatusBadRequest,
		"ZIPURL_Short_URL_Creation_Failed",
		map[string]string{
			"en": "Failed to create short url",
			"bn": "Failed to create short url",
		},
		nil,
	)

	UserAlreadyExistsWithEmail = NewCodeObject(
		http.StatusConflict,
		"ZIPURL_User_Already_Exists_With_This_Email",
		map[string]string{
			"en": "An user with this email address already exists.",
			"bn": "এই ইমেইল দিয়ে ইতোমধ্যে একজন ইউজার নিবন্ধিত রয়েছেন।",
		},
		nil,
	)
)
