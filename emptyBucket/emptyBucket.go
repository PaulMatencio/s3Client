package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/s3Client/lib"
	"log"
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
	N           int  = 100
)

func main() {

	type Response struct {
		Key		 string
		Err      error
	}

	flag.StringVar(&bucket, "b", "", s3Client.ABUCKET)
	flag.StringVar(&location, "s", "site1", s3Client.ALOCATION)
	flag.Parse()
	if len(bucket) == 0  {
		log.Fatalln(errors.New("bucket is missing"))
	}

	/* get the Config */
	s3Config, err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	/* login to S3  */
	s3Login := s3Client.New(s3Config,location)
	minioc := s3Login.GetS3Client()  				// get minio s3Clie
	runtime.GOMAXPROCS(4)


	//  list the Objects of a  bucket
	// Create a done channel to control 'ListObjects' go routine.
	doneCh := make(chan struct{})
	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	// List the buckets witout  prefix

	prefix:=""
	filenames := []string{}
	for objInfo := range minioc.ListObjects(bucket, prefix, true, doneCh) {
		if objInfo.Err != nil {
			fmt.Println(objInfo.Err)
			return
		}
		filenames = append(filenames,objInfo.Key)
	}

	start0 := time.Now()
	ch := make(chan Response)
	runtime.GOMAXPROCS(4)
	N := len(filenames)
	T := 1
	/* */
	// s3client.TraceOn(os.Stdout)
	for obj := 0; obj < N; obj++ {
		start = time.Now()
			filename := filenames[obj]

			go func( string) {

				if err == nil {
					err = minioc.RemoveObject(bucket,filename)
					if err != nil {
						s3Client.PrintNotOk(filename, err)
					}
				} else {
					s3Client.PrintNotOk(filename, err)
				}
				ch <- Response{filename, err}
			}(filename)
		}
	if (N == 0 ) {
		log.Printf("Bucket %s was empty",bucket)
		return
	}
	/*  wait  until all remove are done  */
	for {

		select {
				case /* r:= */ <-ch:
					{
						// fmt.Println(r.Filename,T,N)
						if T == N {
							log.Printf("%d objects have been removed from  bucket:%s in %s\n", N, bucket,time.Since(start0))
							return
						}
						T++
					}
				case <-time.After(50 * time.Millisecond):
					fmt.Printf("w")
		}
	}

}
