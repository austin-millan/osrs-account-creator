package selenium_helper

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/tebeka/selenium"
	"gitlab.com/dracarys-botter/osrs-account-creator/pkg"
)

func newProxyCapability(caps *selenium.Capabilities, ip, port, user, pass string) error {
	proxy := selenium.Proxy{SOCKSVersion: 5, Type: selenium.Manual}
	if caps == nil {
		return fmt.Errorf("Missing capabilities")
	}
	if i, err := strconv.Atoi(port); err != nil {
		return fmt.Errorf("Invalid port: %s", port)
	} else {
		proxy.SocksPort = i
	}
	if net.ParseIP(ip).To4() == nil {
		return fmt.Errorf("Invalid IP address: %s", ip)
	}
	proxy.SOCKS = ip
	proxy.SOCKSUsername = user
	proxy.SOCKSPassword = pass
	caps.AddProxy(proxy)
	return nil
}

// GetWebdriver TODO
func GetWebdriver() (output selenium.WebDriver, err error) {
	caps := selenium.Capabilities{"browserName": "firefox"}
	port, err := getUnusedPort()
	opts := []selenium.ServiceOption{
		selenium.GeckoDriver("geckodriver"), // Specify the path to GeckoDriver in order to use Firefox.
	}
	_, err = selenium.NewSeleniumService("selenium-server.jar", port, opts...)
	if err != nil {
		return nil, fmt.Errorf("Unable to load chromedriver: %s", err.Error())
	}
	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:"+strconv.Itoa(port)+"/wd/hub")
	if err != nil {
		return nil, fmt.Errorf("Unable to get Webdriver: %s", err.Error())
	}
	return wd, err
}

// GetSocksProxiedWebdriver TODO
func GetSocksProxiedWebdriver(config *pkg.ProxyConfig) (output selenium.WebDriver, err error) {
	if config == nil {
		return nil, fmt.Errorf("Must provide ProxyConfig")
	}
	caps := selenium.Capabilities{"browserName": "firefox"}
	err = newProxyCapability(&caps, config.IP, config.Port, config.User, config.Pass)
	if err != nil {
		return nil, err
	}
	opts := []selenium.ServiceOption{
		selenium.StartFrameBuffer(),         // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver("geckodriver"), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr),          // Output debug information to STDERR.
	}
	port, err := getUnusedPort()
	_, err = selenium.NewSeleniumService("selenium-server.jar", port, opts...)
	if err != nil {
		return nil, fmt.Errorf("Unable to load chromedriver: %s", err.Error())
	}
	// if service != nil {
	// 	defer service.Stop()
	// }
	if err != nil {
		return nil, fmt.Errorf("Unable to find unused port: %s", err.Error())
	}
	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:"+strconv.Itoa(port)+"/wd/hub")
	if err != nil {
		return nil, fmt.Errorf("Unable to get Webdriver: %s", err.Error())
	}
	// defer wd.Quit()

	return wd, err
}

func getUnusedPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		return 0, err
	}
	return port, nil
}

