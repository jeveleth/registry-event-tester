down:
	docker-compose down

up:
	docker-compose up --build -d 

logs-all:
	docker-compose logs -f

logs-reg:
	docker-compose logs -f registry

logs-app:
	docker-compose logs -f app

image-push:
	docker pull busybox
	docker tag busybox:latest localhost:5000/myfirstimage:latest
	docker push localhost:5000/myfirstimage:latest

start:
	-${MAKE} down
	-${MAKE} up
