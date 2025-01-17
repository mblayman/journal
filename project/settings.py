from pathlib import Path

import environs

BASE_DIR = Path(__file__).resolve().parent.parent

env = environs.Env()
env.read_env()

SECRET_KEY = env("SECRET_KEY")
DEBUG = env.bool("DEBUG", False)
DEBUG_TOOLBAR = env.bool("DEBUG_TOOLBAR", False)
ALLOWED_HOSTS: list[str] = env.list("ALLOWED_HOSTS", [])

# Application definition

INSTALLED_APPS = [
    "django.contrib.admin",
    "django.contrib.auth",
    "django.contrib.contenttypes",
    "django.contrib.sessions",
    "django.contrib.messages",
    "django.contrib.sites",
    "django.contrib.staticfiles",
    "anymail",
    "django_extensions",
    "djstripe",
    "huey.contrib.djhuey",
    "simple_history",
    "waffle",
    "journal.accounts",
    "journal.core",
    "journal.entries",
    "journal.sentry",
]

MIDDLEWARE = [
    "django.middleware.security.SecurityMiddleware",
    "whitenoise.middleware.WhiteNoiseMiddleware",
    "django.contrib.sessions.middleware.SessionMiddleware",
    "django.middleware.common.CommonMiddleware",
    "django.middleware.csrf.CsrfViewMiddleware",
    "django.contrib.auth.middleware.AuthenticationMiddleware",
    "denied.middleware.DeniedMiddleware",
    "django.contrib.messages.middleware.MessageMiddleware",
    "django.middleware.clickjacking.XFrameOptionsMiddleware",
    "waffle.middleware.WaffleMiddleware",
]

# Enable the debug toolbar only in DEBUG mode.
if DEBUG and DEBUG_TOOLBAR:
    INSTALLED_APPS.append("debug_toolbar")
    MIDDLEWARE.insert(0, "debug_toolbar.middleware.DebugToolbarMiddleware")
    INTERNAL_IPS = ["127.0.0.1"]

ROOT_URLCONF = "project.urls"

TEMPLATES = [
    {
        "BACKEND": "django.template.backends.django.DjangoTemplates",
        "DIRS": [BASE_DIR / "templates"],
        "APP_DIRS": True,
        "OPTIONS": {
            "context_processors": [
                "django.template.context_processors.debug",
                "django.template.context_processors.request",
                "django.contrib.auth.context_processors.auth",
                "django.contrib.messages.context_processors.messages",
            ],
        },
    },
]

# As of Django 4.1, the cached loader is used in development mode.
# runserver works around this in some manner, but Gunicorn does not.
# Override the loaders to get non-cached behavior.
if DEBUG:
    # app_dirs isn't allowed to be True when the loaders key is present.
    TEMPLATES[0]["APP_DIRS"] = False
    TEMPLATES[0]["OPTIONS"]["loaders"] = [
        "django.template.loaders.filesystem.Loader",
        "django.template.loaders.app_directories.Loader",
    ]

WSGI_APPLICATION = "project.wsgi.application"

# Database
# https://docs.djangoproject.com/en/4.2/ref/settings/#databases

DATABASES = {
    "default": {
        "ENGINE": "django.db.backends.sqlite3",
        "NAME": env.path("DB_DIR", BASE_DIR) / "db.sqlite3",
        "OPTIONS": {
            "init_command": """
                PRAGMA journal_mode=wal;
                PRAGMA busy_timeout=5000;
                PRAGMA synchronous=normal;
                PRAGMA cache_size=-20000;
                PRAGMA temp_store=memory;
            """,
            "transaction_mode": "IMMEDIATE",
        },
    }
}

# Admin

ADMIN_URL_PATH_TOKEN = env("ADMIN_URL_PATH_TOKEN")

# Auth

AUTH_PASSWORD_VALIDATORS = [
    {
        "NAME": "django.contrib.auth.password_validation.UserAttributeSimilarityValidator",  # noqa
    },
    {
        "NAME": "django.contrib.auth.password_validation.MinimumLengthValidator",
    },
    {
        "NAME": "django.contrib.auth.password_validation.CommonPasswordValidator",
    },
    {
        "NAME": "django.contrib.auth.password_validation.NumericPasswordValidator",
    },
]
AUTH_USER_MODEL = "accounts.User"
AUTHENTICATION_BACKENDS = [
    "django.contrib.auth.backends.ModelBackend",
    "sesame.backends.ModelBackend",
]
LOGIN_REDIRECT_URL = "/"
LOGOUT_REDIRECT_URL = "/"

# Email

