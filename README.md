# GCP Resource Manager with Fiber

## API List

- Healthcheck: /healthcheck
- GCE Start: /api/v1/gcp/:ProjectId/gce/state?action=start
- GCE Stop: /api/v1/gcp/:ProjectId/gce/state?action=stop

## Create GCP IAM

### Role

compute.instances.list
compute.instances.start
compute.instances.stop
compute.zoneOperations.get

## How to deploy on GCP with Cloud Run

```bash
gcloud run deploy gcp-resource-manager \
--image=$ImageName \
--service-account=$ServiceAccount \
--ingress=internal \
--region=asia-northeast3
```
