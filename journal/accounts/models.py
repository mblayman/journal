from __future__ import annotations

from django.conf import settings
from django.contrib.auth.models import AbstractUser
from django.db import models
from django.db.models.signals import post_save
from django.dispatch import receiver
from django_extensions.db.models import ActivatorModel
from hashid_field import HashidAutoField
from simple_history.models import HistoricalRecords


class AccountManager(models.Manager):
    def active(self):
        """Get all the active accounts."""
        qs = self.get_queryset()
        return qs.filter(status__in=self.model.ACTIVE_STATUSES)

    def promptable(self):
        """Get all the accounts that can receive prompts."""
        active = self.active()
        return active.filter(verified=True)


class Account(models.Model):
    """Account holds the user's state"""

    class Status(models.IntegerChoices):
        TRIALING = 1
        ACTIVE = 2
        EXEMPT = 3
        CANCELED = 4
        TRIAL_EXPIRED = 5

    ACTIVE_STATUSES = (Status.TRIALING, Status.ACTIVE, Status.EXEMPT)

    id = HashidAutoField(primary_key=True, salt=f"account{settings.HASHID_FIELD_SALT}")
    user = models.OneToOneField(
        "accounts.User",
        on_delete=models.CASCADE,
    )
    status = models.IntegerField(
        choices=Status.choices,
        default=Status.TRIALING,
        db_index=True,
    )
    verified = models.BooleanField(
        help_text="An account is verified if they log in through email at least once",
        default=False,
        db_index=True,
    )

    objects = AccountManager()
    history = HistoricalRecords()


class User(AbstractUser, ActivatorModel):
    pass


@receiver(post_save, sender=User)
def create_account(sender, instance, created, raw, **kwargs):
    """A new user gets an associated account."""
    if created and not raw:
        Account.objects.create(user=instance)
