Make sure use install cockroach db in your machine

```bash
brew install cockroach

# start cockroach 
cockroach start --insecure --host=localhost

# create database
cockroach sql --insecure

> CREATE database gochat;
> GRANT ALL ON database gochat to roach1;

# running go code
go get
go run main.go
```



