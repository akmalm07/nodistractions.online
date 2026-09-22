#!/usr/bin/env bash
set -euo pipefail
PROJECT_ID="${1:?usage: provision-gcs.sh PROJECT_ID BUCKET [location]}"
BUCKET="${2:?usage: provision-gcs.sh PROJECT_ID BUCKET [location]}"
LOCATION="${3:-us-central1}"
gcloud auth application-default login
gcloud config set project "$PROJECT_ID"
gcloud storage buckets create "gs://$BUCKET" --location="$LOCATION" --uniform-bucket-level-access
printf "Bucket provisioned. Configure GCP_PROJECT_ID and GCP_STORAGE_BUCKET in .env.\n"
