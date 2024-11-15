FROM node:22 AS nodejs

WORKDIR /app

COPY frontend/package.json frontend/package-lock.json ./

RUN --mount=type=cache,target=/root/.npm \
    npm install --loglevel verbose

COPY frontend frontend/
COPY templates templates/

RUN npm --prefix frontend run build

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

COPY --from=nodejs /app/static/site.css static/

# Some configuration is needed to make Django happy, but these values have no
# impact to collectstatic so we can use dummy values.
RUN \
    ADMIN_URL_PATH_TOKEN=fake-token \
    ANYMAIL_WEBHOOK_SECRET=a-secret-to-everybody \
    DJSTRIPE_WEBHOOK_SECRET=whsec_asecrettoeverybody \
    HASHID_FIELD_SALT=a-secret-to-everybody \
    SECRET_KEY=a-secret-to-everybody \
    SENDGRID_API_KEY=a-secret-to-everybody \
    SENTRY_DSN=dsn_example \
    SENTRY_ENABLED=off \
    STRIPE_LIVE_PUBLISHABLE_KEY=pk_live_a-secret-to-everybody \
    STRIPE_LIVE_SECRET_KEY=sk_live_a-secret-to-everybody \
    STRIPE_TEST_PUBLISHABLE_KEY=pk_test_a-secret-to-everybody \
    STRIPE_TEST_SECRET_KEY=sk_test_a-secret-to-everybody \
    python manage.py collectstatic --noinput

USER app

ENTRYPOINT ["/app/bin/docker-entrypoint"]
EXPOSE 8000
CMD ["/app/bin/server"]
