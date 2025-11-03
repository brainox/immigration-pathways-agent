# Deploying to Heroku

This project is a Go web server (entrypoint: `cmd/server`) that listens on the port provided by the `PORT` environment variable. The repository includes a `Procfile` to run the compiled binary created by the Heroku Go buildpack.

Minimal steps to deploy to Heroku:

1. Install the Heroku CLI: https://devcenter.heroku.com/articles/heroku-cli

2. Login and create the app:

   heroku login
   heroku create <your-app-name>

3. Set the required config vars (at minimum the Gemini API key):

   heroku config:set GEMINI_API_KEY=your_real_key_here

4. Push your code to Heroku (assuming `main` branch):

   git push heroku main

   The Heroku Go buildpack will run `go install` on your packages and place the compiled binary in `bin/server`. The `Procfile` runs `bin/server`.

5. Scale and view logs:

   heroku ps:scale web=1
   heroku logs --tail

Local testing notes:

- You can run locally using `PORT=8080 go run ./cmd/server` or create a `.env` based on `.env.example` and use a tool like `direnv` or `env $(cat .env) go run ./cmd/server`.
- Ensure `GEMINI_API_KEY` is set locally for the Gemini client to work.

Troubleshooting:

- If the app crashes on startup on Heroku, check `heroku logs --tail` for the error. Common issues: missing `GEMINI_API_KEY`, port binding errors (the app must listen on $PORT), or build failures due to incompatible Go version.
- If you need a specific Go version, set `go 1.xx` in `go.mod` (this repo uses go 1.21).
