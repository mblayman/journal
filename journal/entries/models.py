from __future__ import annotations

import random

from django.db import models

from journal.accounts.models import User


class EntryManager(models.Manager):
    """A manager to provide custom methods for Entry."""

    def get_random_for(self, user: User) -> Entry | None:
        queryset = self.get_queryset().filter(user=user)
        count = queryset.count()
        if count == 0:
            return None

        index = random.choice(range(0, count))  # noqa: S311
        return queryset[index]


class Entry(models.Model):
    """An entry stores the user's writing for the day"""

    class Meta:
        verbose_name_plural = "entries"

    body = models.TextField()
    when = models.DateField(db_index=True)
    user = models.ForeignKey(
        "accounts.User",
        on_delete=models.CASCADE,
        related_name="entries",
    )

    objects = EntryManager()


class PromptManager(models.Manager):
    """A manager to provide custom methods for Prompt."""

    def exists_for(self, user, when):
        """Check if a prompt exists for this user on this date."""
        return self.get_queryset().filter(user=user, when=when).exists()


class Prompt(models.Model):
    """A record to track that a journal prompt was sent to a user"""

    when = models.DateField()
    user = models.ForeignKey(
        "accounts.User",
        on_delete=models.CASCADE,
        related_name="prompts",
    )
    message_id = models.TextField()

    objects = PromptManager()
