# Genesis Software Engineering School 2023 API (gses-2023)

## Description

This repository contains the code for a simple API developed for the Genesis Software Engineering School 2023.

## Installation

1. Clone the repository: `git clone https://github.com/vadimpk/gses-2023.git`
2. Navigate to the project directory: `cd gses-2023`

## Usage

### Prerequisites

1. Create `local/` directory in the `core/` directory of the project. This is where the database file will be stored.
2. Get API key from [CoinAPI](https://www.coinapi.io/). Or you can use my API key, which you can find in `config/config.go`.
3. Sign up for [MailGun](https://www.mailgun.com/) account. You will need to verify your domain and get your API key.
4. Fill `.env` files in both `core/` and `crypto/` directories with the following variables:

```
GSES_COIN_API_KEY=<your_coin_api_key>

GSES_MAILGUN_DOMAIN=<your_mailgun_domain>
GSES_MAILGUN_API_KEY=<your_mailgun_api_key>
GSES_MAILGUN_FROM=<your_mailgun_from_email>
```

### Run

1. Run `docker-compose up` to start the application locally.

### List of endpoints:

- `:8081/api/rate` (GET): get current bitcoin rate in UAH
- `:8080/api/subscribe` (POST): subscribe to mailing list
- `:8080/api/sendEmails` (POST): send emails with current currrency rate to all subscribers

## Architecture

<img width="1087" alt="architecture" src="https://github.com/GenesisEducationKyiv/main-project-vadimpk/assets/65962115/3f8f629d-0f56-463c-a0c6-8d4f4b1213aa">

