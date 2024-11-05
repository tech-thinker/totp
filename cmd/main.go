package main

import (
	"fmt"

	"github.com/tech-thinker/totp"
)

func main() {
	secret, _ := totp.GenerateSecret()
    duration := 30
	code, err := totp.TOTP(secret, duration)
	if err != nil {
		fmt.Println("Error generating TOTP:", err)
		return
	}
	fmt.Println("Generated TOTP Code:", code)
    fmt.Println("Verify:", totp.Validate(secret, duration, code))
}

