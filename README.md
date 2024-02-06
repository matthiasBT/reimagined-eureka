# reimagined-eureka

Secure secrets storage

Run a PostgreSQL instance:  
```
docker run --name reimagened-eureka-postgres -e POSTGRES_PASSWORD=mypassword -d -p 5432:5432 postgres:16.1
```

Run a server instance
```
./bin/server-macos-arm64 -d postgresql://postgres:mypassword@localhost:5432 -a localhost:8888
```

Run a client:  
```
./bin/client-macos-arm64 -a http://localhost:8888 -d /Users/temur.uzbekov/GolandProjects/reimagined-eureka/client.db
```