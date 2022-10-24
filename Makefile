docker-build: download-migrate
	mkdir -p ./tmp/docker
	mv ./migrate.linux-amd64 ./tmp/docker/migrate
	cp -r ./pkg/db/table ./tmp/docker
	docker build -t user-syncer:latest .

download-migrate:
	tar zxvf ./pkg/db/migrate/migrate.linux-amd64.tar.gz

