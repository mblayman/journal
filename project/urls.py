from django.conf import settings
from django.contrib import admin
from django.urls import include, path

from journal.accounts.views import create_checkout_session, success
from journal.core.views import index, privacy, terms
from journal.entries.views import import_entries

urlpatterns = [
    path("", index, name="index"),
    #
    # Boring legal stuff
    #
    path("privacy/", privacy, name="privacy"),
    path("terms/", terms, name="terms"),
    #
    # Accounts
    #
    path(
        "accounts/create-checkout-session/",
        create_checkout_session,
        name="create_checkout_session",
    ),
    path("success/", success, name="success"),
    #
    # Entries
    #
    path("import/", import_entries, name="import_entries"),
    #
    # Third party routes
    #
    path("accounts/", include("allauth.urls")),
    path("anymail/", include("anymail.urls")),
    path("stripe/", include("djstripe.urls", namespace="djstripe")),
    #
    # Admin
    #
    path(f"{settings.ADMIN_URL_PATH_TOKEN}/admin/", admin.site.urls),
]
