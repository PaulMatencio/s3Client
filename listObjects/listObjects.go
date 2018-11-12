package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/s3Client/lib"
	"log"
	"os"
	"time"
)


func main() {

	var (
		bucket	string
		location 	string
		prefix		string
		limit 		int
		trace		bool
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
		log.Fatalln(errors.New("bucket is missing"))
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

	// Create a done channel to control 'ListObjects' go routine.
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	// List all objects from a bucket-name with a matching prefix.
	start:= time.Now()
	n:= 0
	for objInfo := range minioc.ListObjects(bucket, prefix, true, doneCh) {
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
