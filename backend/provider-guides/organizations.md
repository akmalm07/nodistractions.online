# organizations starter

Starter domain boundaries for organizations, invitations, and role checks belong in the service layer. Add selected persistence operations behind repository interfaces.

This generator intentionally writes placeholders only. Copy .env.example to .env and provide your own credentials at runtime. Keep provider-specific calls behind the generated repository/provider boundary so replacing this integration does not change routes or domain services.
