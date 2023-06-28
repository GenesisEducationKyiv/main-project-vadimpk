# Genesis Software Engineering School 2023 API (gses-2023)

## Description

This repository contains the code for a simple API developed for the Genesis Software Engineering School 2023.

## Installation

1. Clone the repository: `git clone https://github.com/vadimpk/gses-2023.git`
2. Navigate to the project directory: `cd gses-2023`

## Usage

### Prerequisites

1. Create `local/` directory in the root directory of the project. This is where the database file will be stored.
2. Get API key from [CoinAPI](https://www.coinapi.io/). Or you can use my API key, which you can find in `config/config.go`.
3. Sign up for [MailGun](https://www.mailgun.com/) account. You will need to verify your domain and get your API key.
4. Fill `.env` file in the root directory of the project with the following variables:

```
GSES_COIN_API_KEY=<your_coin_api_key>
GSES_MAILGUN_DOMAIN=<your_mailgun_domain>
GSES_MAILGUN_API_KEY=<your_mailgun_api_key>
GSES_MAILGUN_FROM=<your_mailgun_from_email>
```

### Run locally

1. Run `make run` to start the server locally.

### Run in Docker

1. Run `make docker-build` to build the docker image.
2. Run `make docker-run` to start the server in docker container.

### List of endpoints:

- `/api/rate` (GET): get current bitcoin rate in UAH
- `/api/subscribe` (POST): subscribe to mailing list
- `/api/sendEmails` (POST): send emails with current currrency rate to all subscribers

## Architecture

![image](https://github.com/vadimpk/gses-2023/assets/65962115/3a45eac1-c43a-4a48-82af-2cfafd3f603e)
