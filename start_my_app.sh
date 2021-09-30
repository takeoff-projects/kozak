terraform init 
terraform plan -out plan.tfplan
terraform apply plan.tfplan -auto-approve

gcloud endpoints services deploy openapi-run.yaml

