
package main

import (
	"errors"
	"flag"
	"github.com/minio/minio-go"
	"github.com/s3Client/lib"
	"log"
	"net/http"
	"time"
	"user/files/lib"
)

func main() {
	var (
		bucketName 		string              /* Bucket name  */
		location 		string              /* S3 location */
		objectName 		string              /* Object name */
		filename        string              /* file name  */
		trace			bool

	)

	flag.StringVar(&bucketName,"b","","-b bucketName")
	flag.StringVar(&location,"s","site1","-s locationName")
	flag.StringVar(&objectName,"o","","-o objectName")
	flag.BoolVar(&trace,"trace",false,"-trace ")
	flag.StringVar(&filename,"fn","","-fn fileName")

	flag.Parse()
	if len(bucketName) == 0  || len(objectName) == 0 {
		flag.Usage()
		log.Fatalln(errors.New("bucketName or objectName cannot be empty"))
	}

	s3Client.TRACE = trace
	s3Config,err := s3Client.GetConfig("config.json")  // get S3 Config
	if err != nil {
		log.Fatalln(err)
	}

	s3Login := s3Client.LoginS3(s3Config,location)
	minioc := s3Login.GetS3Client()                            // get minio s3Client after the login  to  set transport option
	tr := &http.Transport{
		DisableCompression: true,
	}
	minioc.SetCustomTransport(tr)


	start := time.Now()
	r := s3Client.S3Request{}
	options := &minio.GetObjectOptions{}
	// options.SetRange(0,10)
	r.S3BuildGetRequest(&s3Login,  bucketName,  objectName,  options)
	buf,err := s3Client.GetObject(r)

	if err != nil {
		log.Fatalln("Get Key %s error %v",objectName,err)
	}
	if len(filename) > 0 {
		files.WriteFile(filename,buf.Bytes(),0644)
	}
	log.Printf("Total Duration:  %s  buffer size: %d ",time.Since(start),buf.Len())
}
