
package main

import (
	"errors"
	"flag"
	"github.com/minio/minio-go"
	"github.com/s3Client/lib"
	"log"
	"strings"
	"time"
)

func main() {
	var (
		bucket	string /* Bucket name  */
		location   	string /* S3 location */
		filename   	string /* output file */
		object 	string /* Object name */
		trace		bool
	)

	flag.StringVar(&bucket, "b", "", s3Client.ABUCKET)
	flag.StringVar(&location, "s", "site1", s3Client.ALOCATION)
	flag.StringVar(&object, "o", "", s3Client.AOBJECT)
	flag.StringVar(&filename, "f", "", s3Client.AFILE)
	flag.BoolVar(&trace,"t",false,s3Client.TRACEON)

	flag.Parse()
	if len(bucket) == 0 || len(filename) == 0 {
		flag.Usage()
		log.Fatalln(errors.New("bucketName or filename cannot be empty"))
	}

	if len(object) == 0 {
		sl := strings.Split(filename, "/")
		object = sl[len(sl)-1]
	}

	s3Client.TRACE =trace

	/* get Config */
	s3Config, err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	s3Login := s3Client.New(s3Config,location)
	minioc := s3Login.GetS3Client()  	// get minio s3Client

	minioc.SetCustomTransport(s3Client.TR)
	opts := minio.PutObjectOptions{}
	opts.ContentType="application/octet-stream"
	start := time.Now()
	if  _,err := minioc.FPutObject(bucket, object, filename,opts); err != nil {
		log.Fatalln(err)
	}
	log.Printf("Duration:  %s", time.Since(start))
}