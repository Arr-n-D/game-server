.PHONY: init_dev init_prod plan_dev plan_prod apply_dev apply_prod apply_auto_approve_dev apply_auto_approve_prod destroy_dev destroy_prod

INIT_DEV := -backend-config="key=heavenly-dragons-game-server-dev"
# INIT_PROD := -backend-config="bucket=progitek-terraform-prod-kjs4jer0" -backend-config="key=dentitek-backend-prod"
PARAMS_DEV := -var-file=./vars/development.tfvars
# PARAMS_PROD := -var-file=./vars/production.tfvars

init_dev:
	@terraform init $(INIT_DEV) -reconfigure

init_prod:
	@terraform init $(INIT_PROD) -reconfigure

plan_dev:
	@terraform plan $(PARAMS_DEV)

plan_prod:
	@terraform plan $(PARAMS_PROD)

apply_auto_approve_dev:
	@terraform apply -auto-approve $(PARAMS_DEV)

apply_auto_approve_prod:
	@terraform apply -auto-approve $(PARAMS_PROD)

apply_dev:
	@terraform apply $(PARAMS_DEV)

apply_prod:
	@terraform apply $(PARAMS_PROD)

destroy_dev:
	@terraform destroy $(PARAMS_DEV)

destroy_prod:
	@terraform destroy $(PARAMS_PROD)