EMAIL_BACKEND = env.str("EMAIL_BACKEND", "anymail.backends.sendgrid.EmailBackend")
DEFAULT_FROM_EMAIL = "noreply@email.journeyinbox.com"
SERVER_EMAIL = "noreply@email.journeyinbox.com"

# Internationalization
# https://docs.djangoproject.com/en/4.2/topics/i18n/

LANGUAGE_CODE = "en-us"
TIME_ZONE = "UTC"
USE_I18N = True
USE_TZ = True

# Security
CSRF_COOKIE_SECURE = env.bool("CSRF_COOKIE_SECURE", True)
SECURE_HSTS_INCLUDE_SUBDOMAINS = True
SECURE_HSTS_PRELOAD = True
SECURE_HSTS_SECONDS = env.int("SECURE_HSTS_SECONDS", 60 * 60 * 24 * 365)
SECURE_PROXY_SSL_HEADER = ("HTTP_X_FORWARDED_PROTO", "https")
# The health check was failing with a 301 to HTTPS.
# With kamal-proxy in front and the .app domain, this should not be needed.
SECURE_SSL_REDIRECT = False
SESSION_COOKIE_SECURE = env.bool("SESSION_COOKIE_SECURE", True)

SILENCED_SYSTEM_CHECKS: list[str] = [
    # STRIPE_TEST_SECRET_KEY and STRIPE_LIVE_SECRET_KEY settings exist
    # and djstripe wants them not to exist.
    "djstripe.I002",
    # Disable warning about SECURE_SSL_REDIRECT.
    # The combo of kamal-proxy using Let's Encrypt and the `.app` domain
    # only working with HTTPS means that the warning can be ignored safely.
    "security.W008",
]

# Static files (CSS, JavaScript, Images)
# https://docs.djangoproject.com/en/4.2/howto/static-files/

STATIC_ROOT = BASE_DIR / "staticfiles"
STATIC_URL = "static/"
STATICFILES_DIRS = [BASE_DIR / "static"]
STORAGES = {
    "staticfiles": {
        "BACKEND": "whitenoise.storage.CompressedManifestStaticFilesStorage",
    },
}

# Default primary key field type
# https://docs.djangoproject.com/en/4.2/ref/settings/#default-auto-field

DEFAULT_AUTO_FIELD = "django.db.models.BigAutoField"

# Sites

SITE_ID = 1

# App settings

# Is the app in a secure context or not?
IS_SECURE = env.bool("IS_SECURE", True)

# dj-stripe

STRIPE_LIVE_SECRET_KEY = env("STRIPE_LIVE_SECRET_KEY")
STRIPE_TEST_SECRET_KEY = env("STRIPE_TEST_SECRET_KEY")
STRIPE_LIVE_MODE = env.bool("STRIPE_LIVE_MODE", True)
STRIPE_PUBLISHABLE_KEY = (
    env("STRIPE_LIVE_PUBLISHABLE_KEY")
    if STRIPE_LIVE_MODE
    else env("STRIPE_TEST_PUBLISHABLE_KEY")
)
DJSTRIPE_WEBHOOK_SECRET = env("DJSTRIPE_WEBHOOK_SECRET")
# This setting is recommended in the dj-stripe docs as the best default.
DJSTRIPE_FOREIGN_KEY_TO_FIELD = "id"
# This setting is recommended in the dj-stripe docs as the best default.
DJSTRIPE_USE_NATIVE_JSONFIELD = True
PRICE_LOOKUP_KEY = "monthly-v1"

# django-anymail

ANYMAIL = {
    "SENDGRID_API_KEY": env("SENDGRID_API_KEY"),
    "WEBHOOK_SECRET": env("ANYMAIL_WEBHOOK_SECRET"),
}

# django-extensions

GRAPH_MODELS = {
    "app_labels": ["accounts", "core", "entries"],
    "rankdir": "BT",
    "output": "models.png",
}

# django-hashid-field

HASHID_FIELD_SALT = env("HASHID_FIELD_SALT")

# django-sesame

SESAME_TOKEN_NAME = "token"  # noqa S105
SESAME_MAX_AGE = 60 * 60  # 1 hour
# If JourneyInbox allows email changes in the future,
# we may want to change this default.
# SESAME_INVALIDATE_ON_EMAIL_CHANGE = False

# django-waffle

WAFFLE_CREATE_MISSING_FLAGS = True

# Huey
HUEY = {
    "huey_class": "huey.SqliteHuey",
    "filename": env.path("DB_DIR", BASE_DIR) / "huey.sqlite3",
    "immediate": False,
}

# sentry-sdk

SENTRY_ENABLED = env.bool("SENTRY_ENABLED", True)
SENTRY_DSN = env("SENTRY_DSN")
