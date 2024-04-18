from datetime import date

from django.urls import reverse

from journal.entries.models import Entry


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
