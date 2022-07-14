postgres :
	docker run --name rsm -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres

startdb:
	docker exec -it rsm psql -U postgres

createdb:
	docker exec -it rsm createdb rsmstore -U postgres

dropdb:
	docker exec -it rsm dropdb rsmstore -U postgres

test:
	go test -v -cover ./...



.PHONY: postgres startdb