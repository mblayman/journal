from allauth.account.adapter import DefaultAccountAdapter

from .constants import MAX_TRIALING_USERS
from .models import Account


class AccountAdapter(DefaultAccountAdapter):
    def is_open_for_signup(self, request):
        """Limit signup based on the number of active trials."""
        return Account.objects.trialing().count() < MAX_TRIALING_USERS
