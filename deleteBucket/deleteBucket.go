
package main

import (
	"errors"
	"flag"
	"github.com/s3Client/lib"
	"log"
)

var (
	bucket	string
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
	log.Printf("bucket %s is removed\n from %s//%s ", bucket, http, endpoint)
}

func main() {
	flag.StringVar(&bucket,"b","",s3Client.ABUCKET)
	flag.StringVar(&location,"s","site1",s3Client.ALOCATION)
	flag.Parse()
	if len(bucket) == 0 {
		flag.Usage()
		log.Fatalln(errors.New("Bucket is missing"))
	}

	/* get s3 Config */
	s3Config,err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	/* login to s3 */
	s3Login := s3Client.New(s3Config,location)
	minioc := s3Login.GetS3Client()  			      	// get s3Client

	if err = minioc.RemoveBucket(bucket);err != nil {
		log.Println(err)
		return
	}

	printOk()

}

