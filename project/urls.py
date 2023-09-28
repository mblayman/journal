from django.conf import settings
from django.conf.urls.static import static
from django.contrib import admin
from django.urls import include, path

from journal.core.views import index

urlpatterns = [
    path("", index, name="index"),
    path("accounts/", include("allauth.urls")),
    path(f"{settings.ADMIN_URL_PATH_TOKEN}/admin/", admin.site.urls),
    path("anymail/", include("anymail.urls")),
] + static(settings.STATIC_URL, document_root=settings.STATIC_ROOT)
