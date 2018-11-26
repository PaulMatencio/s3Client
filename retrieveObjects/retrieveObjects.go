
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/minio/minio-go"
	"github.com/moses/user/goLog"
	"github.com/s3Client/lib"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"
)

var (
	bucket 	string
	location 	string
	endpoint 	string
	site1 		s3Client.Host
	ssl 		bool
	start       time.Time
	prefix 		string
	trace		bool
)

type Response struct {
	Key 	string
	Buffer  *bytes.Buffer
	Err      error
}

func main() {

	flag.StringVar(&bucket, "b", "", s3Client.ABUCKET)
	flag.StringVar(&location, "s", "site1", s3Client.ALOCATION)
	flag.StringVar(&prefix, "p", "", s3Client.APREFIX)
	flag.BoolVar(&trace, "t", false, s3Client.TRACEON)
	flag.Parse()
	if len(bucket) == 0  {
		flag.Usage()
		log.Fatalln(errors.New("bucket is missing"))
	}

	/*
		init logging
	 */
	goLog.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	/* get the Config */
	s3Config, err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	/* login to S3  */
	s3Login := s3Client.New(s3Config,location)
	minioc := s3Login.GetS3Client()  				// get minio s3 client
	runtime.GOMAXPROCS(4)

	// list the Objects of a  bucket
	// Create a done channel to control 'ListObjects' go routine.
	doneCh := make(chan struct{})
	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	// List the buckets with a  prefix

	keys := []string{}
	for objInfo := range minioc.ListObjects(bucket, prefix, true, doneCh) {
		if objInfo.Err != nil {
			fmt.Println(objInfo.Err)
			return
		}
		keys = append(keys,objInfo.Key)
	}

	start0 := time.Now()
	N,T,S := len(keys),1,0
	if (N == 0) {
		log.Printf("Bucket %s is empty",bucket)
		return
	}
	ch := make(chan Response)
	runtime.GOMAXPROCS(4)

	if trace {
		minioc.TraceOn(os.Stdout)
	}

	for obj := 0; obj < N; obj++ {
		start = time.Now()
		key := keys[obj]
		go func( string) {

			r := s3Client.S3GetRequest{}
			//  todo Get Options
			options := &minio.GetObjectOptions{}
			r.S3BuildGetRequest(&s3Login,  bucket,  key,  options)
			buf,err := r.GetObject()
			ch <- Response{key,buf,err}

		}(key)
	}


	/*  wait  until all objects are retrieved  */
	for {

		select {
		case    r:= <-ch:
			{
				log.Printf("Download Object Key: %s / Object size %d bytes",r.Key, r.Buffer.Len())

				if T == N {
					log.Printf("Downloaded %d objects ( Total size %d bytes)  from %s in %s\n", N, S, bucket,time.Since(start0))
					return
				}
				T++
				S +=r.Buffer.Len()
			}
		case <-time.After(50 * time.Millisecond):
			fmt.Printf("w")
		}
	}

}

