build-img:
	docker build -t go-server-template .

run-container:
	docker run -d --name go-server-template go-server-template

backup:
	git add .
	git commit -m "backup"
	git push

gen-rest-doc:
	swag init -g internal/delivery/rest/*.go -o ./docs

clear-db:
	if [ -e "./gorm.db" ]; then \
		rm "./gorm.db"; \
		echo "Database file ./gorm.db deleted."; \
	else \
		echo "Database file ./gorm.db does not exist. Nothing to delete."; \
	fi

run:
	go run .

test-run: clear-db run

gen-apidoc:
	swag init -g internal/delivery/*

build-image:
	docker build -t zanzibar-dag .
	
update-dependencies:
	go get -u ./...

backup:
	git add .
	git commit -m "backup"
	git push

gen-grpc:
	protoc --go_out=. --go-grpc_out=. ./domain/delivery/proto/service.proto

grpc-doc:
	docker run --rm \
	-v $(pwd)/grpc-doc:/out \
	-v $(pwd)/internal/delivery/proto:/protos \
	pseudomuto/protoc-gen-doc

rest-doc:
	swag init -g internal/delivery/rest/*