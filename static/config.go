package static

const (
	ServerName                      = "cwmServer"
	ServerEventName                 = "cwmServerEvent"
	MaxResendAuthenCode             = 3
	TimeResetPhoneCreateAccountLock = 60 //60 mins
	AuthencodeTimeOut               = 60 //60 sec
	MaxNumberAuthencodeFail         = 3
	NumberPrefix                    = "911"
	TimeToReCreateNONCE             = 5  //5 mins
	NONCETTL                        = 20 //60 mins
	JWTTTL                          = 20 //60 mins
	//MaxFileSize = 1 << 10 //1 KB
	MaxFileSize   = 150 * 1024 << 10 //150 MB
	S3NamePrefxix = "cwm_ttl_"
)

func JWTKey() []byte {
	return []byte("CWM JWT KEY")
}
