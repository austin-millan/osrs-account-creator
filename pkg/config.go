package pkg

import "fmt"

const (
	// RunescapeURL TODO
	RunescapeURL = "https://secure.runescape.com/m=account-creation/create_account"
	// RunecapeCaptchaSiteKey TODO
	RunecapeCaptchaSiteKey = "6Lcsv3oUAAAAAGFhlKrkRb029OHio098bbeyi_Hv"
	// CaptchaRetries TODO
	CaptchaRetries = 5
	// ProtonMailLoginURL TODO
	ProtonMailLoginURL = "https://mail.protonmail.com/login"
)

// ClientDriverMode TODO
type ClientDriverMode string

const (
	// RequestMode TODO
	RequestMode ClientDriverMode = "requests"
	// SeleniumMode TODO
	SeleniumMode ClientDriverMode = "selenium"
)

// ProxyConfig TODO
type ProxyConfig struct {
	IP   string
	Port string
	User string
	Pass string
}

// AccountData TODO
type AccountData struct {
	Password      string
	BirthdayDay   string
	BirthdayMonth string
	BirthdayYear  string
}

// AccountConfig TODO
type AccountConfig struct {
	Email     string
	ProxyConfig ProxyConfig
	AccountData AccountData
}

// NewAccountOutput TODO
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

// ShowAccountOutput TODO
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
