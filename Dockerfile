FROM golang:1.24.1-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM python:3.12-slim

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    UV_PROJECT_ENVIRONMENT=/usr/local

COPY --from=ghcr.io/astral-sh/uv:0.5.2 /uv /bin/uv

WORKDIR /app

RUN addgroup --gid 222 --system app \
    && adduser --uid 222 --system --group app

RUN mkdir -p /app && chown app:app /app

COPY --chown=app:app pyproject.toml uv.lock /app/

RUN --mount=type=cache,target=/root/.cache/uv \
    uv sync --frozen --no-dev

COPY --chown=app:app . /app/

# Some configuration is needed to make Django happy, but these values have no
# impact to collectstatic so we can use dummy values.
RUN \
    ADMIN_URL_PATH_TOKEN=fake-token \
    ANYMAIL_WEBHOOK_SECRET=a-secret-to-everybody \
    HASHID_FIELD_SALT=a-secret-to-everybody \
    SECRET_KEY=a-secret-to-everybody \
    SENDGRID_API_KEY=a-secret-to-everybody \
    SENTRY_DSN=dsn_example \
    SENTRY_ENABLED=off \
    python manage.py collectstatic --noinput

USER app

COPY --from=builder --chown=app:app /app/app .

ENTRYPOINT ["/app/bin/docker-entrypoint"]
EXPOSE 8000
CMD ["/app/bin/server"]
