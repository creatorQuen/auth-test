# auth-test

To run the application, you have to define the environment variables and launch cmd/main.go.
In the project used migration "go-migrate" and "go-bindata".

- HOST_DB           `[IP address of the database]`
- LISTEN_PORT       `[Port of the machine]`
- NAME_DB           `[Name of the database]`
- PASSWORD_DB       `[Database password]`
- PORT_DB           `[Port of the database]`
- USER_DB           `[Database username]`
- SALT              `[Salt for hash password]`
- SIGN_KEY          `[Signing key for jwt]`
- TOKEN_TIME_LEFT   `[Token time expiration for jwt]`
