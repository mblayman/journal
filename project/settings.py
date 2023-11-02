from pathlib import Path

import dj_database_url
import environ

BASE_DIR = Path(__file__).resolve().parent.parent

env = environ.Env(
    ACCOUNT_DEFAULT_HTTP_PROTOCOL=(str, "https"),
    ALLOWED_HOSTS=(list, []),
    CSRF_COOKIE_SECURE=(bool, True),
    DATABASE_CONN_MAX_AGE=(int, 600),
    DATABASE_SSL_REQUIRE=(bool, True),
    DEBUG=(bool, False),
    EMAIL_BACKEND=(str, "anymail.backends.sendgrid.EmailBackend"),
    SECURE_HSTS_SECONDS=(int, 60 * 60 * 24 * 365),
    SECURE_SSL_REDIRECT=(bool, True),
    SENTRY_ENABLED=(bool, True),
    SESSION_COOKIE_SECURE=(bool, True),
)
environ.Env.read_env(BASE_DIR / ".env")

SECRET_KEY = env("SECRET_KEY")
DEBUG = env("DEBUG")
ALLOWED_HOSTS: list[str] = env("ALLOWED_HOSTS")

# Application definition

INSTALLED_APPS = [
    "django.contrib.admin",
    "django.contrib.auth",
    "django.contrib.contenttypes",
    "django.contrib.sessions",
    "django.contrib.messages",
    "django.contrib.sites",
    "django.contrib.staticfiles",
    "allauth",
    "allauth.account",
    # Needed by default templates even though we're not using a social provider.
    "allauth.socialaccount",
    "anymail",
    "django_extensions",
    "simple_history",
    "journal.accounts",
    "journal.core",
    "journal.entries",
    "journal.sentry",
]

MIDDLEWARE = [
    "django.middleware.security.SecurityMiddleware",
    "django.contrib.sessions.middleware.SessionMiddleware",
    "django.middleware.common.CommonMiddleware",
    "django.middleware.csrf.CsrfViewMiddleware",
    "django.contrib.auth.middleware.AuthenticationMiddleware",
    "django.contrib.messages.middleware.MessageMiddleware",
    "django.middleware.clickjacking.XFrameOptionsMiddleware",
]

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
    "default": dj_database_url.config(
        conn_max_age=env("DATABASE_CONN_MAX_AGE"),
        ssl_require=env("DATABASE_SSL_REQUIRE"),
    )
}

# Admin

ADMIN_URL_PATH_TOKEN = env("ADMIN_URL_PATH_TOKEN")

# Auth

