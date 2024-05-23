from django.contrib import admin

from .models import Entry, Prompt


@admin.register(Entry)
class EntryAdmin(admin.ModelAdmin):
    list_display = ["id", "user", "when"]
    list_filter = ["when"]
    raw_id_fields = ["user"]


@admin.register(Prompt)
class PromptAdmin(admin.ModelAdmin):
    list_display = ["id", "user", "when"]
    list_filter = ["when"]
    raw_id_fields = ["user"]
