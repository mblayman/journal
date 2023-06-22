import datetime

import factory


class EntryFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = "entries.Entry"

    body = factory.Faker("paragraph")
    when = factory.LazyFunction(datetime.date.today)
    user = factory.SubFactory("journal.accounts.tests.factories.UserFactory")
