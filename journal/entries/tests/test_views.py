from datetime import date, timedelta

from django.urls import reverse
from django.utils import timezone

from journal.accounts.tests.factories import UserFactory
from journal.entries.models import Entry
from journal.entries.tests.factories import EntryFactory


class TestImportEntries:
    def test_ok(self, admin_client, admin_user):
        """The import view creates entries."""
        data = [{"body": "hello from my journal", "when": "2024-04-17"}]

        response = admin_client.post(
            reverse("import_entries"),
            data=data,
            content_type="application/json",
        )

        assert response.status_code == 200
        assert Entry.objects.filter(user=admin_user).count() == 1
        entry = Entry.objects.first()
        assert entry.body == "hello from my journal"
        assert entry.when == date(2024, 4, 17)


class TestExportEntries:
    def test_unauthenticated(self, client):
        """Only allow authenticated users."""

        response = client.get(reverse("export_entries"))

        assert response.status_code == 302

    def test_other_user(self, client):
        """A user does not get another's entries."""
        another_user = UserFactory()
        EntryFactory(user=another_user)
        user = UserFactory()
        client.force_login(user)

        response = client.get(reverse("export_entries"))

        data = response.json()
        assert len(data) == 0

    def test_ok(self, client):
        today = timezone.localdate()
        user = UserFactory()
        entry_tomorrow = EntryFactory(user=user, when=today + timedelta(days=1))
        entry_today = EntryFactory(user=user, when=today)
        client.force_login(user)

        response = client.get(reverse("export_entries"))

        assert response.status_code == 200
        data = response.json()
        assert len(data[0].keys()) == 2
        assert data[0]["body"] == entry_today.body
        assert data[1]["body"] == entry_tomorrow.body
