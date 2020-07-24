package account

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"encoding/base64"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"gitlab.com/dracarys-botter/osrs-account-creator/pkg"
	"google.golang.org/api/gmail/v1"
	req "gitlab.com/dracarys-botter/osrs-account-creator/pkg/requests_helper"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// getStringInBetween Returns empty string if no start string found
func getStringInBetween(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str, end)
	if e == -1 {
		return
	}
	return str[s:e]
}

func getVerificationLink(input string) (output string) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		txt := scanner.Text()
		if strings.Contains(txt, "submit_code.ws") {
			result := getStringInBetween(txt, "(", ")")
			return result
		}
	}
	return ""
}

// GetAllEmailsGmail to retrieve all Gmail emails
func GetAllEmailsGmail(srv *gmail.Service, accountConfig pkg.AccountConfig) (allMessages []*gmail.Message, err error) {
	allMessages = make([]*gmail.Message, 0)
	messages, err := srv.Users.Messages.List("me").Do()
	if err != nil {
		return allMessages, err
	}
	for _, message := range messages.Messages {
		allMessages = append(allMessages, message)
	}
	nextToken := messages.NextPageToken
	for {
		if nextToken == "" {
			break
		}
		messages, err = srv.Users.Messages.List("me").PageToken(nextToken).Do()
		for _, message := range messages.Messages {
			allMessages = append(allMessages, message)
		}
		nextToken = messages.NextPageToken
	}
	return allMessages, nil
}

// FilterOSRSEmailsGmail filters emails from a list of gmail messages
func FilterOSRSEmailsGmail(srv *gmail.Service, messages []*gmail.Message) (osrsEmails []*gmail.Message, err error) {
	osrsEmails = make([]*gmail.Message, 0)
	for _, message := range messages {
		fullMsg, err := srv.Users.Messages.Get("me", message.Id).Format("full").Do()
		if err != nil {
			return osrsEmails, fmt.Errorf("Got error getting messages: %v", err)
		}
		if strings.Contains(strings.ToLower(fullMsg.Snippet), "runescape") {
			osrsEmails = append(osrsEmails, message)
		}
	}
	return osrsEmails, nil
}

// VerifyOSRSEmailGmail to verify all OSRS emails
func VerifyOSRSEmailGmail(srv *gmail.Service, config pkg.AccountConfig) (err error) {
	emails, err := GetAllEmailsGmail(srv, config)
	if err != nil {
		return fmt.Errorf("Unable to get all emails: %s", err)
	}
	filteredEmails, err := FilterOSRSEmailsGmail(srv, emails)
	for _, message := range filteredEmails {
		fullMsg, err := srv.Users.Messages.Get("me", message.Id).Format("full").Do()
		if err != nil {
			fmt.Printf("Got error getting messages: %v", err)
		}
		if fullMsg.Payload != nil {
			for _, part := range fullMsg.Payload.Parts {
				decoded, _ := base64.StdEncoding.DecodeString(part.Body.Data)
				decodedStr := string(decoded)
				link := getVerificationLink(decodedStr)
				if len(link) > 1 {
					parsedLink := strings.TrimSpace(link)
					req.VerifyAccount(parsedLink, config.ProxyConfig)
					break
				}
			}

		}
	}
	return nil
}

// DoAccountVerificationGmail TODO
func DoAccountVerificationGmail(credentialsPath string, config pkg.AccountConfig) (err error) {
	srv, err := LoginGmail(credentialsPath)
	if err != nil {
		return err
	}
	err = VerifyOSRSEmailGmail(srv, config)
	if err != nil {
		return err
	}
	return nil
}

// LoginGmail TODO
func LoginGmail(credentialsPath string) (srv *gmail.Service, err error) {
	b, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	fmt.Printf("%s", string(b))

	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err = gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	return srv, nil
}

