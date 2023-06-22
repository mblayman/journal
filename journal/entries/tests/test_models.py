from journal.entries.tests.factories import EntryFactory


class TestEntry:
    def test_factory(self):
        """A factory produces a valid journal entry."""
        entry = EntryFactory()

        assert entry.body != ""
        assert entry.when is not None
        assert entry.user is not None
