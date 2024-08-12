all: build

build:
	docker build --platform=linux/amd64 -t cr.yandex/crpd9bbfo3gaist0rdn6/tg-efs:latest .
	docker push cr.yandex/crpd9bbfo3gaist0rdn6/tg-efs:latest