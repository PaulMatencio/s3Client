
package main

import (
	"flag"
	"fmt"
	"github.com/minio/minio-go"
	"github.com/s3Client/lib"
	"log"
	"errors"
)

var (
	bucketName 	string
	location 	string
	endpoint 	string
	site1 		s3Client.Host
	ssl 		bool
)

func printOk() {
	http:= "http"
	if ssl {
		http="https"
	}
	log.Printf("bucket %s is removed\n from %s//%s ", bucketName, http, endpoint)
}

func main() {
	flag.StringVar(&bucketName,"b","","-b bucketName")
	flag.StringVar(&location,"s","site1","-s locationName")
	flag.Parse()
	if len(bucketName) == 0 {
		flag.Usage()
		log.Fatalln(errors.New("Bucket is missing"))
	}

	/* get Config */
	s3Config,err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	/* create an S3 session */
	s3:= s3Client.SetS3Session(s3Config,location)
	s3client, err := minio.New(s3.Endpoint, s3.AccessKeyID, s3.SecretKey, s3.SSL)
	if err != nil {
		log.Fatalln(err)
	}

	err = s3client.RemoveBucket(bucketName)
	if err != nil {
		fmt.Println(err)
		return
	}

	printOk()

}

