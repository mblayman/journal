import factory


class UserFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = "accounts.User"
