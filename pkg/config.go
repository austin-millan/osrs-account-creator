package pkg

import "fmt"

const (
	// RunescapeURL is Runescape account creation page
	RunescapeURL = "https://secure.runescape.com/m=account-creation/create_account"
	// RunecapeCaptchaSiteKey is Runescape's captcha key
	RunecapeCaptchaSiteKey = "6Lcsv3oUAAAAAGFhlKrkRb029OHio098bbeyi_Hv"
	// CaptchaRetries declares number of attempts to retry solving a captcha
	CaptchaRetries = 5
	// ProtonMailLoginURL is ProtonMail account login page
	ProtonMailLoginURL = "https://mail.protonmail.com/login"
)

// ClientDriverMode enum
type ClientDriverMode string

const (
	// RequestMode to do account creation/verification via HTTP requests
	RequestMode ClientDriverMode = "requests"
	// SeleniumMode to do account creation/verification via Selenium webdriver
	SeleniumMode ClientDriverMode = "selenium"
)

// ProxyConfig for proxy (used in each `ClientDriverMode`s)
type ProxyConfig struct {
	IP   string
	Port string
	User string
	Pass string
}

// AccountData is used internally for storing form data
type AccountData struct {
	Password      string
	BirthdayDay   string
	BirthdayMonth string
	BirthdayYear  string
}

// AccountConfig is used as input for registering accounts
type AccountConfig struct {
	Email     string
	ProxyConfig ProxyConfig
	AccountData AccountData
}

// NewAccountOutput is the output of a Register account operation, user should save this
type NewAccountOutput struct {
	Email         string
	ProxyIP       string
	ProxyPort     string
	ProxyUser     string
	ProxyPass     string
	BirthdayDay   string
	BirthdayMonth string
	BirthdayYear  string
	Recaptcha     string
}

// ShowAccountOutput is a helper function for displaying NewAccountOutput
func ShowAccountOutput(input NewAccountOutput) {
	fmt.Printf("Email: %s", input.Email)
	fmt.Printf("ProxyIP: %s", input.ProxyIP)
	fmt.Printf("ProxyPort: %s", input.ProxyPort)
	fmt.Printf("ProxyUser: %s", input.ProxyUser)
	fmt.Printf("ProxyPass: %s", input.ProxyPass)
	fmt.Printf("BirthdayDay: %s", input.BirthdayDay)
	fmt.Printf("BirthdayMonth: %s", input.BirthdayMonth)
	fmt.Printf("BirthdayYear: %s", input.BirthdayYear)
	fmt.Printf("Recaptcha: %s", input.Recaptcha)
}