// CreateAccount TODO
func CreateAccount(account pkg.AccountConfig, twoCaptchaAPIKey string) (output *pkg.NewAccountOutput, err error) {
	output = &pkg.NewAccountOutput{}
	var wd selenium.WebDriver
	if account.ProxyConfig.IP != "" {
		wd, err = GetSocksProxiedWebdriver(&pkg.ProxyConfig{
			IP:   account.ProxyConfig.IP,
			Port: account.ProxyConfig.Port,
			User: account.ProxyConfig.User,
			Pass: account.AccountData.Password,
		})
	} else {
		wd, err = GetWebdriver()
	}
	if err != nil {
		return output, err
	}
	wd.Get(pkg.RunescapeURL)
	time.Sleep(time.Second)
	forms, err := wd.FindElements(selenium.ByXPATH, "//*[@id=\"create-email\"]")
	if err != nil {
		panic(err)
	}

	// var text string
	for _, form := range forms {
		inputs, err := form.FindElements(selenium.ByXPATH, "//input")
		if err != nil {
			fmt.Printf("input could not be found %s", err.Error())
		}
		for _, input := range inputs {
			InputName, err := input.GetAttribute("name")
			if err != nil {
				fmt.Printf("input name could not be found: %s", err.Error())
			}
			InputType, err := input.GetAttribute("type")
			if err != nil {
				fmt.Printf("input type could not be found %s", err.Error())
			} else if InputType == "email" {
				input.SendKeys(account.Email)
				time.Sleep(time.Second)
			} else if InputType == "password" {
				input.SendKeys(account.AccountData.Password)
				time.Sleep(time.Second)
			} else if InputName == "onlyOneEmail" {
				input.SendKeys("1")
				time.Sleep(time.Second)
			} else if InputName == "day" {
				input.SendKeys(account.AccountData.BirthdayDay)
				time.Sleep(time.Second)
			} else if InputName == "month" {
				input.SendKeys(account.AccountData.BirthdayMonth)
				time.Sleep(time.Second)
			} else if InputName == "year" {
				input.SendKeys(account.AccountData.BirthdayYear)
				time.Sleep(time.Second)
			} else if InputName == "create-submit" {
				input.SendKeys("create")
				time.Sleep(time.Second)
			}
		}
		// Cookie
		button, err := wd.FindElement(selenium.ByXPATH, "/html/body/div[1]/div/a")
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
		err = button.Click()

		button, err = wd.FindElement(selenium.ByXPATH, "//*[@id=\"create-submit\"]")
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
		err = button.Click()
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
	}

	inputs, err := wd.FindElements(selenium.ByXPATH, "//*[@id=\"recaptcha-token\"]")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	for _, input := range inputs {
		InputName, err := input.GetAttribute("name")
		if err != nil {
			fmt.Printf("input name could not be found: %s", err.Error())
		}
		InputType, err := input.GetAttribute("type")
		fmt.Printf("%s, %s", InputName, InputType)
	}
	inputs, err = wd.FindElements(selenium.ByID, "recaptcha-token")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	input, err := wd.FindElement(selenium.ByID, "recaptcha-token")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	fmt.Printf("%v", input)
	for _, input := range inputs {
		InputName, err := input.GetAttribute("name")
		if err != nil {
			fmt.Printf("input name could not be found: %s", err.Error())
		}
		InputType, err := input.GetAttribute("type")
		fmt.Printf("%s, %s", InputName, InputType)
	}
	inputs, err = wd.FindElements(selenium.ByClassName, "rc-imageselect-challenge")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	inputs, err = wd.FindElements(selenium.ByXPATH, "//p[contains(text(), 'Please complete the reCAPTCHA box.')]")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	input, err = wd.FindElement(selenium.ByXPATH, "//p[contains(text(), 'Please complete the reCAPTCHA box.')]")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	fmt.Printf("%v", input)

	for _, input := range inputs {
		InputName, err := input.GetAttribute("name")
		if err != nil {
			fmt.Printf("input name could not be found: %s", err.Error())
		}
		InputType, err := input.GetAttribute("type")
		fmt.Printf("%s, %s", InputName, InputType)
	}
	if input, err := wd.FindElement(selenium.ByXPATH, "//*[@id=\"recaptcha-token\"]"); err != nil {
		fmt.Printf("Error: %v", err)
	} else {
		InputName, err := input.GetAttribute("name")
		if err != nil {
			fmt.Printf("input name could not be found: %s", err.Error())
		}
		InputType, err := input.GetAttribute("type")
		fmt.Printf("%s, %s", InputName, InputType)
		for _, input := range inputs {
			InputName, err := input.GetAttribute("name")
			if err != nil {
				fmt.Printf("input name could not be found: %s", err.Error())
			}
			InputType, err := input.GetAttribute("type")
			fmt.Printf("%s, %s", InputName, InputType)
		}
	}

	input, err = wd.FindElement(selenium.ByXPATH, "/html/body/div/div[2]/iframe")
	if err != nil {
		fmt.Printf("%s", err.Error())
	} else {
		showElement(input)
	}

	input, err = wd.FindElement(selenium.ByID, "google-recaptcha")
	if err != nil {
		fmt.Printf("%s", err.Error())
	} else {
		showElement(input)
	}

	input, err = wd.FindElement(selenium.ByID, "recaptcha_box")
	if err != nil {
		fmt.Printf("%s", err.Error())
	} else {
		showElement(input)
	}

	input, err = wd.FindElement(selenium.ByID, "recaptcha_box")
	if err != nil {
		fmt.Printf("%s", err.Error())
	} else {
		showElement(input)
	}
	return output, nil
}

