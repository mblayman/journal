from denied.decorators import allow
from django.conf import settings
from django.contrib import admin
from django.urls import include, path
from sesame.views import LoginView

from journal.accounts.views import (
    account_settings,
    check_email,
    create_billing_portal_session,
    create_checkout_session,
    success,
)
from journal.core.views import about, faq, index, privacy, terms, up
from journal.entries.views import export_entries, import_entries

urlpatterns = [
    path("", index, name="index"),
    path("about/", about, name="about"),
    path("faq/", faq, name="faq"),
    path("up", up, name="up"),
    #
    # Boring legal stuff
    #
    path("privacy/", privacy, name="privacy"),
    path("terms/", terms, name="terms"),
    #
    # Accounts
    #
    path("check-email", check_email, name="check-email"),
    path("login", LoginView.as_view(), name="sesame-login"),
    path(
        "accounts/create-checkout-session/",
        create_checkout_session,
        name="create_checkout_session",
    ),
    path(
        "accounts/stripe-billing-portal/",
        create_billing_portal_session,
        name="create_billing_portal_session",
    ),
    path("settings/", account_settings, name="settings"),
    path("success/", success, name="success"),
    #
    # Entries
    #
    path("import/", import_entries, name="import_entries"),
    path("export/", export_entries, name="export_entries"),
    #
    # Third party routes
    #
    path("accounts/", allow(include("allauth.urls"))),
    path("anymail/", allow(include("anymail.urls"))),
    path("stripe/", allow(include("djstripe.urls", namespace="djstripe"))),
    #
    # Admin
    #
    path(f"{settings.ADMIN_URL_PATH_TOKEN}/admin/", allow(admin.site.urls)),
]

# Enable the debug toolbar only in DEBUG mode.
if settings.DEBUG and settings.DEBUG_TOOLBAR:
    urlpatterns = [
        path("__debug__/", allow(include("debug_toolbar.urls")))
    ] + urlpatterns
