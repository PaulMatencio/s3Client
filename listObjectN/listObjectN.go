
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
		location string
		prefix string
		limit int
	)

	/* check input parameters */
	flag.StringVar(&bucketName,"b","","-b bucketName")
	flag.StringVar(&location,"s","site1","-s locationName")
	flag.StringVar(&prefix,"prefix","","-prefix prefixName")
	flag.IntVar(&limit,"limit",100,"-limit number")
	flag.Parse()

	if len(bucketName) == 0  {
		flag.Usage()
		log.Fatalln(errors.New("Bucket name is missing"))
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

	// List 'N' number of objects from a bucket-name with a matching prefix.

	listObjectsN := func(bucket, prefix string, recursive bool, N int) (objsInfo []minio.ObjectInfo, err error) {
		// Create a done channel to control 'ListObjects' go routine.
		doneCh := make(chan struct{}, 1)

		// Free the channel upon return.
		defer close(doneCh)

		i := 1
		for object := range s3client.ListObjects(bucket, prefix, recursive, doneCh) {
			if object.Err != nil {
				return nil, object.Err
			}
			i++
			// Verify if we have printed N objects.
			if i == N {
				// Indicate ListObjects go-routine to exit and stop
				// feeding the objectInfo channel.
				doneCh <- struct{}{}
			}
			objsInfo = append(objsInfo, object)
		}
		return objsInfo, nil
	}

	// List recursively first n entries for prefix 'my-prefixname'.
	recursive := true
	objsInfo, err := listObjectsN(bucketName,prefix, recursive, limit)
	if err != nil {
		fmt.Println(err)
	}

	for r := range objsInfo{
		object := objsInfo[r]
		fmt.Printf("Name:%s Size:%d  ContentType:%s\n  Metadata %s\n" ,object.Key,object.Size,object.ContentType,object.Metadata)
	}
}