import factory
from django.db.models.signals import post_save
from djstripe.models import Event


@factory.django.mute_signals(post_save)
class AccountFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = "accounts.Account"

    user = factory.SubFactory("journal.accounts.tests.factories.UserFactory")


class EventFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = Event

    id = factory.Sequence(lambda n: f"evt_{n}")
    data = factory.LazyFunction(lambda: {})


class UserFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = "accounts.User"

    email = factory.Sequence(lambda n: f"user_{n}@testing.com")
    username = factory.Sequence(lambda n: f"user_{n}")