AUTH_PASSWORD_VALIDATORS = [
    {
        "NAME": "django.contrib.auth.password_validation.UserAttributeSimilarityValidator",
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
    "allauth.account.auth_backends.AuthenticationBackend",
]
LOGIN_REDIRECT_URL = "/"
LOGOUT_REDIRECT_URL = "/"

# Email

EMAIL_BACKEND = env("EMAIL_BACKEND")
DEFAULT_FROM_FIELD = "noreply@journeyinbox.com"
SERVER_EMAIL = "noreply@journeyinbox.com"

# Internationalization
# https://docs.djangoproject.com/en/4.2/topics/i18n/

LANGUAGE_CODE = "en-us"
TIME_ZONE = "UTC"
USE_I18N = True
USE_TZ = True

# Security
CSRF_COOKIE_SECURE = env("CSRF_COOKIE_SECURE")
SECURE_HSTS_INCLUDE_SUBDOMAINS = True
SECURE_HSTS_PRELOAD = True
SECURE_HSTS_SECONDS = env("SECURE_HSTS_SECONDS")
SECURE_PROXY_SSL_HEADER = ("HTTP_X_FORWARDED_PROTO", "https")
SECURE_SSL_REDIRECT = env("SECURE_SSL_REDIRECT")
SESSION_COOKIE_SECURE = env("SESSION_COOKIE_SECURE")

# Static files (CSS, JavaScript, Images)
# https://docs.djangoproject.com/en/4.2/howto/static-files/

STATIC_ROOT = BASE_DIR / "staticfiles"
STATIC_URL = "static/"
# STATICFILES_DIRS = [BASE_DIR / "static"]

# Default primary key field type
# https://docs.djangoproject.com/en/4.2/ref/settings/#default-auto-field

DEFAULT_AUTO_FIELD = "django.db.models.BigAutoField"

# Sites

SITE_ID = 1

# django-allauth

# ACCOUNT_ADAPTER => default
# ACCOUNT_AUTHENTICATED_LOGIN_REDIRECTS => default
ACCOUNT_AUTHENTICATION_METHOD = "email"
ACCOUNT_CONFIRM_EMAIL_ON_GET = True
# ACCOUNT_EMAIL_CONFIRMATION_ANONYMOUS_REDIRECT_URL => default
# ACCOUNT_EMAIL_CONFIRMATION_AUTHENTICATED_REDIRECT_URL => default
# ACCOUNT_EMAIL_CONFIRMATION_EXPIRE_DAYS => default
# ACCOUNT_EMAIL_CONFIRMATION_HMAC => default
ACCOUNT_EMAIL_REQUIRED = True
ACCOUNT_EMAIL_VERIFICATION = "mandatory"
ACCOUNT_EMAIL_SUBJECT_PREFIX = "JourneyInbox - "
ACCOUNT_DEFAULT_HTTP_PROTOCOL = env("ACCOUNT_DEFAULT_HTTP_PROTOCOL")
# ACCOUNT_EMAIL_CONFIRMATION_COOLDOWN => default
# ACCOUNT_EMAIL_MAX_LENGTH => default
# ACCOUNT_MAX_EMAIL_ADDRESSES => default
# ACCOUNT_FORMS => default
# ACCOUNT_LOGIN_ATTEMPTS_LIMIT => default
# ACCOUNT_LOGIN_ATTEMPTS_TIMEOUT => default
ACCOUNT_LOGIN_ON_EMAIL_CONFIRMATION = True
# ACCOUNT_LOGOUT_ON_GET => default
# ACCOUNT_LOGOUT_ON_PASSWORD_CHANGE => default
# ACCOUNT_LOGIN_ON_PASSWORD_RESET => default
# ACCOUNT_LOGOUT_REDIRECT_URL => default
# ACCOUNT_PASSWORD_INPUT_RENDER_VALUE => default
ACCOUNT_PRESERVE_USERNAME_CASING = False
# ACCOUNT_PREVENT_ENUMERATION => default
# ACCOUNT_RATE_LIMITS => default
ACCOUNT_SESSION_REMEMBER = True
# ACCOUNT_SIGNUP_EMAIL_ENTER_TWICE => default
# ACCOUNT_SIGNUP_FORM_CLASS => default
ACCOUNT_SIGNUP_PASSWORD_ENTER_TWICE = False
# ACCOUNT_SIGNUP_REDIRECT_URL => default
# ACCOUNT_TEMPLATE_EXTENSION => default
# ACCOUNT_USERNAME_BLACKLIST => default
# ACCOUNT_UNIQUE_EMAIL => default
ACCOUNT_USER_DISPLAY = lambda user: user.email  # noqa
# ACCOUNT_USER_MODEL_EMAIL_FIELD => default
# ACCOUNT_USER_MODEL_USERNAME_FIELD => default
# ACCOUNT_USERNAME_MIN_LENGTH => default
ACCOUNT_USERNAME_REQUIRED = False
# ACCOUNT_USERNAME_VALIDATORS => default
# SOCIALACCOUNT_* => default

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

# sentry-sdk

SENTRY_ENABLED = env("SENTRY_ENABLED")
SENTRY_DSN = env("SENTRY_DSN")
