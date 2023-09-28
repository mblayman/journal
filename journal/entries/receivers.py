def handle_inbound(sender, event, esp_name, **kwargs):
    message = event.message
    print(
        "Received message from %s (envelope sender %s) with subject '%s'"
        % (message.from_email, message.envelope_sender, message.subject)
    )
