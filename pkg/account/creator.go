package account

import (
	"fmt"
	"math/rand"
	"time"

	"gitlab.com/dracarys-botter/osrs-account-creator/pkg"
	req "gitlab.com/dracarys-botter/osrs-account-creator/pkg/requests_helper"
)


// RegisterAccounts wrapper for registering many accounts
func RegisterAccounts(accounts []pkg.AccountConfig, mode pkg.ClientDriverMode, twoCaptchaAPIKey string) (output []pkg.NewAccountOutput, err error) {
	output = make([]pkg.NewAccountOutput, 0)
	for _, acc := range accounts {
		newAccount, err := RegisterAccount(acc, mode, twoCaptchaAPIKey)
		if err != nil {
			fmt.Printf("Got error registering account %v: %s", acc, err.Error())
			continue
		}
		if newAccount != nil {
			output = append(output, *newAccount)
		}
	}
	return output, nil
}

func newPassword(length int) string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#&"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func newBirthday() (day, month, year string) {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	day = fmt.Sprintf("%d", (seededRand.Int()%28)+1)
	month = fmt.Sprintf("%d", (seededRand.Int()%12)+1)
	year = fmt.Sprintf("%d", (2010-(seededRand.Int()%50))+1)
	return
}

// RegisterAccount register an account and solve the captcha problem
func RegisterAccount(account pkg.AccountConfig, mode pkg.ClientDriverMode, twoCaptchaAPIKey string) (output *pkg.NewAccountOutput, err error) {
	pass := newPassword(12)
	day, month, year := newBirthday()
	account.AccountData = pkg.AccountData{
		Password:      pass,
		BirthdayDay:   day,
		BirthdayMonth: month,
		BirthdayYear:  year,
	}
	switch mode {
	case pkg.RequestMode:
		if output, err := req.CreateAccount(account, twoCaptchaAPIKey); err != nil {
			return nil, fmt.Errorf("err: %v", err)
		} else {
			fmt.Printf("output: %v", output)
		}
	}
	return nil, nil
}
