from denied.decorators import allow
from django.http import HttpRequest, HttpResponse
from django.shortcuts import render


@allow
def index(request: HttpRequest) -> HttpResponse:
    """The entry point for the website."""
    template_name = "core/index.html"
    return render(request, template_name, {})


@allow
def up(request):
    """A healthcheck to show when the app is up and able to respond to requests."""
    return render(request, "core/up.html", {})
