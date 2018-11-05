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
	bucketName 	string
	location 	string
	endpoint 	string
	site1 		s3Client.Host
	ssl 		bool
	start       time.Time
	N           int  = 100
)

func main() {

	type Response struct {
		Filename string
		Err      error
	}

	flag.StringVar(&bucketName, "b", "", "-b bucketName")
	flag.StringVar(&location, "s", "site1", "-s locationName")
	flag.Parse()
	if len(bucketName) == 0  {
		log.Fatalln(errors.New("bucketName is missing"))
	}

	/* get the Config */
	s3Config, err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	/* login to S3  */
	s3Login := s3Client.LoginS3(s3Config,location)
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
	for objInfo := range minioc.ListObjects(bucketName, prefix, true, doneCh) {
		if objInfo.Err != nil {
			fmt.Println(objInfo.Err)
			return
		}
		filenames = append(filenames,objInfo.Key)
	}

	start0 := time.Now()
	messages := make(chan Response)
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
					err = minioc.RemoveObject(bucketName,filename)
					if err != nil {
						s3Client.PrintNotOk(filename, err)
					}
				} else {
					s3Client.PrintNotOk(filename, err)
				}
				messages <- Response{filename, err}
			}(filename)
		}
	if (N == 0 ) {
		log.Printf("Bucket %s was empty",bucketName)
		return
	}
	/*  wait  until all remove are done  */
	for {

		select {
				case /* r:= */ <-messages:
					{
						// fmt.Println(r.Filename,T,N)
						if T == N {
							log.Printf("Remove %d objects from %s in %s\n", N, bucketName,time.Since(start0))
							return
						}
						T++
					}
				case <-time.After(50 * time.Millisecond):
					fmt.Printf("w")
				}
		}

}
