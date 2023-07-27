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

        index = random.choice(range(0, count))
        return queryset[index]


class Entry(models.Model):
    """An entry stores the user's writing for the day"""

    class Meta:
        verbose_name_plural = "entries"

    body = models.TextField()
    when = models.DateField()
    user = models.ForeignKey(
        "accounts.User",
        on_delete=models.CASCADE,
        related_name="entries",
    )

    objects = EntryManager()
