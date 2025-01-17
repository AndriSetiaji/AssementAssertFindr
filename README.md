# AssementAssertFindr
### tech stack:
- Language: Golang ✔️
- Framework: Gin ✔️
- ORM: GORM ✔️
- Database: PostgreSQL ✔️

### Note
## postgressql
- run postgressql in local with docker compose.
- for auth and settings already write on docker-compose.yml
## golang
- in this repo already include go mod and go sum for running
- adjust file .env for connection with postgressql in local
- can direct running with go run [your-path]/src/main.go
- running in localhost:8080

## ERD on PostgresSql
![image](https://github.com/user-attachments/assets/8f1831c1-51b2-4d08-8a50-35e6da3ee62f)

## API get list all data - GET /api/post 
```
curl --location 'http://localhost:8080/api/post'
```
![image](https://github.com/user-attachments/assets/c24441f8-5020-4f34-8ad5-88f98ddb59f1)


## API add data - POST /api/post/
```
curl --location 'http://localhost:8080/api/post' \
--header 'Content-Type: application/json' \
--data '{
    "title": "title test new 001",
    "content":"content test new 001",
    "tags":["add new 001", "add new 001-A"]
}'
```
![image](https://github.com/user-attachments/assets/63cd3287-55d9-4d79-bd33-35714c8056d9)


## API get data by id - GET /api/post/:postId
```
curl --location 'http://localhost:8080/api/post/26'
```
![image](https://github.com/user-attachments/assets/028a8a32-24c8-43ce-8117-10d3f822829f)


## API update data - PUT /api/post/
```
curl --location --request PUT 'http://localhost:8080/api/post/26' \
--header 'Content-Type: application/json' \
--data '{
    "title": "request title new update test 01",
    "content":"request content new update test 01",
    "tags":["request tag new update test 01", "request tag new update test 02", "request tag new update test 03"]
}'
```
![image](https://github.com/user-attachments/assets/8454d539-1288-46cc-954e-acfb4a9b8492)

## API delete data - DELETE /api/post/:postId
```
curl --location --request DELETE 'http://localhost:8080/api/post/26'
```
![image](https://github.com/user-attachments/assets/c7c9da91-dbf1-4432-9c0f-c83c3d17aa5c)


