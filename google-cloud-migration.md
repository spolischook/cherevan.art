# Hugo Site Migration Plan: Netlify to Firebase Hosting

This document outlines the steps to migrate the Hugo static website from Netlify to Firebase Hosting, leveraging Google's infrastructure for cost-effective and high-performance hosting with easy integration into the Google ecosystem.

## Phase 1: Firebase Setup & Initial Manual Deployment

This phase focuses on setting up Firebase Hosting for your project and performing an initial manual deployment to ensure everything is working correctly.

### 1. Firebase Project Setup

*   **Navigate to the Firebase Console:**
    *   Go to the [Firebase Console](https://console.firebase.google.com/).
*   **Add or Select a Project:**
    *   You can **Add a new Firebase project** or **select your existing Google Cloud Project** (`cherevan-art-workspace`) to add Firebase services to it. Adding Firebase to an existing GCP project is often preferred for better integration.
        *   If adding to an existing GCP project: Click "Add project", then select your `cherevan-art-workspace` project from the dropdown or search for it.
        *   If creating a new Firebase project: Click "Add project", give it a name (e.g., `cherevan-art-website`), and follow the prompts. You can choose whether or not to enable Google Analytics for this project (recommended).
*   **Enable Billing (if not already on the GCP project):**
    *   Firebase Hosting's free tier (Spark plan) is very generous. However, if you anticipate exceeding free limits or want to use paid Firebase features later, ensure your underlying Google Cloud Project (`cherevan-art-workspace`) has billing enabled. For just hosting a static site with your described traffic, the Spark plan should suffice.

### 2. Install Firebase CLI

*   The Firebase Command Line Interface (CLI) is used to manage and deploy your Firebase projects.
*   **Installation (if not already installed):**
    *   Open your terminal or command prompt.
    *   Install the Firebase CLI via npm (Node Package Manager, which comes with Node.js). If you don't have Node.js, install it first from [nodejs.org](https://nodejs.org/).
      ```bash
      npm install -g firebase-tools
      ```
*   **Verify Installation:**
    ```bash
    firebase --version
    ```

### 3. Login to Firebase

*   **Authenticate the Firebase CLI with your Google account:**
    ```bash
    firebase login
    ```
    *   This will open a browser window for you to log in and authorize the CLI.
    *   If you have multiple Google accounts, ensure you log in with the one associated with your Firebase/GCP project.

### 4. Initialize Firebase Hosting in Your Hugo Project

*   **Navigate to your Hugo project root directory in the terminal:**
    ```bash
    cd /path/to/your/cherevan.art
    ```
*   **Run the Firebase initialization command:**
    ```bash
    firebase init hosting
    ```
*   **Follow the prompts:**
    *   **"Are you ready to proceed?"**: `Yes`
    *   **"Which Firebase project do you want to associate as default for this directory?"**: Use the arrow keys to select your Firebase project (`cherevan-art-workspace` or the new one you created) and press Enter.
    *   **"What do you want to use as your public directory?"**: Enter `public` (this is Hugo's default output directory). Press Enter.
    *   **"Configure as a single-page app (rewrite all urls to /index.html)?"**: `No` (Hugo generates individual HTML files, so this is typically not needed). Press Enter.
    *   **"Set up automatic builds and deploys with GitHub?"**: `No` (for now, we will set this up manually in Phase 2. If you choose Yes, it will try to create a GitHub Action for you, but we want more control initially).
    *   **"File public/index.html already exists. Overwrite?"**: If it asks this because you've already run `hugo`, choose `No`.
*   This process will create two files in your Hugo project root:
    *   `.firebaserc`: Configures which Firebase project this local directory is associated with.
    *   `firebase.json`: Configures your Firebase Hosting settings (like the public directory).

### 5. Configure Custom Domain & SSL

*   **Navigate to Firebase Hosting in the Console:**
    *   Go to the [Firebase Console](https://console.firebase.google.com/), select your project.
    *   In the left-hand menu, go to "Build" > "Hosting".
*   **Add Custom Domain:**
    *   Click the "Add custom domain" button.
    *   Enter your domain (e.g., `cherevan.art`).
    *   You can also add `www.cherevan.art` if you want both to point to your site.
*   **Verify Domain Ownership:**
    *   Firebase will provide you with TXT records (or other methods) to add to your DNS provider (DigitalOcean in your case) to verify ownership. Follow the instructions.
    *   This verification can take some time for DNS to propagate.
*   **SSL Certificates:**
    *   Once your domain is verified and DNS is pointing to Firebase, Firebase will automatically provision and renew SSL certificates for your custom domain(s) **free of charge**.
    *   The status will change from "Needs setup" or "Verifying" to "Connected".

### 6. Manual Deployment & Test

*   **Build your Hugo site:**
    *   In your Hugo project root directory, run:
      ```bash
      hugo
      ```
      This generates your static site into the `public/` directory.
*   **Deploy to Firebase Hosting:**
    *   In your Hugo project root directory, run:
      ```bash
      firebase deploy --only hosting
      ```
      *   The `--only hosting` flag ensures only hosting content is deployed (useful if you use other Firebase services like Functions).
*   **Verify:**
    *   After deployment, the CLI will give you a Firebase-hosted URL (e.g., `your-project-id.web.app`) and your custom domain URL(s) once they are connected.
    *   Access your site via your custom domain (e.g., `https://cherevan.art`) to verify.
    *   Test navigation, HTTPS, and check for any broken links or missing assets.

## Phase 2: Adapt GitHub Actions Workflows for Firebase Deployment (using Workload Identity Federation)

This phase details modifying your GitHub Actions workflow (`deploy.yml`) to deploy to Firebase Hosting using Google Cloud Workload Identity Federation. This method is more secure than using long-lived tokens like `FIREBASE_TOKEN` as it uses short-lived credentials.

### 1. Google Cloud Setup for Workload Identity Federation

**A. Environment Variables (ensure these are set in your local terminal for the setup commands):**

```bash
export PROJECT_ID="cherevan-art-workspace"
export SERVICE_ACCOUNT_EMAIL="github-actions-firebase-deploy@${PROJECT_ID}.iam.gserviceaccount.com"
export GITHUB_REPO="spolischook/cherevan.art" # Your GitHub org/username and repo name

# These will be derived/set by the commands below
export POOL_ID="github-actions-pool"
export PROVIDER_ID="github-provider"
```

**B. Create Workload Identity Pool:**

If it doesn't exist already:
```bash
gcloud iam workload-identity-pools create "${POOL_ID}" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --display-name="GitHub Actions Pool"
```

**C. Get Full Workload Identity Pool ID:**

```bash
export WORKLOAD_IDENTITY_POOL_FULL_ID=$(gcloud iam workload-identity-pools describe "${POOL_ID}" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --format="value(name)")
```

**D. Create Workload Identity Provider in the Pool:**

If it doesn't exist already:
```bash
gcloud iam workload-identity-pools providers create-oidc "${PROVIDER_ID}" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --workload-identity-pool="${POOL_ID}" \
  --display-name="GitHub Provider" \
  --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository" \
  --issuer-uri="https://token.actions.githubusercontent.com"
```

**E. Create or Verify Service Account:**

List existing service accounts to check:
```bash
gcloud iam service-accounts list --project="${PROJECT_ID}"
```
If your intended service account (`github-actions-firebase-deploy@cherevan-art-workspace.iam.gserviceaccount.com`) doesn't exist, or if you used a slightly different name (e.g., `deployer` vs `deploy`), create it (adjust display name if needed):
```bash
gcloud iam service-accounts create "github-actions-firebase-deploy" \
  --project="${PROJECT_ID}" \
  --display-name="GitHub Actions Firebase Deployer"
```
*Ensure your `SERVICE_ACCOUNT_EMAIL` variable matches the actual email of the created/existing service account.* 

**F. Grant Service Account the Firebase Hosting Admin Role:**

```bash
gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
  --member="serviceAccount:${SERVICE_ACCOUNT_EMAIL}" \
  --role="roles/firebasehosting.admin"
```

**G. Allow GitHub Actions to Impersonate the Service Account:**

This binds the `roles/iam.workloadIdentityUser` to your service account, allowing principals matching the specified GitHub repository to impersonate it.
```bash
gcloud iam service-accounts add-iam-policy-binding "${SERVICE_ACCOUNT_EMAIL}" \
  --project="${PROJECT_ID}" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/${WORKLOAD_IDENTITY_POOL_FULL_ID}/attribute.repository/${GITHUB_REPO}"
```

### 2. Update GitHub Actions Workflow (`.github/workflows/deploy.yml`)

The following workflow configuration uses the Workload Identity Federation setup:

```yaml
name: Deploy to Firebase Hosting

on:
  push:
    branches:
      - main

jobs:
  deploy:
    runs-on: ubuntu-latest
    permissions:
      contents: 'read'
      id-token: 'write' # Required for OIDC token generation

    steps:
    - name: Checkout repository
      uses: actions/checkout@v3
      with:
        submodules: 'recursive'
        fetch-depth: 1

    - name: Set up Hugo
      uses: peaceiris/actions-hugo@v2
      with:
        hugo-version: 'latest' # or your specific version
        # extended: true # Uncomment if you use SCSS/SASS

    - name: Setup Node.js
      uses: actions/setup-node@v3
      with:
        node-version: 18 # Or your preferred Node.js version
        cache: 'npm'

    # If you have npm dependencies for your Hugo build (e.g., PostCSS)
    # - name: Install npm dependencies
    #   run: npm ci

    - name: Build Hugo site
      run: hugo --minify --baseURL="https://www.cherevan.art"

    - name: Authenticate to Google Cloud
      id: auth
      uses: 'google-github-actions/auth@v2'
      with:
        workload_identity_provider: 'projects/769799666531/locations/global/workloadIdentityPools/github-actions-pool/providers/github-provider' # Full provider ID
        service_account: 'github-actions-firebase-deploy@cherevan-art-workspace.iam.gserviceaccount.com' # Your service account email

    - name: Install Firebase CLI
      run: npm install -g firebase-tools

    - name: Deploy to Firebase Hosting
      run: firebase deploy --only hosting --project cherevan-art-workspace --non-interactive --debug
```

### 3. GitHub Secrets Cleanup

*   You can now **delete** the `FIREBASE_TOKEN` secret from your GitHub repository settings (Settings > Secrets and variables > Actions), as it's no longer used by this workflow.

This completes the migration to Firebase Hosting using a secure Workload Identity Federation setup for your GitHub Actions CI/CD pipeline.

## Phase 3: Post-Deployment Verification & Cleanup

*(Original content for Phase 3 can remain or be reviewed)*