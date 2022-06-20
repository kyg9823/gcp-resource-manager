# GCP Resource Manager with Fiber

## API List

- Healthcheck: /healthcheck
- GCE Start: /api/v1/gcp/:ProjectId/gce/state?action=start
- GCE Stop: /api/v1/gcp/:ProjectId/gce/state?action=stop
- EC2 Start: /api/v1/aws/ec2/state?action=start
- EC2 Stop: /api/v1/aws/ec2/state?action=stop

## Create GCP IAM

### Role

compute.instances.list
compute.instances.start
compute.instances.stop
compute.zoneOperations.get

```bash
gcloud iam roles create GCP_RESOURCE_MANAGER --project=${PROJECT_ID} \
--title=GCP_RESOURCE_MANAGER \
--description="Role for GCP Resource Manager" \
--permissions=compute.instances.list,compute.instances.start,compute.instances.stop,compute.zoneOperations.get
```

### Service Account

```bash
gcloud iam service-accounts create gcp-resource-manager \
--display-name="GCP Resource Manager" 
gcloud iam service-accounts add-iam-policy-binding gcp-resource-manager@${PROJECT_ID}.iam.gserviceaccount.com \
--role=projects/${PROJECT_ID}/roles/GCP_RESOURCE_MANAGER \
--member=allAuthenticatedUsers
```

## How to deploy on GCP with Cloud Run

```bash
gcloud run deploy gcp-resource-manager \
--image=$ImageName \
--service-account=$ServiceAccount \
--ingress=internal \
--no-allow-unauthenticated \
--region=asia-northeast3 \
```

## How to call on GCP with Cloud Run
