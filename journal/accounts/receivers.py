def verify_account(sender, request, user, **kwargs):
    """Mark an account as verified."""
    user.account.verified = True
    user.account.save()
