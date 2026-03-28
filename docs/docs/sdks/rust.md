---
sidebar_position: 5
title: Rust
description: Posta Rust SDK
---

# Rust SDK

Official async Rust client for Posta, built with [reqwest](https://docs.rs/reqwest).

## Installation

Add to your `Cargo.toml`:

```toml
[dependencies]
posta = "0.1"
tokio = { version = "1", features = ["full"] }
```

**Requires:** Rust 2021 edition

## Quick Start

```rust
use posta::{Client, SendEmailRequest};

#[tokio::main]
async fn main() -> Result<(), posta::Error> {
    let client = Client::new("https://posta.example.com", "your-api-key");

    let resp = client.send_email(&SendEmailRequest {
        from: "sender@example.com".into(),
        to: vec!["recipient@example.com".into()],
        subject: "Hello from Posta".into(),
        html: Some("<h1>Hello!</h1><p>This is a test email.</p>".into()),
        ..Default::default()
    }).await?;

    println!("Email sent: id={} status={}", resp.id, resp.status);
    Ok(())
}
```

## Send Template Email

```rust
use std::collections::HashMap;
use posta::SendTemplateEmailRequest;

let resp = client.send_template_email(&SendTemplateEmailRequest {
    template: "welcome".into(),
    to: vec!["user@example.com".into()],
    from: Some("noreply@example.com".into()),
    language: Some("en".into()),
    template_data: Some(HashMap::from([
        ("name".into(), "Alice".into()),
    ])),
    ..Default::default()
}).await?;
```

## Batch Send

```rust
use posta::{BatchRequest, BatchRecipient};

let resp = client.send_batch(&BatchRequest {
    template: "newsletter".into(),
    from: Some("news@example.com".into()),
    recipients: vec![
        BatchRecipient {
            email: "user1@example.com".into(),
            template_data: Some(HashMap::from([("name".into(), "Bob".into())])),
            ..Default::default()
        },
        BatchRecipient {
            email: "user2@example.com".into(),
            language: Some("fr".into()),
            template_data: Some(HashMap::from([("name".into(), "Carol".into())])),
            ..Default::default()
        },
    ],
    ..Default::default()
}).await?;

println!("Sent: {}, Failed: {}", resp.sent, resp.failed);
```

## Check Delivery Status

```rust
let status = client.get_email_status("email-uuid").await?;
println!("Status: {}", status.status);
```

## Error Handling

```rust
match client.get_email_status("invalid-uuid").await {
    Err(posta::Error::Api { status_code, info }) => {
        println!("Status: {status_code}");
        if let Some(info) = info {
            println!("Message: {}", info.message);
        }
    }
    Err(posta::Error::Http(e)) => println!("Network error: {e}"),
    Err(posta::Error::Decode(e)) => println!("Decode error: {e}"),
    Ok(status) => println!("Status: {}", status.status),
}
```

## Error Types

| Variant | Description |
|---------|-------------|
| `Error::Api { status_code, info }` | API returned a non-2xx response |
| `Error::Http(e)` | Network or connection error |
| `Error::Decode(e)` | Failed to parse response JSON |
