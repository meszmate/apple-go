# apple-go

`apple-go` is a unofficial Golang package to validate authorization tokens and manage the authorization of Apple Sign In server side. It provides utility functions and models to retrieve user information and validate authorization codes.

## Installation

Install with go modules:

```
go get github.com/meszmate/apple-go
```

## Usage

The package follow the Go approach to resolve problems, the usage is pretty straightforward, you start initiating a client with:

```go
package main

import (
    "github.com/meszmate/apple-go"
)

func main() {
    appleAuth, err := apple.New("<APP-ID>", "<TEAM-ID>", "<KEY-ID>", "/path/to/apple-sign-in-key.p8")
    if err != nil {
        panic(err)
    }
}
```

Using base64 env variable:

```go
package main

import (
    "log"
    "os"

    "github.com/meszmate/apple-go"
)

func main() {
    // Load the key from an environment variable encoded with:
    //   base64 -i AuthKey_ABCDE12345.p8
    //   export APPLE_KEY=<base64-string>
    auth, err := apple.NewB64(
        "com.example.app",        // App ID / Bundle ID
        "TEAM123456",             // Apple Team ID
        "KEY_ABCDE12345",         // Key ID
        os.Getenv("APPLE_KEY"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // use auth ...
    _ = auth
}
```

Generate the Apple Sign-In URL and send the user to it:

```go
package main

import (
    "fmt"
    "github.com/meszmate/apple-go"
)

func main() {
    cfg := apple.AuthorizeURLConfig{
        ClientID:     "com.example.app.login",        // Services ID
        RedirectURI:  "https://example.com/auth/apple/callback",
        Scope:        []string{"email", "name"},
        State:        "csrf-123",
        Nonce:        "nonce-abc",
    }

    loginURL := apple.AuthorizeURL(cfg)
    fmt.Println("Redirect the user to:", loginURL) 
    // output: Redirect the user to: https://appleid.apple.com/auth/authorize?response_type=code&response_mode=form_post&client_id=com.example.app.login&redirect_uri=https%3A%2F%2Fexample.com%2Fauth%2Fapple%2Fcallback&state=csrf-123&nonce=nonce-abc&scope=email+name
}
```

To validate an authorization code, retrieving refresh and access tokens:

```go
package main

import (
    "github.com/meszmate/apple-go"
)

func main() {
    appleAuth, err := apple.New("<APP-ID>", "<TEAM-ID>", "<KEY-ID>", "/path/to/apple-sign-in-key.p8")
    if err != nil {
        panic(err)
    }

    // Validate authorization code from a mobile app.
    tokenResponse, err := appleAuth.ValidateCode("<AUTHORIZATION-CODE>")
    if err != nil {
        panic(err)
    }

    // Validate authorization code from web app with redirect uri.
    tokenResponse, err := appleAuth.ValidateCodeWithRedirectURI("<AUTHORIZATION-CODE>", "https://redirect-uri")
    if err != nil {
        panic(err)
    }
}
```

The returned `tokenResponse` provides the access token, to make requests on behalf of the user with Apple servers, the refresh token, to retrieve a new access token after expiration, trought the `ValidateRefreshToken` method, and the id token, which is a JWT encoded string with user information. To retrieve the user information from this id token we provide a utility function `GetUserInfoFromIDToken`:

```go
package main

import (
    "fmt"

    "github.com/meszmate/apple-go"
)

func main() {
    appleAuth, err := apple.New("<APP-ID>", "<TEAM-ID>", "<KEY-ID>", "/path/to/apple-sign-in-key.p8")
    if err != nil {
        panic(err)
    }

    // Validate authorization code from a mobile app.
    tokenResponse, err := appleAuth.ValidateCode("<AUTHORIZATION-CODE>")
    if err != nil {
        panic(err)
    }

    user, err := apple.GetUserInfoFromIDToken(tokenResponse.idToken)
    if err != nil {
        panic(err)
    }

    // User Apple unique identification.
    fmt.Println(user.UID)
    // User email if the user provided it.
    fmt.Println(user.Email)
}
```
