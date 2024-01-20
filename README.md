
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