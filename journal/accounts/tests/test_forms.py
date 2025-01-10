from journal.accounts.forms import SiginForm
from journal.accounts.models import User


class TestSigninForm:
    def test_create_user(self):
        """A non-existing email creates a new User record."""
        # Check the username uniqueness constraint.
        User.objects.create(email="somethingelse@somewhere.com")
        email = "newuser@somewhere.com"
        data = {"email": email}
        form = SiginForm(data=data)
        is_valid = form.is_valid()

        form.save()
        # Ensure only 1 account is created.
        form.save()

        assert is_valid
        assert User.objects.filter(email=email).count() == 1

    def test_existing_user(self):
        """When a user account exists for an email, use that user."""
        user = User.objects.create(email="test@testing.com")
        data = {"email": user.email}
        form = SiginForm(data=data)
        is_valid = form.is_valid()

        form.save()

        assert is_valid
        assert User.objects.filter(email=user.email).count() == 1

    def test_triggers_signin_link_task(self):
        """The magic link job fires."""

        # FIXME: assert that the outbox has 1 email in it *to* the right email address.

    def test_invalid_email(self):
        """An invalid email is rejected."""
        data = {"email": "not-an-email"}
        form = SiginForm(data=data)

        is_valid = form.is_valid()

        assert not is_valid
        assert "valid email" in form.errors["email"][0]
