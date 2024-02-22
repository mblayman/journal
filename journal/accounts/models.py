from django.conf import settings
from django.contrib.auth.models import AbstractUser
from django.db import models
from django.db.models.signals import post_save
from django.dispatch import receiver
from django_extensions.db.models import ActivatorModel
from djstripe import webhooks
from hashid_field import HashidAutoField
from simple_history.models import HistoricalRecords


class AccountManager(models.Manager):
    def active(self):
        """Get all the active accounts."""
        qs = self.get_queryset()
        return qs.filter(status__in=self.model.ACTIVE_STATUSES)


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

    objects = AccountManager()
    history = HistoricalRecords()


class User(AbstractUser, ActivatorModel):
    pass


@receiver(post_save, sender=User)
def create_account(sender, instance, created, **kwargs):
    """A new user gets an associated account."""
    if created:
        Account.objects.create(user=instance)


@webhooks.handler("checkout.session.completed")
def handle_checkout_session_completed(event, **kwargs):
    """Transition the account to an active state.

    This event occurs after a user provides their checkout payment information.
    """
    event_data = event.data["object"]
    # The payments gateway sets the user ID in the client reference ID field.
    user_id = int(event_data["client_reference_id"])
    Account.objects.filter(user_id=user_id).update(status=Account.Status.ACTIVE)
