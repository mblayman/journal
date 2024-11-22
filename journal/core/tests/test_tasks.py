from journal.core.tasks import clear_db_sessions


class TestClearDBSessions:
    def test_ok(self):
        """Sanity check"""
        clear_db_sessions()
