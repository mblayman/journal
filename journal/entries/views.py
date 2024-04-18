import json

from dateutil.parser import parse
from django.contrib.auth.decorators import login_required, user_passes_test
from django.http import HttpRequest, HttpResponse

from .models import Entry


@login_required
@user_passes_test(lambda user: user.is_staff)
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
