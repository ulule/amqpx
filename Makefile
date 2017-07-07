test:
	@(scripts/test)

lint:
	@(scripts/lint)

ci:
	@(go get -d -u github.com/stretchr/testify/require)
	@(go get -d -u github.com/pkg/errors)
	@(go get -d -u github.com/streadway/amqp)
	@(go test -v -race .)
