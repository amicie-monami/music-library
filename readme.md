# Installation
Clone this repository using the command. Install dependencies
```
git clone git@github.com:amicie-monami/music-library.git
cd/music-library
go mod tidy
```
Modify the configuration .env file, create a local database instance. Start the server using the command 
```
go run cmd/main.go
```
Use commands to add test values to the database (psql cli required)
```
cd/migrations
psql -d $PG_DATABASE -U $PG_USER -f test_data.sql
```
Once the server is running, you can view the 'open api' (Swagger) documentation at
```
localhost:8080/swagger/
```