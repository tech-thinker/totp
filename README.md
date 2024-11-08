# TOTP

TOTP is a Time Based One Time Password Algorithm, which can be used in client side or server side applications.

## Features
- Generate TOTP secret.
- Generate TOTP code.
- Verify TOTP code.

## How to use it?
- Install dependencies
```sh
go get github.com/tech-thinker/totp@v1.0.0
```
- Write your code
```go
package main    You, 05/11/24 09:52 • initial commit

import (
    "fmt"

    "github.com/tech-thinker/totp"
)

func main() {
    secret, err := totp.GenerateSecret()
    if err != nil {
        fmt.Println("Error generating secret:", err)
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
```
