import datetime

import factory


class EntryFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = "entries.Entry"

    body = factory.Faker("paragraph")
    when = factory.LazyFunction(datetime.date.today)
    user = factory.SubFactory("journal.accounts.tests.factories.UserFactory")


class PromptFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = "entries.Prompt"

    when = factory.LazyFunction(datetime.date.today)
    user = factory.SubFactory("journal.accounts.tests.factories.UserFactory")
    message_id = factory.Sequence(lambda n: f"message_id-{n}")
