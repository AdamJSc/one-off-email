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

Amend env values so that they reflect the following:

* `MAILGUN_API_KEY` = API key registered within Mailgun
* `MAILGUN_SENDER_DOMAIN` = Sender domain configured within Mailgun
* `SENDER_NAME` = Sender's name to issue the email from
* `SENDER_EMAIL` = Sender's email address to issue the email from
* `REPLY_TO_NAME` = Reply-to name to include when issuing email
* `REPLY_TO_EMAIL` = Reply-to email address to include when issuing email
* `EMAIL_SUBJECT` = Subject to use when sending emails️
* `MESSAGE_SIGN_OFF` = Sign-off name within email message

Create message templates:

```
cp data/templates/example.message.html data/templates/message.html
cp data/templates/example.message.txt data/templates/message.txt
```

Change the template names of each new file respectively:

```diff
data/templates/message.html

-{{define "example_message_html"}}
+{{define "message_html"}}
```

```diff
data/templates/message.txt

-{{define "example_message_txt"}}
+{{define "message_txt"}}
```

Create recipients file:

```
cp data/recipients.example.yml data/recipients.yml
```

Amend new recipients file, ensuring that contents are valid YAML and retain the same `.name` and `.email` format per recipient.

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
