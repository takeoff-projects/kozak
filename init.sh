terraform init && terraform apply -auto-approve

gcloud config set project roi-takeoff-user25
gcloud endpoints services deploy openapi-run.yaml

