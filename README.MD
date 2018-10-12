# GOLANG LIBRARY TOOLS
1. httpclient  
   httpclient  library very simple chaining method & support ssl  
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
                          SetTimeout(3 * time.Second).
                          //SetPemCertificate("cert/cert_rsa.pub"). // if use ssl
                          Get("https://httpstat.us/200?sleep=5000")
      
      if err != nil && rsp.IsTimeout() {
      		fmt.Println(err)
      		return
      }
      
      if err != nil {
      		fmt.Println(err)
      		return
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


### Author  
* Daud Valentino labsbrain@gmail.com 