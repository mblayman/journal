import json

from dateutil.parser import parse
from denied.authorizers import any_authorized, staff_authorized
from denied.decorators import authorize
from django.http import HttpRequest, HttpResponse, JsonResponse
from django.utils import timezone
from django.views.decorators.csrf import csrf_exempt

from .models import Entry


@csrf_exempt
@authorize(staff_authorized)
def import_entries(request: HttpRequest) -> HttpResponse:
    """Create new entries from JSON data.

    Issue #126: validate JSON
    Issue #126: validation that the entries are correct (i.e., just `body` and `when`)
    Issue #126: add validation on the number of allowed entries. 365 * 100
    Issue #126: fail on any date parse errors
    Issue #126: limit the allowed length of each post.
    Issue #126: maybe change HTTP status code?
    """
    data = json.loads(request.body)

    entries = [
        Entry(body=entry["body"], when=parse(entry["when"]), user=request.user)
        for entry in data
    ]
    Entry.objects.bulk_create(entries)
    return HttpResponse(b"ok")


@authorize(any_authorized)
def export_entries(request: HttpRequest) -> HttpResponse:
    """Export all of a user's entries as a JSON file."""
    entries = list(
        Entry.objects.filter(user=request.user).order_by("when").values("body", "when")
    )
    # safe=False is ok here because we are passing a list. JsonResponse only
    # wants dictionaries, but we know for certain the JSON spec allows for lists.
    response = JsonResponse(data=entries, safe=False)
    today = timezone.localdate()
    response["Content-Disposition"] = (
        f'attachment; filename="journeyinbox-{today}.json"'
    )
    return response
