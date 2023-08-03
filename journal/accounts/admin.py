from django.contrib import admin
from django.contrib.auth.admin import UserAdmin as BaseUserAdmin
from simple_history.admin import SimpleHistoryAdmin

from .models import Account, User


@admin.register(User)
class UserAdmin(BaseUserAdmin):
    fieldsets = BaseUserAdmin.fieldsets + (
        (
            "Extra",
            {"fields": ("status", "activate_date", "deactivate_date")},
        ),
    )


@admin.register(Account)
class AccountAdmin(SimpleHistoryAdmin):
    pass
