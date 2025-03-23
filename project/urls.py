from denied.decorators import allow
from django.conf import settings
from django.contrib import admin
from django.urls import include, path

from journal.core.views import index, up

urlpatterns = [
    path("", index, name="index"),
    path("up", up, name="up"),
    path("anymail/", allow(include("anymail.urls"))),
    path("stripe/", allow(include("djstripe.urls", namespace="djstripe"))),
    path(f"{settings.ADMIN_URL_PATH_TOKEN}/admin/", allow(admin.site.urls)),
]

# Enable the debug toolbar only in DEBUG mode.
if settings.DEBUG and settings.DEBUG_TOOLBAR:
    urlpatterns = [
        path("__debug__/", allow(include("debug_toolbar.urls")))
    ] + urlpatterns
