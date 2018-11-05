
package main

import (
	"errors"
	"flag"
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

	s3Login := s3Client.LoginS3(s3Config,location)
	minioc := s3Login.GetS3Client()  // get minio s3Client

	/*  Create a bucket at location*/
	err = minioc.MakeBucket(bucketName, location)
	if err != nil {
		log.Fatalln(err)
	}

	printOk(s3Login.GetS3Config())

}