func showElement(input selenium.WebElement) {
	InputName, _ := input.GetAttribute("name")
	InputType, _ := input.GetAttribute("type")
	InputID, _ := input.GetAttribute("id")
	InputClass, _ := input.GetAttribute("class")
	InputStatus, _ := input.GetAttribute("status")
	fmt.Printf("%s", InputName)
	fmt.Printf("%s", InputType)
	fmt.Printf("%s", InputID)
	fmt.Printf("%s", InputClass)
	fmt.Printf("%s", InputStatus)
}

// ProtonmailConfig TODO
type ProtonmailConfig struct {
	// Email TODO
	Email string
	// Password TODO
	Password string
	// EncryptionPassword TODO
	EncryptionPassword string
}

// LoginProtonmail TODO
func LoginProtonmail(cfg ProtonmailConfig) (err error) {
	// var solution string
	var wd selenium.WebDriver
	wd, err = GetWebdriver()
	if err != nil {
		return err
	}
	err = wd.Get(pkg.ProtonMailLoginURL)
	if err != nil {
		return err
	}
	time.Sleep(time.Second*3)

	form, err := wd.FindElement(selenium.ByXPATH, "//*[@id=\"username\"]")
	if err != nil {
		panic(err)
	}
	form.SendKeys(cfg.Email)

	form, err = wd.FindElement(selenium.ByXPATH, "//*[@id=\"password\"]")
	if err != nil {
		panic(err)
	}
	form.SendKeys(cfg.Password)

	time.Sleep(time.Second*1)

	// LOGIN
	button, err := wd.FindElement(selenium.ByXPATH, "//*[@id=\"login_btn\"]")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	time.Sleep(time.Second*1)
	err = button.Click()
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	time.Sleep(time.Second*4)

	form, err = wd.FindElement(selenium.ByXPATH, "//*[@id=\"mailboxPassword\"]")
	if err != nil {
		panic(err)
	}
	form.SendKeys(cfg.EncryptionPassword)
	time.Sleep(time.Second*3)

	button, err = wd.FindElement(selenium.ByXPATH, "//*[@id=\"unlock_btn\"]")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	err = button.Click()
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	time.Sleep(5)
	emails, err := wd.FindElements(selenium.ByCSSSelector, "div.conversation:nth-child(1)")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	for _, email := range emails {
		showElement(email)
	}

	emails, err = wd.FindElements(selenium.ByName, "conversation-meta")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	for _, email := range emails {
		showElement(email)
	}

	emails, err = wd.FindElements(selenium.ByXPATH, "/html/body/div[2]/div[2]/div/div[1]/section/div[1]/div[2]")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	for _, email := range emails {
		showElement(email)
	}

	emails, err = wd.FindElements(selenium.ByCSSSelector, "div.conversation:nth-child(1) > div:nth-child(3)")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	for _, email := range emails {
		showElement(email)
	}

	emails, err = wd.FindElements(selenium.ByName, "subject-text ellipsis")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	for _, email := range emails {
		showElement(email)
	}

	emails, err = wd.FindElements(selenium.ByName, "row top")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	for _, email := range emails {
		showElement(email)
	}


	emails, err = wd.FindElements(selenium.ByID, "conversation-list-columns")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	for _, email := range emails {
		showElement(email)
	}


	emails, err = wd.FindElements(selenium.ByXPATH, "//*[@id=\"conversation-list-columns\"]")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	for _, email := range emails {
		showElement(email)
	}

	//*[@id="conversation-list-columns"]


	return nil
}
