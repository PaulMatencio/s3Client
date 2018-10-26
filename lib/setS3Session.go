package s3Client

type S3 struct {
	Site Host         	 `json:"host"`
	Endpoint string    	 `json:"endpoint"`
	AccessKeyID  string  `json:"accessKeyID"`
	SecretKey string     `json:"secretKey"`
	SSL bool             `json:"ssl"`
}

func SetS3Session(s3Config Config,location string) (S3) {
	 var (
	 	s S3
	 	site Host
	 	)
	 site = StructToMap(&s3Config.Hosts)[location]
		s.Site = site
		s.Endpoint = site.URL
		s.AccessKeyID = site.AccessKey
		s.SecretKey = site.SecretKey
		s.SSL = site.SSL
	return s
}
