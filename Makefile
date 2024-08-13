# Test Database
.PHONY: stoptestdb_container starttestdb_container rmtestdb_container maketestdb_container createtestdb droptestdb exectestdb
stoptestdb_container:
	docker stop lifestyledb_test

starttestdb_container:
	docker start lifestyledb_test

rmtestdb_container:
	docker rm lifestyledb_test

maketestdb_container:
	docker run --name lifestyledb_test -e POSTGRES_PASSWORD=testpassword -e POSTGRES_USER=root -p 5432:5432 -d postgres:16.1 

createtestdb:
	docker exec -it lifestyledb_test psql -U root -c "create database lfdb;"

droptestdb:
	docker exec -it lifestyledb_test psql -u root -p -e "drop database lfdb;"

exectestdb:
	docker exec -it lifestyledb_test psql -U root

# Dev Database
.PHONY: stopdevdb_container startdevdb_container rmdevdb_container makedevdb_container createdevdb dropdevdb execdevdb
stopdevdb_container:
	docker stop lifestyledb

startdevdb_container:
	docker start lifestyledb

rmdevdb_container:
	docker rm lifestyledb

makedevdb_container:
	docker run --name lifestyledb -e POSTGRES_PASSWORD=Goldendragon1375 -e POSTGRES_USER=root -p 5432:5432 -d postgres:16.1 

createdevdb:
	docker exec -it lifestyledb psql -U root -c "create database lfdb;"

dropdevdb:
	docker exec -it lifestyledb psql -u root -p -e "drop database lfdb;"

execdevdb:
	docker exec -it lifestyledb psql -U root


# Fly.io database
.PHONY: connect_to_fitfin_postgres
connect_to_fitfin_postgres:
	fly postgres connect -a saeidlifestyledb -u lifestylebackend -p kvqUQ2B2Oy9TtSv 


# sqlc migrations
.PHONY: newmigrations psqlmigrateup psqlmigratedown
newmigrations:
	migrate create -ext sql -dir db/migration -seq lifeStyleMigrations

psqlmigrateup:
	migrate -path db/migration -database "postgres://root:Goldendragon1375@0.0.0.0:5432/lfdb?sslmode=disable" -verbose up 

psqlmigratedown:
# if you want to drop the last one, do "1" after down (it's the other way around)
	migrate -path db/migration -database "postgres://root:Goldendragon1375@0.0.0.0:5432/lfdb?sslmode=disable" -verbose down
# ! If in migration you saw the dirty, it will tell you the dirty version. then force it to the PREVIOUS version
# Example:
# panic: Dirty database version 13. Fix and force version
# migrate -path db/migration -database "postgres://root:Goldendragon1375@0.0.0.0:5432/lfdb?sslmode=disable" force 12

# Compilation
fitfin:
	go build -o fitfin main.go 

.PHONY: serve test
serve:
	go run main.go

test:
# TODO: Add steps here
# 1. create psql container
# 2. run container
# 3. create database 
	go test -v ./handlers -count=1


# Redis
run_dev_redis:
	docker run --name redisdev -d -p 127.0.0.1:6000:6379 redis:7.2.4 redis-server --requirepass redisdev

start_redis_dev:
	docker start redisdev

stop_redis_dev:
	docker stop redisdev

exec_redis_dev:
	docker exec -it redisdev redis-cli -a redisdev
# mysql old code
# migrate:
# 	migrate -path lifeStyleBack/db/migration -database "mysql://root:Goldendragon1375@tcp(0.0.0.0:3306)/lifestyle" -verbose up


# Fly io
fly_restartbackend:
	fly apps restart lifestylebackend

fly_restartfrontend:
	fly apps restart lifestylefrontend


fly_backlog:
	fly logs -a lifestylebackend

fly_frontlog:
	fly logs -a lifestylefrontend




