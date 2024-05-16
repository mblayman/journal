import factory
from djstripe.models import Customer, Price, Product


class CustomerFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = Customer

    id = factory.Sequence(lambda n: f"cus_{n}")


class ProductFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = Product

    id = factory.Sequence(lambda n: f"product_{n}")


class PriceFactory(factory.django.DjangoModelFactory):
    class Meta:
        model = Price

    id = factory.Sequence(lambda n: f"price_{n}")
    active = True
    livemode = False
    product = factory.SubFactory(ProductFactory)
