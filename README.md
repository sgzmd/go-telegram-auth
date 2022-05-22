# go-telegram-auth: web authentication with Telegram

This library allows for simple web page authentication with Telegram. 

Documentation is not ready yet, but please check [this example](https://github.com/sgzmd/tgauth) which 
was the starting point for this library.


```mermaid
sequenceDiagram
    User ->>+IndexPage: Navigates
    IndexPage ->>+CheckAuth: Check auth status
    CheckAuth ->>CheckAuth: Verify if cookie is present
    Alt Cookie is valid
    CheckAuth -->>-IndexPage: Cookie Valid
    IndexPage -->>-User: Welcome, username
    else No valid cookie
    CheckAuth -->>LoginPage: redirect to /login
    LoginPage -->>LoginPage: Set Cookie
    LoginPage -->>IndexPage: Redirect
    end
```