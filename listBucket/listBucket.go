
package main

import (
	"flag"
	"github.com/minio/minio-go"
	"github.com/s3Client/lib"
	"log"
)

var location string

func main() {

	flag.StringVar(&location,"s","site1","-s locationName")
	flag.Parse()

	s3Config,err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	/* create an S3 session */
	s3:= s3Client.SetS3Session(s3Config,location)
	// fmt.Println(s3.Endpoint,s3.AccessKeyID,s3.SecretKey,s3.SSL)
	s3client, err := minio.New(s3.Endpoint, s3.AccessKeyID, s3.SecretKey, s3.SSL)
	if err != nil {
		log.Fatalln("Create Session: ", err)
	}
	buckets, err := s3client.ListBuckets()
	if err != nil {
		log.Fatalln(err)
	}
	for _, bucket := range buckets {
		log.Println(bucket)
	}
}