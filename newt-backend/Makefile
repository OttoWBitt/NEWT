ip = localhost
mysql = root:asdfg321@tcp(${ip}:3306)/newt

build:
	docker build . -t newt

run:
	@mkdir ~/newtFiles ||true
	docker run -d --name newt_backend -v ~/newtFiles:/app/files --network="host" -p 3001:3001 -e MYSQL="${mysql}" --entrypoint ./Newt newt

remove-containers:
	@docker container rm newt_backend -f ||true
	@docker rmi $$(docker images --filter "dangling=true" -q --no-trunc) ||true

compose:
	docker-compose up -d

build-and-run: remove-containers build run

build-and-run-all: compose remove-containers build run

stop-containers:
	@docker stop newt_backend ||true
	@docker stop --time=5 newt_phpmyadmin ||true
	@docker stop newt_mariadb ||true
	docker-compose down
	@docker rmi $$(docker images --filter "dangling=true" -q --no-trunc) ||true

remove-dangling:
	@docker rmi $$(docker images --filter "dangling=true" -q --no-trunc) ||true
