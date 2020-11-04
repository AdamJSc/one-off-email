# one-off-email

## About

Send a one-off email to a provided recipient list.

## Requirements

* Golang 1.14

## Running Locally

### Create config files

From the project root...

Create env file:

```
cp data/.env.example data/.env
```

Create message templates:

```
cp data/templates/example.message.html data/templates/message.html
cp data/templates/example.message.txt data/templates/message.txt
```

Change the template names of each new file respectively:

```
data/templates/message.html

-{{define "example_message_html"}}
+{{define "message_html"}}
```

```
data/templates/message.txt

-{{define "example_message_txt"}}
+{{define "message_txt"}}
```

### Run in preview mode

From the project root...

```
go run main.go
```

...then visit http://localhost:8080 in your web browser to preview the templates

### Run and send emails

From the project root...

```
go run main.go -send
```
