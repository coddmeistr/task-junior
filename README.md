
## Installation

  1. Clone


```bash
  git clone https://github.com/maxik12233/task-junior
```
    
2. Create .env file in root directory

3. Write all required environment variables in this file
```bash
  DB_CONNECTION_STRING="postgres://postgres:password@host:port/database?sslmode=disable"
  PORT=4000 (default - 4000)
  IS_DEBUG=false (default - false)
```
Only  ```DB_CONNECTION_STRING``` is required and should contain ```sslmode=disable``` to prevent unnecessary errors.

```IS_DEBUG``` if true, then it enables Gin's logger and recovery, so it should be false.

4. Build binary from ```cmd/main.go``` file and run it.



## API Reference

### Get person/persons

```http
  GET /person
```
#### Allowed queries
```id int``` - if correct, you get one requested instance of this id.

```page int``` - pagination page.

```per_page int``` - how many instances in one page.

```sort_by string``` - specify sort field.

```sort_order string``` - specify sort order (asc or desc).


### Change person instance

```http
  PUT /person
```
Request JSON body schema:
```http
  {
    "id" int,
    "name" string,
    "surname" string,
    "patronymic" string (optional),
    "gender" string,
    "age" int,
    "nationality" string
  }
```



### Delete person instance

```http
  DELETE /person
```
Request JSON body schema:
```http
  {
    "id" int
  }
```

### Create person

```http
  POST /stats
```
Request JSON body schema:
```http
  {
    "name" string,
    "surname" string
    "patronymic" string (optional)
  }
```


