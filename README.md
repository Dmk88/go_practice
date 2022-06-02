# Currency monitoring service

run docker-compose:
```
docker-compose up -d
```
create db:
```
./currencymonitor/create-tables.sh
```
url example - start monitoring 
```
http://localhost:8080/start?period=10m&frequency=1m
```
url example - getting results 
```
http://localhost:8080/results?monitoring_id=452d97f5-0a17-4179-b39c-f9734999600e
```
period and frequency - Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".