package tools

import (
  "encoding/base64"
  "regexp"
  "net/http"
  "crypto/sha256"
  "fmt"
  "strconv"
  "appengine/datastore"
  "appengine"
  //"github.com/gorilla/securecookie"
  "github.com/martini-contrib/sessions"
  "time"
  "crypto/rand"
  )


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

func (jsreply *JsonReply)CreateJson(w *http.ResponseWriter, r *http.Request, message string, object interface{}) {
  (*jsreply).Status = message
  (*jsreply).Data = object
}

func SetSession(userid int64, w *http.ResponseWriter, r *http.Request,session *sessions.Session) {
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

	(*session).Set("UID",strconv.FormatInt(userid, 10))
	(*session).Set("SID",strconv.FormatInt(keyPut.IntID(), 10))

	loginUser:=loggedinusers{
		UID: userid,
		SID : keyPut.IntID(),
		Extime : time.Now().Unix() + 3600,
	}
	if _, errPut := datastore.Put(c, LoginKey(c, loginUser.SID), &loginUser); errPut != nil {
		fmt.Fprint((*w), errPut)
	}
}

func ClearSession(SID int64, w *http.ResponseWriter, r *http.Request,session *sessions.Session) {
	c := appengine.NewContext(r)
	(*session).Delete("UID")
  (*session).Delete("SID")
	if deleteErr := datastore.Delete(c, LoginKey(c, SID)); deleteErr != nil {
		fmt.Fprint(*w, deleteErr)
	}
}
