from denied.authorizers import any_authorized
from denied.decorators import authorize
from django.shortcuts import render


@authorize(any_authorized)
def account_settings(request):
    """Show the user's settings."""
    return render(request, "accounts/settings.html", {})


@authorize(any_authorized)
def success(request):
    """Show the success after account activation."""
    return render(request, "accounts/success.html", {})
