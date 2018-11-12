
package main

import (
	"flag"
	"github.com/s3Client/lib"
	"log"
)

var (
	location 	string
	trace		bool
)

func main() {

	flag.StringVar(&location,"s","site1",s3Client.ALOCATION)
	flag.BoolVar(&trace,"t",false,s3Client.TRACEON)
	flag.Parse()

	s3Client.TRACE = trace

	s3Config,err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	s3Login := s3Client.New(s3Config,location)

	// s3Login := New.S3Login(location)
	minioc := s3Login.GetS3Client()  // get minio s3Client


	buckets, err := minioc.ListBuckets()
	if err != nil {
		log.Fatalln(err)
	}
	for _, bucket := range buckets {
		log.Println(bucket)
	}
}