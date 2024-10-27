# db-email
Experiment to compare the time execution of sending emails with and without applying concurrency programming in Go. Stable internet connection is required to run this project.

## Stacks
- [Go](https://go.dev/)
- [go-mail](https://github.com/go-gomail/gomail)
- [Gorm](https://gorm.io/)
- [Viper](https://github.com/spf13/viper)
- PostgreSQL
- Docker

## Installation an running
1. Configure smtp host, port, and other credentials from .env
   ```
   cp .env.example .env
   ```
   Then open .env to:
   - configure ```EMAIL_DEFAULT_SENDER``` to your GMail account
   - configure ```EMAIL_DEFAULT_RECIPIENT``` to default email destination
   - configure ```EMAIL_APP_PASSWORD``` from password that you should already get from [Google App Password](https://myaccount.google.com/apppasswords)
   - configure ```NUM_EMAIL_SEEDS``` to set how many emails you want to sent to ```EMAIL_DEFAULT_RECIPIENT```. 
2. Run the app with docker. Make sure you have already install Docker and docker compose in your system. 
   ```
   docker compose up --build
   ```
3. Wait for a few seconds. This process will takes more seconds or even minutes depending on your internet connection and ```NUM_EMAIL_SEEDS``` that you have been set.  
   ```
   INFO[0057] successfully send 10 emails and failed to send 0 (total: 10) in 56.745299 seconds without concurrency
   INFO[0071] successfully send 10 emails and failed to send 0 (total: 10) in 13.746935 seconds with concurrency
   ```
