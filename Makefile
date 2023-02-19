buildlambda:
	GOOS=linux go build -ldflags="-s -w" -o ./bin/$(function) functions/$(function)/main.go && zip -jrm ./bin/$(function).zip ./bin/$(function)

previewdeployterraform:
	cd terraform && terraform init && terraform plan

deployterraform:
	cd terraform && terraform apply -auto-approve

previewdestroyterraform:
	cd terraform && terraform init && terraform destroy

destroyterraform:
	cd terraform && terraform destroy -auto-approve