from django.http import HttpRequest, HttpResponse
from django.shortcuts import render


def index(request: HttpRequest) -> HttpResponse:
    """The entry point for the website."""
    template_name = "core/index.html"
    return render(request, template_name, {})


def up(request):
    """A healthcheck to show when the app is up and able to respond to requests."""
    return render(request, "core/up.html", {})
