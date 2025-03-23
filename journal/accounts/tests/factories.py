import factory
from django.db.models.signals import post_save


@factory.django.mute_signals(post_save)
class AccountFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = "accounts.Account"

    user = factory.SubFactory("journal.accounts.tests.factories.UserFactory")


class UserFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = "accounts.User"

    email = factory.Sequence(lambda n: f"user_{n}@testing.com")
    username = factory.Sequence(lambda n: f"user_{n}")
