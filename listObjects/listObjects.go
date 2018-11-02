package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/minio/minio-go"
	"github.com/s3Client/lib"
	"log"
	"errors"
	"os"
	"time"
)


func main() {

	var (
		bucketName 	string
		location 	string
		prefix		string
		limit 		int
		trace		bool
	)

	/* check input parameters */
	flag.StringVar(&bucketName,"b","","-b bucketName")
	flag.StringVar(&location,"s","site1","-s locationName")
	flag.StringVar(&prefix,"prefix","","-prefix prefixName")
	flag.IntVar(&limit,"limit",100,"-limit number")
	flag.BoolVar(&trace,"t",false,"-t ")
	flag.Parse()
	if len(bucketName) == 0  {
		flag.Usage()
		log.Fatalln(errors.New("bucketName cannot be empty"))
	}

	/* get config  */
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

	if trace  {
		s3client.TraceOn(os.Stdout)
	}

	// Create a done channel to control 'ListObjects' go routine.
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	// List all objects from a bucket-name with a matching prefix.
	start:= time.Now()
	n:= 0
	for objInfo := range s3client.ListObjects(bucketName, prefix, true, doneCh) {
		if objInfo.Err != nil {
			fmt.Println(objInfo.Err)
			return
		}
		n++
		metadata, _ := json.Marshal(objInfo.Metadata)
		fmt.Printf("key : %s  Content Type: %s Last Modified %s Size: %d  Metadata %s\n",objInfo.Key, objInfo.ContentType,  objInfo.LastModified,  objInfo.Size,metadata)
	}
	fmt.Printf("Listing %d objects in %s\n" ,n,time.Since(start))
	return
}
