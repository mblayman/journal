# from pprint import pprint
from typing import Any

from anymail.signals import AnymailInboundEvent


def handle_inbound(
    sender: Any, event: AnymailInboundEvent, esp_name: str, **kwargs: Any
) -> None:
    message = event.message
    if message is not None:
        # body - come from the message -
        # what field is it stored in? what is the format of that field?
        # message.text OR message.html OR maybe message.stripped_text

        # message.text contains the plain text AND the prompt.
        # We need to trim out the prompt.
        # Idea: trim all lines after journal@email.journeyinbox.com appears.
        # That's imperfect, but a reasonably safe solution.

        # Check event.esp_event
        print("ESP event was")
        # pprint(event.esp_event)
        if event.esp_event is not None:
            for header in event.esp_event.headers:
                print(header)

            print(event.esp_event.body.decode())

        # when - It's Wednesday, Oct. 18, 2023, how are you? (2023-10-18)
        # user -
        # Option 1 - event.from_email - does this match some User in the db?
        # Option 2 - when the prompt email is sent, SendGrid gives a mail ID.
        # print(
        #     "Received message from %s (envelope sender %s) with subject '%s'"
        #     % (message.from_email, message.envelope_sender, message.subject)
        # )
