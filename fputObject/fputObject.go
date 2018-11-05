
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
		bucketName 	string /* Bucket name  */
		location   	string /* S3 location */
		filename   	string /* output file */
		objectName 	string /* Object name */
		trace		bool
	)

	flag.StringVar(&bucketName, "b", "", "-b bucketName")
	flag.StringVar(&location, "s", "site1", "-s locationName")
	flag.StringVar(&objectName, "o", "", "-o objectName")
	flag.StringVar(&filename, "fn", "", "-fn filename")
	flag.BoolVar(&trace,"trace",false,"-trace ")

	flag.Parse()
	if len(bucketName) == 0 || len(filename) == 0 {
		flag.Usage()
		log.Fatalln(errors.New("bucketName or filename cannot be empty"))
	}

	if len(objectName) == 0 {
		sl := strings.Split(filename, "/")
		objectName = sl[len(sl)-1]
	}

	s3Client.TRACE =trace

	/* get Config */
	s3Config, err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	s3Login := s3Client.LoginS3(s3Config,location)
	minioc := s3Login.GetS3Client()  	// get minio s3Client

	minioc.SetCustomTransport(s3Client.TR)
	opts := minio.PutObjectOptions{}
	opts.ContentType="application/octet-stream"
	start := time.Now()
	if  _,err := minioc.FPutObject(bucketName, objectName, filename,opts); err != nil {
		log.Fatalln(err)
	}
	log.Printf("Duration:  %s", time.Since(start))
}