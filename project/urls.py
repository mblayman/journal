from denied.decorators import allow
from django.conf import settings
from django.contrib import admin
from django.contrib.auth.views import LogoutView
from django.urls import include, path
from sesame.views import LoginView

from journal.accounts.views import (
    account_settings,
    success,
)
from journal.core.views import index, up
from journal.entries.views import export_entries, import_entries

urlpatterns = [
    path("", index, name="index"),
    path("up", up, name="up"),
    #
    # Accounts
    #
    path("login", LoginView.as_view(), name="sesame-login"),
    path("logout", allow(LogoutView.as_view()), name="logout"),
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
