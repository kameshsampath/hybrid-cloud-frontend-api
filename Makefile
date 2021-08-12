.PHONY swagger
swagger:	
	swag init -g server.go

.PHONY build-image
build-image:	
	docker build --rm --nocache -t ghcr.io/kameshsampath/hybrid-cloud-frontend-api .
	docker push ghcr.io/kameshsampath/hybrid-cloud-frontend-api