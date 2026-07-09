package response

import "net/http"

// InvalidCredentials indicates the username or password is incorrect.
var (
	InvalidUrlsProvided = NewCodeObject(
		http.StatusBadRequest,
		"ZIPURL_INVUP401",
		map[string]string{
			"en": "Invalid short url provided",
			"bn": "ভুল শর্ট ইউআরএল প্রদান করা হয়েছে",
		},
		nil,
	)

	TechnicalError400 = NewCodeObject(
		http.StatusBadRequest,
		"ZIPURL_TE400",
		map[string]string{
			"en": "Something went wrong, please try again later",
			"bn": "কিছু ভুল হয়েছে, অনুগ্রহ করে পরে আবার চেষ্টা করুন",
		},
		nil,
	)

	DataValidationErr400 = NewCodeObject(
		http.StatusBadRequest,
		"ZIPURL_DVE400",
		map[string]string{
			"en": "Invalid request data",
			"bn": "অবৈধ অনুরোধ ডেটা",
		},
		nil,
	)

	UserSignUpFailed = NewCodeObject(
		http.StatusBadRequest,
		"ZIPURL_USUF",
		map[string]string{
			"en": "User Sign Up Failed",
			"bn": "ইউজার সাইন আপ ব্যর্থ হয়েছে",
		},
		nil,
	)

	ShortURLCreationFailed = NewCodeObject(
		http.StatusBadRequest,
		"ZIPURL_SUCF",
		map[string]string{
			"en": "Failed to create short url",
			"bn": "Failed to create short url",
		},
		nil,
	)

	UserAlreadyExistsWithEmail = NewCodeObject(
		http.StatusBadRequest,
		"ZIPURL_UEE",
		map[string]string{
			"en": "An user with this email address already exists.",
			"bn": "এই ইমেইল দিয়ে ইতোমধ্যে একজন ইউজার নিবন্ধিত রয়েছেন।",
		},
		nil,
	)
)
