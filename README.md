# Install go di ubuntu

  sudo apt update

  sudo apt install golang

  go version

  export GOPATH=$HOME/go
  export PATH=$PATH:$GOPATH/bin

  mkdir -p $HOME/go/{bin,pkg,src}


# Create file hello.go

  package main
  
  import "fmt"
  
  func main() {
      fmt.Println("Hello, World!")
  }

# Run go

go run hello.go

# Login_golang
Api login dan regiater dengan language golang dan database postgres

# Tambahkan create file gifnoc.yml
copy yang di bawah ke dalam gifnoc.yml

server:
  port: 9990
  ssl_port: 9990
  hostname: "localhost"

database:
  dbhost: "localhost"
  dbport: "5432"
  dbname: "login_golang"
  dbuser: "login_golang"
  dbpassword: "login_golang"

# Api bisa di test 
1. http://localhost:9990/login
   POST
   Param:
   {
    "username" : "dhani123",
    "password" : "1234567890"
   }
2. http://localhost:9991/register
   POST
   Param:
   {
    "name" : "Muhammad Hamdani",
    "username" :"dhani123",
    "password" : "1234567890",
    "email" : "dhani@gmail.com",
    "phone" : "0813XXXXX"
   }
