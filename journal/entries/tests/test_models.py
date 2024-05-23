from datetime import timedelta

from django.utils import timezone

from journal.accounts.tests.factories import UserFactory
from journal.entries.models import Entry
from journal.entries.tests.factories import EntryFactory


class TestEntry:
    def test_factory(self):
        """A factory produces a valid journal entry."""
        entry = EntryFactory()

        assert entry.body != ""
        assert entry.when is not None
        assert entry.user is not None

    def test_get_random_for(self):
        """The manager can get a random entry."""
        today = timezone.localdate()
        user = UserFactory()
        entry_1 = EntryFactory(user=user, when=today)
        entry_2 = EntryFactory(user=user, when=today + timedelta(days=1))

        entry = Entry.objects.get_random_for(user)

        assert entry in {entry_1, entry_2}

    def test_get_random_for_only_user(self):
        """The random entry only belongs to the user."""
        entry_1 = EntryFactory()
        EntryFactory()

        entry = Entry.objects.get_random_for(entry_1.user)

        assert entry == entry_1

    def test_get_random_for_no_entries(self):
        """No entry is returned when there are no entries to pick from."""
        user = UserFactory()

        entry = Entry.objects.get_random_for(user)

        assert entry is None
