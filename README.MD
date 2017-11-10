#GOLANG LIBRARY TOOLS
1. httpclient  
   httpclient  library very simple chaining method adn support ssl  
   **example:**
  ```go 
  
  package main
  
  import (
        "encoding/json"
        "time"
        
        "brainlabs/httpClient"
  )
  
  func main(){
        
      client := httpClient.NewClient()      
      rsp, err := client.
                          SetHeader("Content-Type", "application/json").
                          SetHeader("Channel", "tester").
                          SetTimeout(10 * time.Second).
                          //SetPemCertificate("cert/cert_rsa.pub"). // if use ssl
                          Get("https://github.com")
      
      if err != nil {
          fmt.Println(err)
      }
      
      // get raw response
      fmt.Println(rsp.GetRaw())
      
     
      var testData = struct {
      		ID   string `json:"id"`
      		Name string `json:"name"`
      }{}
      
      err := rsp.GetUnmarshalJSON(&testData)
      
      if err != nil {
      		fmt.Println(err)
      }
      
      x, _ := json.MarshalIndent(testData, "", " ")
      
      fmt.Println(string(x))
  
  }
  
  ```
