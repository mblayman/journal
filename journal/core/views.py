from denied.decorators import allow
from django.http import HttpRequest, HttpResponse
from django.shortcuts import render

from journal.accounts import constants


@allow
def about(request: HttpRequest) -> HttpResponse:
    """The about page... duh"""
    context = {}
    return render(request, "core/about.html", context)


@allow
def faq(request: HttpRequest) -> HttpResponse:
    """Frequently Asked Questions"""
    context = {}
    return render(request, "core/faq.html", context)


@allow
def index(request: HttpRequest) -> HttpResponse:
    """The entry point for the website."""
    context = {"trial_days": constants.TRIAL_DAYS}
    template_name = "core/index_unauthenticated.html"
    if request.user.is_authenticated:
        template_name = "core/index.html"
    return render(request, template_name, context)


@allow
def terms(request: HttpRequest) -> HttpResponse:
    """The terms of service"""
    context = {}
    return render(request, "core/terms.html", context)


@allow
def privacy(request: HttpRequest) -> HttpResponse:
    """The privacy policy"""
    context = {}
    return render(request, "core/privacy.html", context)
