#!/bin/bash

# Send POST request to webhook with multipart/form-data
curl -X POST http://localhost:8080/webhook \
  -u "testuser:testpass" \
  -H "Content-Type: multipart/form-data; boundary=xYzZY" \
  --data-binary @bin/example_payload.txt \
  -o response.txt

# Check response
if grep -q "ok" response.txt; then
  echo "Webhook responded with 'ok'"
else
  echo "Webhook failed:"
  cat response.txt
  exit 1
fi

# Verify database
echo "Checking database entry for user_id=1, when='2025-03-26':"
sqlite3 ./db.sqlite3 <<EOF
SELECT user_id, "when", body FROM entries_entry WHERE user_id = 1 AND "when" = '2025-03-26';
EOF

# Clean up
rm response.txt
