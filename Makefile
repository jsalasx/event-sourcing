build-account-cmd:
	docker build -t kfcregistry.azurecr.io/bank/account-cmd:dev -f account-cmd/Dockerfile .

build-account-query:
	docker build -t kfcregistry.azurecr.io/bank/account-query:dev -f account-query/Dockerfile .

restart-account-cmd:
	kubectl rollout restart deployment account-cmd -n es-bank-account

restart-account-query:
	kubectl rollout restart deployment account-query -n es-bank-account

deploy-account-cmd: build-account-cmd restart-account-cmd
	
deploy-account-query: build-account-query restart-account-query
