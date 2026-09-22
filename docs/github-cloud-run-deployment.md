# Automatic Cloud Run deployments from GitHub

The workflow at `.github/workflows/ci.yml` tests every pull request and push. A push to `main` deploys the backend only after the Go and frontend jobs succeed. The frontend is built and tested, but it needs a separate static host or Cloud Run service if you want it deployed too.

The workflow uses GitHub OIDC and Google Workload Identity Federation (WIF), not a downloaded service-account key. Restrict the federation provider to the exact GitHub repository. The workflow itself only deploys on pushes to `main`.

## 1. Create the Google Cloud resources

Set these values in a local shell before running the commands. `GITHUB_REPOSITORY` must be the exact `owner/repository` name from GitHub.

```bash
export PROJECT_ID="your-gcp-project-id"
export REGION="us-central1"
export REPOSITORY="containers"
export CLOUD_RUN_SERVICE="nodistractions-api"
export RUNTIME_SERVICE_ACCOUNT_NAME="nodistractions-runtime"
export DEPLOYER_SERVICE_ACCOUNT_NAME="github-deployer"
export GITHUB_REPOSITORY="owner/repository"

gcloud config set project "$PROJECT_ID"
gcloud config set api_endpoint_overrides/secretmanager "https://secretmanager.$REGION.rep.googleapis.com/"
gcloud services enable run.googleapis.com artifactregistry.googleapis.com \
  iamcredentials.googleapis.com secretmanager.googleapis.com

gcloud artifacts repositories create "$REPOSITORY" \
  --repository-format=docker --location="$REGION"
gcloud iam service-accounts create "$RUNTIME_SERVICE_ACCOUNT_NAME"
gcloud iam service-accounts create "$DEPLOYER_SERVICE_ACCOUNT_NAME"
```

If the Artifact Registry repository or service accounts already exist, skip their create commands.

## 2. Store the dotenv configuration secret

Create the regional `nodistractions-secrets` Secret Manager secret from the `.env` file. Its payload is parsed as dotenv and must include `DATABASE_URL`, `DATABASE_URL_UNPOOLED`, `SESSION_SECRET`, `JWT_ISSUER`, `JWT_AUDIENCE`, and `FRONTEND_URL`. Never put those values in GitHub repository variables or source files.

```bash
gcloud secrets create nodistractions-secrets \
  --project="$PROJECT_ID" \
  --location="$REGION" \
  --data-file=.env
```

Grant both the running API and the GitHub deployer access to that secret. The deployer uses it only to run the schema migration; the runtime reads it on startup with Application Default Credentials.

```bash
export RUNTIME_SERVICE_ACCOUNT="$RUNTIME_SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com"
export DEPLOYER_SERVICE_ACCOUNT="$DEPLOYER_SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com"

for SERVICE_ACCOUNT in "$RUNTIME_SERVICE_ACCOUNT" "$DEPLOYER_SERVICE_ACCOUNT"; do
  gcloud secrets add-iam-policy-binding nodistractions-secrets \
    --location="$REGION" \
    --member="serviceAccount:$SERVICE_ACCOUNT" \
    --role="roles/secretmanager.secretAccessor"
done
```

## 3. Give the deployer only deployment permissions

```bash
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:$DEPLOYER_SERVICE_ACCOUNT" \
  --role="roles/run.admin"
gcloud artifacts repositories add-iam-policy-binding "$REPOSITORY" \
  --location="$REGION" \
  --member="serviceAccount:$DEPLOYER_SERVICE_ACCOUNT" \
  --role="roles/artifactregistry.writer"
gcloud iam service-accounts add-iam-policy-binding "$RUNTIME_SERVICE_ACCOUNT" \
  --member="serviceAccount:$DEPLOYER_SERVICE_ACCOUNT" \
  --role="roles/iam.serviceAccountUser"
```

## 4. Configure GitHub OIDC federation

```bash
gcloud iam workload-identity-pools create github \
  --location=global --display-name="GitHub Actions"

export WIF_POOL="$(gcloud iam workload-identity-pools describe github \
  --location=global --format='value(name)')"

gcloud iam workload-identity-pools providers create-oidc github \
  --location=global \
  --workload-identity-pool=github \
  --issuer-uri="https://token.actions.githubusercontent.com" \
  --attribute-mapping="google.subject=assertion.sub,attribute.repository=assertion.repository" \
  --attribute-condition="assertion.repository == '$GITHUB_REPOSITORY'"

export WIF_PROVIDER="$(gcloud iam workload-identity-pools providers describe github \
  --location=global --workload-identity-pool=github --format='value(name)')"

gcloud iam service-accounts add-iam-policy-binding "$DEPLOYER_SERVICE_ACCOUNT" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/$WIF_POOL/attribute.repository/$GITHUB_REPOSITORY"
```

## 5. Push the workflow

This configured project already includes its Google Cloud resource identifiers in `.github/workflows/ci.yml`, so no GitHub Actions variables are required. If you fork the project or change any Cloud resource, update those non-secret identifiers in that workflow. Do not store the dotenv payload or database URLs in GitHub.

Every successful push to `main` reads the dotenv configuration from Secret Manager, runs the database migration, pushes an immutable image tagged with the commit SHA, and updates the Cloud Run backend. Configure Cloud Run ingress and unauthenticated access separately; the workflow intentionally does not change that security policy.
