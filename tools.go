package tools

import (
  "encoding/json"
  "encoding/base64"
  "regexp"
  "net/http"
  "crypto/sha256"
  "fmt"
  "strconv"
  "appengine/datastore"
  "appengine"
  "github.com/gorilla/securecookie"
  "time"
  "crypto/rand"
  )

var key1 = []byte("5916569511133184")
var key2 = []byte("4776259720577024")
var CookieHandler = securecookie.New(key1, key2)

var Appname string = "authentication.auth-test-ryan.appspot.com"
//var Appname string = "127.0.0.1:8081"
var FacebookclientID string = "499628346846146"
var FacebookclientSecret string = "4538c6faccc2ea698392220c210e6d54"
var GoogleClientID string = "123186211988-uhijqudedasocrct7vn1s70jp7srsr68.apps.googleusercontent.com"
var GoogleClientSecret string = "uC8mXLjT5WpxLjpbuGNFboGp"
var LinkedinAPIKey string = "75xf5zzyy9kf8v"
var LinkedinAPISecret string = "8StnOtxpp1ftGZHN"
var AdminEmail = "selva@vidao.co"

type jsonreply JsonReply
type loggedinusers Loggedinusers


func OpenidKey(c appengine.Context, openid string) *datastore.Key {
  return datastore.NewKey(c, "openiduid", openid, 0, nil)
}

func SignUpLockKey(c appengine.Context, openid string) *datastore.Key {
  return datastore.NewKey(c, "SignUpLock", openid, 0, nil)
}

func UsersKey(c appengine.Context) *datastore.Key {
  return datastore.NewKey(c, "User", "", 0, nil)
}

func FillUsersKey(c appengine.Context, Userid int64) *datastore.Key {
	return datastore.NewKey(c, "User", "", Userid, nil)
}


func LoginKey(c appengine.Context, Sessionid int64) *datastore.Key {
	return datastore.NewKey(c, "login", "", Sessionid, nil)
}

func CreateloginKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "login", "", 0, nil)
}

func ProfileKey(c appengine.Context, ancestor *datastore.Key) *datastore.Key {
	return datastore.NewKey(c, "Profile", "", 0, ancestor)
}

func UnvalidatedKey(c appengine.Context, Userid int64) *datastore.Key {
	return datastore.NewKey(c, "unvaildatedusers", "", Userid, nil)
}

func Hash256(input string) (output string){
  hash := sha256.New()
  bytepass:=[]byte(input)
  hash.Write(bytepass)
  sum := hash.Sum(nil)
  return fmt.Sprintf("%x",sum)
}

func NonceGenerator(w *http.ResponseWriter) string {
	size := 32 // change the length of the generated random string here
	rb := make([]byte, size)
	_, err := rand.Read(rb)
	if err != nil {
		fmt.Println(err)
	}
	rs := base64.URLEncoding.EncodeToString(rb)

	reg, err := regexp.Compile("[^A-Za-z0-9]+") //remove any non-alphanumerical character
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
	}

	safe := reg.ReplaceAllString(rs, "a")
	safe=safe[:21]
	return safe
}

func StringToInt64(input string) (output int64){
  output, _ = strconv.ParseInt(input, 10, 64)
  return output
}

func Int64ToString(input int64) (output string){
  return strconv.FormatInt(input, 10)
}

func SendJson(w *http.ResponseWriter, r *http.Request, message string, id string, sid string) {

  var jsreply JsonReply
  jsreply.Status = message
  jsreply.Sid = sid
  jsreply.ID = id
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

func SetSession(userid int64, w *http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	g := &loggedinusers{
		SID: 0,
	}

	key := CreateloginKey(c)
	keyPut, err := datastore.Put(c, key, g)
	if err != nil {
    	http.Error((*w), err.Error(), http.StatusInternalServerError)
		return
	}

	value := map[string]string{
		"userid":    strconv.FormatInt(userid, 10),
		"sessionid": strconv.FormatInt(keyPut.IntID(), 10),
	}
	if encoded, err := CookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{

			Name:  "session",
			Value: encoded,
			Path:  "/",
			MaxAge: 3600,
			Domain: ".auth-test-ryan.appspot.com",


		}

   		http.SetCookie((*w), cookie)

		loginUser:=loggedinusers{
			UID: userid,
			SID : keyPut.IntID(),
			Extime : time.Now().Unix() + 3600,
		}
		if _, errPut := datastore.Put(c, LoginKey(c, loginUser.SID), &loginUser); errPut != nil {
      		fmt.Fprint((*w), errPut)
		}
	}
	//sendJson(&w, r, "User Logged In", "0", strconv.FormatInt(keyPut.IntID(), 10))
}

func ClearSession(SID int64, w *http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	cookie := &http.Cookie{
		Name:   "session",
		MaxAge: -1,
		Domain: ".auth-test-selva.appspot.com",
	}

	//fmt.Fprint(w,"cookies deleted")
	http.SetCookie(*w, cookie)

	if deleteErr := datastore.Delete(c, LoginKey(c, SID)); deleteErr != nil {
		fmt.Fprint(*w, deleteErr)
	}
}