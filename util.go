package tools

type JsonReply struct {
  Status string
  Data interface{}
}

type Cookie struct {
  UID string
  SID string
}

type LoggedInUsers struct {
	UID    int64
	SID 	int64
	Extime    int64
}

type Users struct {
	Username string
	Password string
 	UID int64
}

type OpenID struct {
	Openid string
	UID int64
}

type ProfileVcard struct {
	Attribute  string
	Value      string
	Permission int
	OwnerID    int64
	SearchTerm string
}

type UserProfile struct {
  Username  string
  Data  []ProfileVcard
}

type UnvalidatedUsers struct {
	Email  string
  Code  string
	UID int64
}
