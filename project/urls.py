from django.conf import settings
from django.contrib import admin
from django.urls import include, path

from journal.core.views import index, privacy, terms

urlpatterns = [
    path("", index, name="index"),
    path("privacy/", privacy, name="privacy"),
    path("terms/", terms, name="terms"),
    path("accounts/", include("allauth.urls")),
    path(f"{settings.ADMIN_URL_PATH_TOKEN}/admin/", admin.site.urls),
    path("anymail/", include("anymail.urls")),
]
