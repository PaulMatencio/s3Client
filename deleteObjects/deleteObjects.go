
package main

import (
	"flag"
	"fmt"
	"github.com/minio/minio-go"
	"github.com/s3Client/lib"
	"log"
	"errors"
)

func main() {

	var (
		bucketName string
		location   string
		prefix     string
	)

	/* check input parameters */
	flag.StringVar(&bucketName, "b", "", "-b bucketName")
	flag.StringVar(&location, "s", "site1", "-s locationName")
	flag.StringVar(&prefix, "prefix", "", "-prefix prefixName")

	flag.Parse()
	if len(bucketName) == 0 {
		flag.Usage()
		log.Fatalln(errors.New("Bucket is missing"))
	}

	/* get config  */
	s3Config, err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	site1:= s3Client.StructToMap(&s3Config.Hosts)[location]
	endpoint := site1.GetUrl()
	accessKeyID := site1.GetAccesKey()
	secretAccessKey := site1.GetSecretKey()
	ssl := site1.GetSecure()

	s3Client, err := minio.New(endpoint, accessKeyID, secretAccessKey,ssl)
	if err != nil {
		fmt.Println(err)
		return
	}


	objectsCh := make(chan string)

	// Send object names that are needed to be removed to objectsCh
	go func() {
		defer close(objectsCh)
		// List all objects from a bucket-name with a matching prefix.
		for object := range s3Client.ListObjects(bucketName, prefix, true, nil) {
			if object.Err != nil {
				log.Fatalln(object.Err)
			}
			objectsCh <- object.Key
		}
	}()

	for rErr := range s3Client.RemoveObjects(bucketName, objectsCh) {
		fmt.Println("Error detected during deletion: ", rErr)
	}

}