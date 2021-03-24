ip = localhost
mysql = root:asdfg321@tcp(${ip}:3306)/newt

build:
	docker build . -t newt

run:
	docker run -d --name newt_backend -p 3000:3000 -e MYSQL="${mysql}" --entrypoint ./Newt newt

remove-containers:
	@docker container rm newt_backend -f ||true

compose:
	docker-compose up -d

build-and-run: compose remove-containers build run

stop-containers:
	@docker stop newt_backend ||true
	@docker stop --time=5 newt_phpmyadmin_1 ||true
	@docker stop newt_mariadb_1 ||true
	docker-compose down