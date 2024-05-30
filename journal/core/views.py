from django.http import HttpRequest, HttpResponse
from django.shortcuts import render

from journal.accounts import constants
from journal.payments.gateway import PaymentsGateway


def about(request: HttpRequest) -> HttpResponse:
    """The about page... duh"""
    context = {}
    return render(request, "core/about.html", context)


def faq(request: HttpRequest) -> HttpResponse:
    """Frequently Asked Questions"""
    context = {}
    return render(request, "core/faq.html", context)


def index(request: HttpRequest) -> HttpResponse:
    """The entry point for the website."""
    payments_gateway = PaymentsGateway()
    context = {
        "payments_publishable_key": payments_gateway.publishable_key,
        "price": payments_gateway.price,
        "trial_days": constants.TRIAL_DAYS,
    }
    template_name = "core/index_unauthenticated.html"
    if request.user.is_authenticated:
        template_name = "core/index.html"
    return render(request, template_name, context)


def terms(request: HttpRequest) -> HttpResponse:
    """The terms of service"""
    context = {}
    return render(request, "core/terms.html", context)


def privacy(request: HttpRequest) -> HttpResponse:
    """The privacy policy"""
    context = {}
    return render(request, "core/privacy.html", context)
