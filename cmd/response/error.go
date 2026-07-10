package response

import "net/http"

// InvalidCredentials indicates the username or password is incorrect.
var (
	InvalidUrlsProvided = NewCodeObject(
		http.StatusBadRequest,
		"ZIP_IUP",
		map[string]string{
			"en": "Invalid short url provided",
			"bn": "ভুল শর্ট ইউআরএল প্রদান করা হয়েছে",
		},
		nil,
	)

	TechnicalError = NewCodeObject(
		http.StatusInternalServerError,
		"ZIP_TE",
		map[string]string{
			"en": "Due to the unexpected Issues We couldn't process your request, please try again later",
			"bn": "কিছু ভুল হয়েছে, অনুগ্রহ করে পরে আবার চেষ্টা করুন",
		},
		nil,
	)

	DataValidationErr = NewCodeObject(
		http.StatusBadRequest,
		"ZIP_DVE",
		map[string]string{
			"en": "Invalid request data",
			"bn": "অবৈধ অনুরোধ ডেটা",
		},
		nil,
	)

	UserSignUpFailed = NewCodeObject(
		http.StatusInternalServerError,
		"ZIP_USF",
		map[string]string{
			"en": "User Sign Up Failed",
			"bn": "ইউজার সাইন আপ ব্যর্থ হয়েছে",
		},
		nil,
	)

	ShortURLCreationFailed = NewCodeObject(
		http.StatusBadRequest,
		"ZIP_SUCF",
		map[string]string{
			"en": "Failed to create short url",
			"bn": "Failed to create short url",
		},
		nil,
	)

	UserAlreadyExistsWithEmail = NewCodeObject(
		http.StatusConflict,
		"ZIP_UAEWE",
		map[string]string{
			"en": "An user with this email address already exists.",
			"bn": "এই ইমেইল দিয়ে ইতোমধ্যে একজন ইউজার নিবন্ধিত রয়েছেন।",
		},
		nil,
	)
)
