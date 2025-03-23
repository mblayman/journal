from denied.decorators import allow
from django.http import HttpRequest, HttpResponse
from django.shortcuts import redirect, render, reverse

from journal.accounts import constants
from journal.accounts.forms import SiginForm


@allow
def index(request: HttpRequest) -> HttpResponse:
    """The entry point for the website."""
    context: dict = {"trial_days": constants.TRIAL_DAYS}
    template_name = "core/index_unauthenticated.html"

    form = SiginForm()
    if request.method == "POST":
        form = SiginForm(request.POST)
        if form.is_valid():
            form.save()
            return redirect(reverse("index"))

    if request.user.is_authenticated:
        template_name = "core/index.html"

    context["form"] = form
    return render(request, template_name, context)


@allow
def up(request):
    """A healthcheck to show when the app is up and able to respond to requests."""
    return render(request, "core/up.html", {})
