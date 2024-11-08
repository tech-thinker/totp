package main

import (
	"fmt"

	"github.com/tech-thinker/totp"
)

func main() {
	secret, err := totp.GenerateSecret()
    if err != nil {
        fmt.Println("Error generating secret:", err)
        return
    }
    duration := 30
	otp, err := totp.TOTP(secret, duration)
	if err != nil {
		fmt.Println("Error generating TOTP:", err)
		return
	}
	fmt.Println("Generated TOTP Code:", otp)
    fmt.Println("Verify:", totp.Validate(secret, duration, otp))
}

