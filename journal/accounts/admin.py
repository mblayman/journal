from django.contrib import admin
from django.contrib.auth.admin import UserAdmin
from simple_history.admin import SimpleHistoryAdmin

from .models import Account, User

admin.site.register(User, UserAdmin)


@admin.register(Account)
class AccountAdmin(SimpleHistoryAdmin):
    pass
