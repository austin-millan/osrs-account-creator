package requests_helper

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"net/http"
	"net/url"
	"time"

	"github.com/austin-millan/twocaptcha/pkg/twocaptcha"
	"gitlab.com/dracarys-botter/osrs-account-creator/pkg"
	"golang.org/x/net/proxy"
)


// NewHTTPClient is a constructor for the Webull-Client client
func NewHTTPClient() (c *http.Client, err error) {
	// Create client
	c = &http.Client{
		Transport: &http.Transport{TLSHandshakeTimeout: 10 * time.Second},
		Timeout: time.Second * 20,
	}
	return c, nil
}

// NewProxiedHTTPClient is a constructor for the Webull-Client client
func NewProxiedHTTPClient(config *pkg.ProxyConfig) (c *http.Client, err error) {
	var dialSocksProxy proxy.Dialer
	addr := fmt.Sprintf("%s:%s", config.IP, config.Port)
	if config.User != "" && config.Pass != "" {
		auth := &proxy.Auth{
			User:     config.User,
			Password: config.Pass,
		}
		dialSocksProxy, err = proxy.SOCKS5("tcp", addr, auth, &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		})
		if err != nil {
			fmt.Println("Error connecting to proxy:", err)
		}

	} else {
		dialSocksProxy, err = proxy.SOCKS5("tcp", addr, nil, proxy.Direct)
		if err != nil {
			fmt.Println("Error connecting to proxy:", err)
		}
	}
	// Create client
	c = &http.Client{
		Transport: &http.Transport{
			Dial:                dialSocksProxy.Dial,
			TLSHandshakeTimeout: 10 * time.Second,
		},
		Timeout: time.Second * 20,
	}
	return
}


// VerifyAccount verifies an account behind a proxy
func VerifyAccount(verifyAccountURL string, config pkg.ProxyConfig) (err error) {
	c, err := NewProxiedHTTPClient(&config)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodGet, verifyAccountURL, nil)
	setRunescapeCommonHeaders(req)
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("Got error sending request %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Unable to read data %s", err.Error())
	}
	stringified := string(body)
	if strings.Contains(stringified, "The link you clicked has already been used") {
		fmt.Printf("Account already verified")
	} else {
		fmt.Printf("Probably a success")
	}
	fmt.Printf("%s", stringified)
	return nil
}

// CreateAccount creates an account and solves captcha AOT
func CreateAccount(account pkg.AccountConfig, twoCaptchaAPIKey string) (output *pkg.NewAccountOutput, err error) {
	var solution string
	cli, err := NewProxiedHTTPClient(&pkg.ProxyConfig{
		IP:   account.ProxyConfig.IP,
		Port: account.ProxyConfig.Port,
		User: account.ProxyConfig.User,
		Pass: account.AccountData.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to get client: %s", err.Error())
	}
	requestBody := map[string]string{
		"email1":               account.Email,
		"onlyOneEmail":         "1",
		"password1":            account.AccountData.Password,
		"onlyOnePassword":      "1",
		"day":                  account.AccountData.BirthdayDay,
		"month":                account.AccountData.BirthdayMonth,
		"year":                 account.AccountData.BirthdayYear,
		"create-submit":        "create",
		"g-recaptcha-response": "",
	}
	instance, err := twocaptcha.NewInstance(twoCaptchaAPIKey, twocaptcha.SettingInfo{TimeBetweenRequests: 10})
	if err != nil {
		return nil, fmt.Errorf("Unable to get 2Captcha client: %s", err.Error())
	}

	currentTime := time.Now()
	var timeSolved time.Time

	for i := 1; i < pkg.CaptchaRetries+1; i++ {
		if solution, err = instance.SolveRecaptchaV2(pkg.RunecapeCaptchaSiteKey, pkg.RunescapeURL); err != nil {
			fmt.Printf("Got error solving recaptcha: %s", err.Error())
		} else {
			requestBody["g-recaptcha-response"] = solution
			timeSolved = time.Now()
			break
		}
	}
	timeToSolve := timeSolved.Sub(currentTime)
	fmt.Printf("Time to solve captcha: %fs", timeToSolve.Seconds())

	formValues := url.Values{}
	for k, v := range requestBody {
		formValues[k] = []string{v}
	}
	encoded := strings.NewReader(formValues.Encode())

	req, err := http.NewRequest(http.MethodPost, pkg.RunescapeURL, encoded)
	setRunescapeCommonHeaders(req)

	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Got error sending request %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to read data %s", err.Error())
	}
	stringified := string(body)
	fmt.Printf("%s", stringified)
	if strings.Contains(stringified, "Account Created - RuneScape") {
		fmt.Printf("Account created")
	} else {
		return nil, fmt.Errorf("Account not created")
	}
	output = &pkg.NewAccountOutput{
		Email:         account.Email,
		ProxyIP:       account.ProxyConfig.IP,
		ProxyPort:     account.ProxyConfig.Port,
		ProxyUser:     account.ProxyConfig.User,
		ProxyPass:     account.ProxyConfig.Pass,
		BirthdayDay:   account.AccountData.BirthdayDay,
		BirthdayMonth: account.AccountData.BirthdayMonth,
		BirthdayYear:  account.AccountData.BirthdayYear,
		Recaptcha:     solution,
	}

	return output, nil
}

// // LoginProtonmail TODO
// func LoginProtonmail(cfg ProtonmailConfig) (output *pkg.NewAccountOutput, err error) {
// 	// var solution string
// 	cli, err := NewHTTPClient()
// 	if err != nil {
// 		return nil, fmt.Errorf("Unable to get client: %s", err.Error())
// 	}

// 	req, err := http.NewRequest(http.MethodGet, pkg.ProtonMailLoginURL, nil)
// 	setProtonmailCommonHeaders(req)

// 	resp, err := cli.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("Got error sending request %s", err.Error())
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("Unable to read data %s", err.Error())
// 	}
// 	stringified := string(body)
// 	fmt.Printf("%s", stringified)
// 	if strings.Contains(stringified, "Account Created - RuneScape") {
// 		fmt.Printf("Account created")
// 	} else {
// 		fmt.Printf("Account NOT created")
// 	}
// 	output = &pkg.NewAccountOutput{
// 		Email:         cfg.Email,
// 	}

// 	return output, nil
// }


func setRunescapeCommonHeaders(req *http.Request) {
	if req == nil {
		return
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36")
	req.Header.Add("Origin", "https://secure.runescape.com/")
	req.Header.Add("Referer", "https://secure.runescape.com/m=account-creation/create_account?theme=oldschool")
	req.Header.Add("DNT", "1")
	req.Header.Add("TE", "Trailers")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	// req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "en-GB,en-US;q=0.8,en;q=0.6")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Host", "secure.runescape.com")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
}

func setProtonmailCommonHeaders(req *http.Request) {
	if req == nil {
		return
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 Safari/537.36")
	req.Header.Add("DNT", "1")
	req.Header.Add("TE", "Trailers")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-GB,en-US;q=0.8,en;q=0.6")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Host", "secure.runescape.com")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
}