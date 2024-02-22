from django.conf import settings
from django.contrib import admin
from django.urls import include, path

from journal.accounts.views import create_checkout_session, success
from journal.core.views import index, privacy, terms

urlpatterns = [
    path("", index, name="index"),
    path("privacy/", privacy, name="privacy"),
    path("terms/", terms, name="terms"),
    path("accounts/", include("allauth.urls")),
    path("success/", success, name="success"),
    path(
        "accounts/create-checkout-session/",
        create_checkout_session,
        name="create_checkout_session",
    ),
    path(f"{settings.ADMIN_URL_PATH_TOKEN}/admin/", admin.site.urls),
    path("anymail/", include("anymail.urls")),
    path("stripe/", include("djstripe.urls", namespace="djstripe")),
]
