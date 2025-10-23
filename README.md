# Golang microservices using RabbitMQ

feat: producer & consumer microservices using separate config.yaml: take 3

Credits to github.com/OtusGolang/webinars_practical_part/29-queues


# install
```cd deployments && sudo docker compose up -d```

```cd ../```
```export DATABASE_URL="postgres://otus_user1:otus_password1@192.168.1.71:5433/events?sslmode=disable"```
```goose -dir migrations postgres $DATABASE_URL up```
```goose -dir migrations postgres $DATABASE_URL up```


# Usage
```go run cmd/calendar/*.go --config=configs/config.yaml```

```go run cmd/producer/*.go --config=configs/producer_config.yaml```

```go run cmd/consumer/*.go --config=configs/consumer_config.yaml```

# Github

```git add .
git commit -m "feat: xxx: take xx"
git push -u origin main```

# Run Producer Once a Day (via Cron / systemd timer)
```
0 7 * * * /usr/local/bin/producer --config=/etc/calendar/producer.yaml
``` 