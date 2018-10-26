
package main

import (
	"flag"
	"errors"
	"github.com/minio/minio-go"
	"github.com/s3Client/lib"
	"log"
)

var (
	bucketName 	string
	location 	string

)

func printOk(s3 s3Client.S3) {
	http:= "http"
	if s3.SSL {
		http="https"
	}
	log.Printf("bucket %s is created\n at endpoint %s//%s ", bucketName, http, s3.Endpoint)
}

func main() {
	flag.StringVar(&bucketName,"b","","-b bucketName")
	flag.StringVar(&location,"s","site1","-s locationName")
	flag.Parse()
	if len(bucketName) == 0 {
		flag.Usage()
		log.Fatalln(errors.New("Bucket name is missing"))
	}

	/* get Config */
	s3Config,err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	/* create an S3 session */
	s3:= s3Client.SetS3Session(s3Config,location)
	s3client, err := minio.New(s3.Endpoint, s3.AccessKeyID, s3.SecretKey,s3.SSL)
	if err != nil {
		log.Fatalln(err)
	}
	/*  Create a bucket at location*/
	err = s3client.MakeBucket(bucketName, location)
	if err != nil {
		log.Fatalln(err)
	}

	printOk(s3)

}

