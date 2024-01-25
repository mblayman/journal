from django.http import HttpRequest, HttpResponse
from django.shortcuts import render

from journal.payments.gateway import PaymentsGateway


def index(request: HttpRequest) -> HttpResponse:
    """The entry point for the website."""
    payments_gateway = PaymentsGateway()
    context = {
        "payments_publishable_key": payments_gateway.publishable_key,
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
