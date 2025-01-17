# Sign up and log in

- Single stream sign up or log in
- No password
- Log in via magic links (expiration on the links)
- Extend the session to 365 days

## Benefits

- No passwords... no password reset flows
- Not multiple places to get into the site.

## Drawbacks

- A lot more emails
- Bots could be a problem

## Implementation

- https://github.com/aaugustin/django-sesame for the magic links
- custom signin page
    - If a user doesn't exist, create one
    - Magic a magic link for the user
    - Email the link
    - Redirect to a "check your email" page
- Remove allauth
- Account migration... invalidate existing passwords
