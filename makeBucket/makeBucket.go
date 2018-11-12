
package main

import (
	"errors"
	"flag"
	"github.com/s3Client/lib"
	"log"
)

var (
	bucket 	string
	location 	string

)

func printOk(s3 s3Client.S3) {
	http:= "http"
	if s3.SSL {
		http="https"
	}
	log.Printf("bucket %s is created\n at endpoint %s//%s ", bucket, http, s3.Endpoint)
}

func main() {
	flag.StringVar(&bucket,"b","",s3Client.ABUCKET)
	flag.StringVar(&location,"s","site1",s3Client.ALOCATION)
	flag.Parse()
	if len(bucket) == 0 {
		flag.Usage()
		log.Fatalln(errors.New("Bucket is missing"))
	}

	/* get Config */
	s3Config,err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	s3Login := s3Client.New(s3Config,location)
	minioc := s3Login.GetS3Client()  // get minio s3Client

	/*  Create a bucket at location*/
	err = minioc.MakeBucket(bucket, location)
	if err != nil {
		log.Fatalln(err)
	}

	printOk(*s3Login.S3)

}

