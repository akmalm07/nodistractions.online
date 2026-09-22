param([Parameter(Mandatory=$true)][string]$ProjectId, [Parameter(Mandatory=$true)][string]$Bucket, [string]$Location="us-central1")
$ErrorActionPreference = "Stop"
gcloud auth application-default login
gcloud config set project $ProjectId
gcloud storage buckets create "gs://$Bucket" --location=$Location --uniform-bucket-level-access
Write-Host "Bucket provisioned. Configure GCP_PROJECT_ID and GCP_STORAGE_BUCKET in .env."
