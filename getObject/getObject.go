
package main

import (
	"errors"
	"flag"
	"github.com/minio/minio-go"
	"github.com/s3Client/lib"
	"log"
	"net/http"
	"time"
	"github.com/moses/user/files/lib"
)

func main() {
	var (
		bucket 			string              /* Bucket name  */
		location 		string              /* S3 location */
		object 			string              /* Object name */
		filename        string              /* file name  */
		trace			bool
	)

	flag.StringVar(&bucket,"b","",s3Client.ABUCKET)
	flag.StringVar(&location,"s","site1",s3Client.ALOCATION)
	flag.StringVar(&object,"o","",s3Client.AOBJECT)
	flag.BoolVar(&trace,"t",false,s3Client.TRACEON)
	flag.StringVar(&filename,"f","",s3Client.AFILE)

	flag.Parse()
	if len(bucket) == 0  || len(object) == 0 {
		flag.Usage()
		log.Fatalln(errors.New("bucketName or objectName cannot be empty"))
	}

	s3Client.TRACE = trace
	s3Config,err := s3Client.GetConfig("config.json")  // get S3 Config
	if err != nil {
		log.Fatalln(err)
	}

	s3Login := s3Client.New(s3Config,location)
	minioc := s3Login.GetS3Client()                            // get minio s3Client after the login  to  set transport option
	tr := &http.Transport{
		DisableCompression: true,
	}
	minioc.SetCustomTransport(tr)


	start := time.Now()
	r := s3Client.S3GetRequest{}
	options := &minio.GetObjectOptions{}
	// options.SetRange(0,10)
	r.S3BuildGetRequest(&s3Login,  bucket,  object,  options)

	buf,err := r.GetObject()

	if err != nil {
		log.Fatalln("Get Key %s error %v",object,err)
	}
	if len(filename) > 0 {
		files.WriteFile(filename,buf.Bytes(),s3Client.FILEMODE)
	}
	log.Printf("Total Duration:  %s  buffer size: %d ",time.Since(start),buf.Len())
}
