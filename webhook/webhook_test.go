package webhook

import (
	"bytes"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mblayman/journal/model"
)

const rawPayload = `--xYzZY
Content-Disposition: form-data; name="to"

JourneyInbox Journal <journal.abcdef1@email.journeyinbox.com>
--xYzZY
Content-Disposition: form-data; name="email"

Received: from mail-pl1-f177.google.com (mxd [209.85.214.177]) by mx.sendgrid.net with ESMTP id x7jN7BtxS8SV-4NoDN_S7g for <journal.abcdef1@email.journeyinbox.com>; Sat, 29 Mar 2025 14:17:18.777 +0000 (UTC)
Received: by mail-pl1-f177.google.com with SMTP id d9443c01a7336-225df540edcso80958095ad.0
        for <journal.abcdef1@email.journeyinbox.com>; Sat, 29 Mar 2025 07:17:18 -0700 (PDT)
DKIM-Signature: v=1; a=rsa-sha256; c=relaxed/relaxed;
        d=gmail.com; s=20230601; t=1743257838; x=1743862638; darn=email.journeyinbox.com;
        h=to:subject:message-id:date:from:in-reply-to:references:mime-version
         :from:to:cc:subject:date:message-id:reply-to;
        bh=gn+oz7lem9ZEm8yueTqOpZsnF/+F0G/wCgZUE+JCpdY=;
        b=E9kD+dWM4kHwQ63e9FJ5//5izAkvNy+3xV+K/mLFWGWjTbod9WIGsIc4uQK0Xjv4fe
         ww1noabcdefKw7cs/k8cILSRgP39A/aVi9lsnG6Eku0wcWiBk2B0L4AeH7Q1yJG91Q1f
         /+aRvluxoDZcQD07zF6g109lB2kVJL1S2vLnrGammrLiNF56csQdzYpbQ54amDVzdRgh
         O/9PHx7jOnNQWabcdefsUzqBHrQLsq9Moep9Wa4hBmd7Um+09bra+zspGro/iYz2B1q
         w4tabcdefgjgYC2tS5ZsbEyfImCBiwuSlDliwHyqzOB7a6babcdefG7Vjom6zhJ3zE
         +RMA==
From: Matthew Layman <test@somewhere.com>
Date: Sat, 29 Mar 2025 10:17:07 -0400
Message-ID: <abcdefrm44jCDsrnabcdefxdZ+jmdATabcdefgn68Xo7RbiPQ@mail.gmail.com>
Subject: Re: It's Wednesday, Mar. 26, 2025. How are you?
To: JourneyInbox Journal <journal.abcdef1@email.journeyinbox.com>
Content-Type: multipart/alternative; boundary="00000000000028c08d06317bd809"

--00000000000028c08d06317bd809
Content-Type: text/plain; charset="UTF-8"

I got up this morning at 8:30 and brushed my teeth, then left to go to Cafe
Ibiza to meet with Jared. Lorem ipsum dolor sit amet, consectetur adipiscing
elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.

Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in
voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint
occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit
anim id est laborum.

On Wed, Mar 26, 2025, 9:00 AM JourneyInbox Journal <
journal.abcdef1@email.journeyinbox.com> wrote:

> Reply to this prompt to update your journal.
>
> On your journey on Saturday, Apr. 14, 2018 (6 years, 11 months ago), you
> wrote:
>
> The weather was nice so I worked outside like crazy today. I got weed cloth
> down for the new garden, chopped down the dogwood bush that never thrived,
> and got the lawn mower prepared for the season.
>
> Mark helped me a lot today. I showed him how to use a hammer while we put
> twine on the strawberry box. He said that helping me was his favorite part
> of the day. I think it was my favorite part too.
>

--00000000000028c08d06317bd809
Content-Type: text/html; charset="UTF-8"
Content-Transfer-Encoding: quoted-printable

<div dir=3D"ltr"><div dir=3D"auto">I got up this morning at 8:30 and brushe=
d my teeth, then left to go to Cafe Ibiza to meet with Jared. [... HTML content ...]
--00000000000028c08d06317bd809--
--xYzZY--
`

func TestWebhookHandler(t *testing.T) {
	username := "testuser"
	password := "testpass"

	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", log.LstdFlags)

	processor := func(emailContent model.EmailContent) {
		expectedTo := "JourneyInbox Journal <journal.abcdef1@email.journeyinbox.com>"
		if emailContent.To != expectedTo {
			t.Errorf("Expected output to contain %q, got %q", expectedTo, emailContent.To)
		}
		expectedSubject := "Re: It's Wednesday, Mar. 26, 2025. How are you?"
		if emailContent.Subject != expectedSubject {
			t.Errorf("Expected output to contain %q, got %q", expectedSubject, emailContent.Subject)
		}
		expectedTextSnippet := "I got up this morning at 8:30 and brushed my teeth"
		if !strings.Contains(emailContent.Text, expectedTextSnippet) {
			t.Errorf("Expected output to contain %q, got %q", expectedTextSnippet, emailContent.Text)
		}

	}

	handler := WebhookHandler(username, password, processor, logger)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	err := writer.SetBoundary("xYzZY")
	if err != nil {
		t.Fatalf("Failed to set boundary: %v", err)
	}

	emailPart := strings.SplitN(rawPayload, `Content-Disposition: form-data; name="email"`, 2)[1]
	emailPart = strings.SplitN(emailPart, "--xYzZY--", 2)[0]
	err = writer.WriteField("email", strings.TrimSpace(emailPart))
	if err != nil {
		t.Fatalf("Failed to write email field: %v", err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatalf("Failed to close multipart writer: %v", err)
	}

	t.Logf("Constructed request body:\n%s", body.String())

	req := httptest.NewRequest(http.MethodPost, "/webhook", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.SetBasicAuth(username, password)

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", recorder.Code)
	}
	if recorder.Body.String() != "ok" {
		t.Errorf("Expected body 'ok', got %q", recorder.Body.String())
	}
}

func TestWebhookHandlerUnauthorized(t *testing.T) {
	username := "testuser"
	password := "testpass"
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", log.LstdFlags)
	processor := func(model.EmailContent) {}
	handler := WebhookHandler(username, password, processor, logger)

	// Create a request without auth
	req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 Unauthorized, got %d", recorder.Code)
	}
	if !strings.Contains(recorder.Header().Get("WWW-Authenticate"), `Basic realm="Restricted"`) {
		t.Errorf("Expected WWW-Authenticate header, got %v", recorder.Header().Get("WWW-Authenticate"))
	}
}

func TestWebhookHandlerMethodNotAllowed(t *testing.T) {
	username := "testuser"
	password := "testpass"
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", log.LstdFlags)
	processor := func(model.EmailContent) {}
	handler := WebhookHandler(username, password, processor, logger)

	req := httptest.NewRequest(http.MethodGet, "/webhook", nil)
	req.SetBasicAuth(username, password)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405 Method Not Allowed, got %d", recorder.Code)
	}
}
