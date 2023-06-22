from django.contrib import admin

from .models import Entry


@admin.register(Entry)
class EntryAdmin(admin.ModelAdmin):
    list_display = ["id", "user", "when"]
    list_filter = ["when"]
    raw_id_fields = ["user"]
