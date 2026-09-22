#!/usr/bin/env bash
set -euo pipefail
: "${PROJECT_ID:?}" "${SERVICE:?}" "${REGION:?}"
RUNTIME_SERVICE_ACCOUNT="${RUNTIME_SERVICE_ACCOUNT:-$SERVICE@$PROJECT_ID.iam.gserviceaccount.com}"

gcloud config set project "$PROJECT_ID"
gcloud run deploy "$SERVICE" \
  --source backend \
  --region "$REGION" \
  --service-account "$RUNTIME_SERVICE_ACCOUNT"
