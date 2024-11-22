from .settings import *  # noqa

# An in-memory database should be good enough for now.
DATABASES = {
    "default": {
        "ENGINE": "django.db.backends.sqlite3",
        "NAME": ":memory:",
    }
}

EMAIL_BACKEND = "django.core.mail.backends.locmem.EmailBackend"

HUEY = {
    "huey_class": "huey.SqliteHuey",
    "filename": ":memory:",
    "immediate": True,
}

STORAGES = {
    "staticfiles": {
        "BACKEND": "django.contrib.staticfiles.storage.StaticFilesStorage",
    },
}

MEDIA_URL = "/test_media"

# This eliminates the warning about a missing staticfiles directory.
WHITENOISE_AUTOREFRESH = True
