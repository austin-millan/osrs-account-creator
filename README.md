# osrs-account-creator

## About

A library for creating OSRS accounts behind a proxy and verifying accounts using GMail API.

## Setup

```bash
go get -u gitlab.com/dracarys-botter/osrs-account-creator
```

## Requirements

Only works behind a SOCKS5 proxy. This is intended to avoid tainting your home IP address.

## Example Usage

```go
package main

import (
    "fmt"

    "gitlab.com/dracarys-botter/osrs-account-creator/pkg"
    account "gitlab.com/dracarys-botter/osrs-account-creator/pkg/account"
)

func main() {
    accConfig := pkg.AccountConfig{
        Email: "youremail@domain.com",
        ProxyConfig: pkg.ProxyConfig{
            IP:   "127.0.0.1",
            User: "someid",
            Pass: "somepass",
            Port: "1080",
        },
    }

    if acc, err := account.RegisterAccount(accConfig, pkg.RequestMode, "your 2captcha.com API key"); err != nil {
        fmt.Printf("Error: %v", err)
    } else {
        if acc != nil {
            pkg.ShowAccountOutput(*acc)
        }
    }
    // your Gmail credentials file, see more here:
    //      https://developers.google.com/gmail/api/quickstart/go
    account.DoAccountVerificationGmail("credentials.json", accConfig)
}
```

## Issues

- Only works behind SOCKS5 proxy, but ideally will allow multiple protocols (HTTP/HTTPS) and the ability to wave the proxy altogether
