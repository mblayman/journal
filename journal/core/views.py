from django.http import HttpRequest, HttpResponse
from django.shortcuts import render


def index(request: HttpRequest) -> HttpResponse:
    """The entry point for the website."""
    for header in request.headers.items():
        print(header)
    context = {}
    return render(request, "core/index.html", context)
