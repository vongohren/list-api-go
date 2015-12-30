package main

import(
  "log"
  "fmt"
)

func main(){
  server := NewServer(":8000")
  err := server.ListenAndServe()
  if err != nil {
    fmt.Println("error");
    log.Fatalln("Error: %v", err)
  }
}
