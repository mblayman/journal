# Convert to Go

## Motivation

I want to try something new. I use Go at work, but the ecosystem is set up.
This would give me a chance to start a Go project, from scratch,
that I can deploy myself.

Because the app is so simple when serving just me,
I think the Go version will use way less memory.
I could deploy this version on a shared droplet that I have access to
and save $8/month.

## Plan

### Foundation

1. Get a Go binary going in this project.
2. Get the binary baked into the Docker image early.
3. Deploy and run the Go binary along with the main process.
   Use `beta.journeyinbox.com` as a temporary domain name.
4. Integrate Sentry.

### MVP

Build out the following:

1. A simple index page.
2. The `up` route for health check purposes.
3. The webhook that can read and properly parse data from SendGrid.
4. The background job (goroutine?) that can send the prompt daily.

## Non-goals

* This is for other people.
* This site is styled to look pretty.
* I have to build everything from scratch.
