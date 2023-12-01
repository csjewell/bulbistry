module github.com/csjewell/bulbistry

go 1.21.4

require (
	internal/config v0.0.0
	internal/database v0.0.0
	internal/urls v0.0.0
	internal/version v0.0.0
)

replace internal/config => ./internal/config
replace internal/database => ./internal/database
replace internal/urls => ./internal/urls
replace internal/version => ./internal/version

require (
	github.com/goccy/go-yaml v1.11.2
	github.com/google/uuid v1.4.0
	github.com/gorilla/mux v1.8.1
	github.com/mattn/go-sqlite3 v1.14.18
	github.com/tg123/go-htpasswd v1.2.1
)

require (
	github.com/GehirnInc/crypt v0.0.0-20230320061759-8cc1b52080c5 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/manifoldco/promptui v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/urfave/cli/v2 v2.25.7 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/crypto v0.16.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
)
