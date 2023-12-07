from .settings import *  # noqa

# An in-memory database should be good enough for now.
DATABASES = {
    "default": {
        "ENGINE": "django.db.backends.sqlite3",
        "NAME": ":memory:",
    }
}

EMAIL_BACKEND = "django.core.mail.backends.locmem.EmailBackend"

STORAGES = {
    "staticfiles": {
        "BACKEND": "django.contrib.staticfiles.storage.StaticFilesStorage",
    },
}

# This eliminates the warning about a missing staticfiles directory.
WHITENOISE_AUTOREFRESH = True
