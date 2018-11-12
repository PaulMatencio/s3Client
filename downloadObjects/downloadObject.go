package main

import (
	"bytes"
	"errors"
"flag"
"fmt"
	"github.com/minio/minio-go"
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
	prefix 		string
	N           int  = 500
)

func main() {

	type Response struct {
		Filename string
		Buffer  *bytes.Buffer
		Err      error
	}

	flag.StringVar(&bucket, "b", "", s3Client.ABUCKET)
	flag.StringVar(&location, "s", "site1", s3Client.ALOCATION)
	flag.StringVar(&prefix, "p", "", s3Client.TRACEON)
	flag.Parse()
	if len(bucket) == 0  {
		flag.Usage()
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

	filenames := []string{}
	for objInfo := range minioc.ListObjects(bucket, prefix, true, doneCh) {
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
	S := 0
	/* */
	// s3client.TraceOn(os.Stdout)
	for obj := 0; obj < N; obj++ {
		start = time.Now()
		filename := filenames[obj]

		go func( string) {

			if err == nil {
				r := s3Client.S3Request{}
				options := &minio.GetObjectOptions{}
				r.S3BuildGetRequest(&s3Login,  bucket,  filename,  options)
				buf,err := s3Client.GetObject(r)
				messages <- Response{filename,buf,err}
			} else {
				messages <- Response{filename,nil, err}
			}

		}(filename)
	}
	if (N == 0 ) {
		log.Printf("Bucket %s is empty",bucket)
		return
	}
	/*  wait  until all remove are done  */
	for {

		select {
		case    r:= <-messages:
			{
				log.Printf("Download Object Key: %s / Object size %d bytes",r.Filename, r.Buffer.Len())

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

