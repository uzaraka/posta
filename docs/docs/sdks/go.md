---
sidebar_position: 2
title: Go
description: Posta Go SDK
---

# Go SDK

Official Go client for Posta, built with [kitoko](https://github.com/jkaninda/kitoko).

## Installation

```bash
go get github.com/goposta/posta-go
```

**Requires:** Go 1.25+

## Quick Start

```go
package main

import (
    "fmt"
    "log"

    posta "github.com/goposta/posta-go"
)

func main() {
    client := posta.New("https://posta.example.com", "your-api-key")

    resp, err := client.SendEmail(&posta.SendEmailRequest{
        From:    "sender@example.com",
        To:      []string{"recipient@example.com"},
        Subject: "Hello from Posta",
        HTML:    "<h1>Hello!</h1><p>This is a test email.</p>",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Email sent: id=%s status=%s\n", resp.ID, resp.Status)
}
```

## Send Template Email

```go
resp, err := client.SendTemplateEmail(&posta.SendTemplateEmailRequest{
    Template: "welcome",
    To:       []string{"user@example.com"},
    From:     "noreply@example.com",
    Language: "en",
    TemplateData: map[string]any{
        "name": "Alice",
    },
})
```

## Batch Send

```go
resp, err := client.SendBatch(&posta.BatchRequest{
    Template: "newsletter",
    From:     "news@example.com",
    Recipients: []posta.BatchRecipient{
        {Email: "user1@example.com", TemplateData: map[string]any{"name": "Bob"}},
        {Email: "user2@example.com", Language: "fr", TemplateData: map[string]any{"name": "Carol"}},
    },
})
fmt.Printf("Sent: %d, Failed: %d\n", resp.Sent, resp.Failed)
```

## Check Delivery Status

```go
status, err := client.GetEmailStatus("email-uuid")
fmt.Printf("Status: %s\n", status.Status)
```

## Error Handling

```go
_, err := client.GetEmailStatus("invalid-uuid")
if err != nil {
    var apiErr *posta.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("Status: %d\n", apiErr.StatusCode)
        if apiErr.Info != nil {
            fmt.Printf("Message: %s\n", apiErr.Info.Message)
        }
    }
}
```
