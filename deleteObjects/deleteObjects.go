
package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/s3Client/lib"
	"log"
)

func main() {

	var (
		bucket 	string
		location   	string
		prefix     	string
		trace		bool
	)

	/* Delete objects with prefix
	   emptyBucket   does not have -prefix

	*/

	flag.StringVar(&bucket, "b", "", s3Client.ABUCKET)
	flag.StringVar(&location, "s", "site1", s3Client.ALOCATION)
	flag.StringVar(&prefix, "p", "", s3Client.APREFIX)
	flag.BoolVar(&trace,"t",false,s3Client.TRACEON)
	flag.Parse()

	if len(bucket) == 0 {
		flag.Usage()
		log.Fatalln(errors.New("Bucket is missing"))
	}

	s3Client.TRACE=trace

	/* get config  */
	s3Config, err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	s3Login := s3Client.New(s3Config,location)
	minioc := s3Login.GetS3Client()

	objectsCh := make(chan string)

	// Send object names that are needed to be removed to objectsCh
	go func() {
		defer close(objectsCh)
		// List all objects from a bucket-name with a matching prefix.
		for object := range minioc.ListObjects(bucket, prefix, true, nil) {
			if object.Err != nil {
				log.Fatalln(object.Err)
			}
			objectsCh <- object.Key
		}
	}()

	for rErr := range minioc.RemoveObjects(bucket, objectsCh) {
		fmt.Println("Error detected during deletion: ", rErr)
	}

}