package tools

import (
  "encoding/json"
  "net/http"
  "crypto/sha256"
  "fmt"
  "strconv"
  )

type JsonReply struct {
  Status string
  Uid    string
  Sid    string
}

func Hash256(input string) (output string){
  hash := sha256.New()
  bytepass:=[]byte(input)
  hash.Write(bytepass)
  sum := hash.Sum(nil)
  return fmt.Sprintf("%x",sum)
}

func StringToInt64(input string) (output int64){
  output, _ = strconv.ParseInt(input, 10, 64)
  return output
}

func Int64ToString(input int64) (output string){
  return strconv.FormatInt(input, 10)
}

func SendJson(w *http.ResponseWriter, r *http.Request, message string, uid string, sid string) {

  var jsreply JsonReply
  jsreply.Status = message
  jsreply.Sid = sid
  jsreply.Uid = uid
  js, err := json.Marshal(jsreply)
  if err != nil {
    http.Error((*w), err.Error(), http.StatusInternalServerError)
    return
  }
  (*w).Header().Add("Content-Type", "application/json")
  (*w).Header().Add("Access-Control-Allow-Origin", "*")
  (*w).Header().Add("X-Requested-With", "XMLHttpRequest")
  (*w).Write(js)
  return
}
