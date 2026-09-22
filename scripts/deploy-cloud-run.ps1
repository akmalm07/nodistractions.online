param(
  [Parameter(Mandatory = $true)][string]$ProjectId,
  [Parameter(Mandatory = $true)][string]$Service,
  [Parameter(Mandatory = $true)][string]$Region,
  [string]$RuntimeServiceAccount = "$Service@$ProjectId.iam.gserviceaccount.com"
)

$ErrorActionPreference = "Stop"

gcloud config set project $ProjectId
gcloud run deploy $Service `
  --source backend `
  --region $Region `
  --service-account $RuntimeServiceAccount
