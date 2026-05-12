# SendTheSong Go Version

This folder contains the Go rewrite of SendTheSong using Gin + MySQL/MariaDB.

## Run

Set the database connection string:

```bash
set SENDTHESONG_DSN=root:@tcp(127.0.0.1:3306)/sendthesong?parseTime=true&loc=Local
set PORT=8080
```

Then run:

```bash
cd goapp
go mod tidy
go run .
```

Open:
- http://localhost:8080/
- http://localhost:8080/submit
- http://localhost:8080/browse

## Notes

- Static assets are reused from the existing root `assets/` folder.
- Uploaded files continue to be served from the root `uploads/` folder.
- The schema must include the `song_*` columns already used by the app.
