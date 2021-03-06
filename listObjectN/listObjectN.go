
package main

import (
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
		bucket string
		location string
		prefix string
		limit int
		trace bool
	)

	/* check input parameters */
	flag.StringVar(&bucket,"b","",s3Client.ABUCKET)
	flag.StringVar(&location,"s","site1",s3Client.ALOCATION)
	flag.StringVar(&prefix,"p","",s3Client.APREFIX)
	flag.IntVar(&limit,"m",100,s3Client.AMAXKEY)
	flag.BoolVar(&trace,"t",false,s3Client.TRACEON)
	flag.Parse()

	if len(bucket) == 0  {
		flag.Usage()
		log.Fatalln(errors.New("Bucket is missing"))
	}

	/* get config  */
	s3Config,err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	s3Login := s3Client.New(s3Config,location)
	minioc := s3Login.GetS3Client()  // get minio s3Client
	if trace  {
		minioc.TraceOn(os.Stdout)
	}

	// List 'N' number of objects from a bucket-name with a matching prefix.
	start := time.Now()

	listObjectsN := func(bucket, prefix string, recursive bool, N int) (objsInfo []minio.ObjectInfo, err error) {
		// Create a done channel to control 'ListObjects' go routine.
		doneCh := make(chan struct{}, 1)

		// Free the channel upon return.
		defer close(doneCh)

		i := 1
		for object := range minioc.ListObjects(bucket, prefix, recursive, doneCh) {
			if object.Err != nil {
				return nil, object.Err
			}
			i++
			// Verify if we have printed N objects.
			log.Println(i,N)
			if i == N  {
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

	objsInfo, err := listObjectsN(bucket,prefix, recursive, limit)
	if err != nil {
		fmt.Println(err)
	}

	for r := range objsInfo{
		object := objsInfo[r]
		log.Printf("Name:%s Size:%d Last modified:%v\n" ,object.Key,object.Size,object.LastModified)
	}

	log.Printf("Listing  %d objects in %s\n ",len(objsInfo),   time.Since(start))
}