# go-mailketing

[![Go Reference](https://pkg.go.dev/badge/github.com/wiradeltaid/go-mailketing.svg)](https://pkg.go.dev/github.com/wiradeltaid/go-mailketing)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Idiomatic Go client library for the [Mailketing API v2](https://mailketing.co.id).

Developed and maintained by **Wira Delta Indonesia**.

## Installation

```bash
go get github.com/wiradeltaid/go-mailketing
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/wiradeltaid/go-mailketing"
)

func main() {
	client := mailketing.NewClient("YOUR_API_TOKEN",
		mailketing.WithDefaultSender("Wira Delta Indonesia", "notification@wiradelta.id"),
	)

	resp, err := client.Send(context.Background(), mailketing.SendEmailRequest{
		Recipient: "user@example.com",
		Subject:   "Selamat Datang di Wira Delta Indonesia",
		Content:   "<h1>Halo!</h1><p>Terima kasih telah bergabung.</p>",
	})
	if err != nil {
		log.Fatalf("Gagal mengirim email: %v", err)
	}

	fmt.Printf("Email berhasil diantrekan! Message ID: %s\n", resp.Data.MessageID)
}
```

## Features

- **Standard Library Only:** Zero external dependencies (only Go standard library `net/http`, `context`, `encoding/json`).
- **Context-Aware:** Full support for timeouts, deadlines, and graceful cancellations.
- **Transactional Email:** Send HTML or plain text emails with attachments and custom tracking IDs.
- **Account Endpoints:** Inspect remaining credits (`GetCredits`) and list verified domains/senders (`GetSenders`).
- **Configurable:** Functional options for custom HTTP client, base URL, and default sender metadata.

## License

MIT License. See [LICENSE](LICENSE) for details.
