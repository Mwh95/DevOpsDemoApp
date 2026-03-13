# DRAFT

# Configure GKE Ingress via Google Cloud Console

This runbook describes how to create and configure the DemoApp Ingress using the Google Cloud Console (UI).

## Prerequisites

- A GKE cluster with Keycloak deployed (e.g. after running `./scripts/deploy-gcp.sh`).
- `kubectl` configured to use the cluster, or use the Cloud Console **Kubernetes Engine > Workloads** and **Services & Ingress**.

## Option 1: Create Ingress from YAML in Console

1. Open [Google Cloud Console](https://console.cloud.google.com/) and select your project.
2. Go to **Kubernetes Engine > Services & Ingress**.
3. Click **CREATE INGRESS** (or **+ CREATE** and choose **Ingress**).
4. Choose **Advanced** (or equivalent) to paste YAML.
5. Paste the contents of [../k8s/gcp/ingress.yaml](../k8s/gcp/ingress.yaml) from this repository.
6. Optionally add annotations:
   - **Static IP**: add `kubernetes.io/ingress.global-static-ip-name: "YOUR_STATIC_IP_NAME"` (create the address under **VPC network > IP addresses** first).
   - **Managed certificate**: add `networking.gke.io/managed-certificates: "YOUR_MANAGED_CERT_NAME"` (create the ManagedCertificate resource first).
7. Click **Create**. The Ingress will get an external IP after a few minutes.

## Option 2: Create Ingress from Cloud Shell

1. In Cloud Console, open **Cloud Shell** (icon at top right).
2. Clone or upload this repo, then run:
   ```bash
   kubectl apply -f k8s/gcp/ingress.yaml
   ```
3. Check the Ingress: `kubectl get ingress demoapp-ingress`.

## Optional: Reserve a static IP

1. **VPC network > IP addresses**.
2. **Reserve external static address**.
3. Name (e.g. `demoapp-ingress-ip`), type **Global** (for HTTP(S) load balancing).
4. Add the annotation to the Ingress as in Option 1.

## Optional: HTTPS with Managed Certificate

1. Create a `ManagedCertificate` resource (e.g. in the same namespace as the Ingress) pointing to your domain.
2. Add the annotation `networking.gke.io/managed-certificates: "your-cert-name"` to the Ingress.
3. Ensure your Ingress has a rule with the host set to your domain so the certificate is used.

See [GKE Ingress configuration](https://cloud.google.com/kubernetes-engine/docs/how-to/ingress-configuration) for details.

## Access

- After the Ingress has an **ADDRESS**: open `http://<ADDRESS>/auth` for Keycloak.
- For a custom domain, point the DNS A record to the Ingress ADDRESS and (if using HTTPS) configure the managed certificate as above.
