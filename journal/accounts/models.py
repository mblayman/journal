from django.contrib.auth.models import AbstractUser
from django.db import models
from django.db.models.signals import post_save
from django.dispatch import receiver


class Account(models.Model):
    """Account holds the user's state"""

    class Status(models.IntegerChoices):
        TRIALING = 1
        ACTIVE = 2
        EXEMPT = 3
        CANCELED = 4
        TRIAL_EXPIRED = 5

    user = models.OneToOneField(
        "accounts.User",
        on_delete=models.CASCADE,
    )
    status = models.IntegerField(
        choices=Status.choices,
        default=Status.TRIALING,
        db_index=True,
    )


class User(AbstractUser):
    pass


@receiver(post_save, sender=User)
def create_account(sender, instance, created, **kwargs):
    """A new user gets an associated account."""
    if created:
        Account.objects.create(user=instance)
