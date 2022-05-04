
check_install:
	which swagger || GO111MODULE="on" go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_install
	GO111MODULE="on" swagger generate spec -o ./swagger.yaml --scan-models

generate_client:
	cd sdk && swagger generate client -f ../swagger.yaml -A todolist-api

all: serve

.PHONY: deploy

run:
	@ grep -v "#" .env | sed 's/.*/export &/' > /tmp/.env
	@ . /tmp/.env && go run main.go --db_user user --db_pwd password --db_name go-mysql-api --db_port 3306
deploy:
	docker-compose -f ./deploy/local/docker-compose.yml up -d


teardown:
	docker-compose -f ./deploy/local/docker-compose.yml stop

