# email:twilio starter

Use Twilio Email—not SendGrid—at `POST https://comms.twilio.com/v1/Emails`. It is asynchronous and returns an operation ID; use the operation location to track delivery. Official guide: https://www.twilio.com/docs/email/api/overview .

This generator intentionally writes placeholders only. Copy .env.example to .env and provide your own credentials at runtime. Keep provider-specific calls behind the generated repository/provider boundary so replacing this integration does not change routes or domain services.
