package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/moses/user/goLog"
	"github.com/s3Client/lib"
	"github.com/s3Client/s3Core/lib"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"
)


func main() {

	var (
		bucket 	string
		location 	string
		prefix		string
		delimiter   string
		after       string
		limit 		int
		fetchOwner  bool
		trace		bool
		loop        bool
	)

	type Response struct {
		Key		 string
		Err      error
	}

	/* check input parameters */
	flag.StringVar(&bucket,"b","",s3Client.ABUCKET)
	flag.StringVar(&location,"s","site1",s3Client.ALOCATION)
	flag.StringVar(&prefix,"p","",s3Client.APREFIX)
	flag.StringVar(&after,"a","","-after <aString>")
	flag.IntVar(&limit,"m",100,s3Client.AMAXKEY)
	flag.BoolVar(&loop,"l",false,"-l loop ")
	flag.BoolVar(&trace,"t",false,s3Client.TRACEON)
	flag.Parse()
	if len(bucket) == 0  {
		flag.Usage()
		log.Fatalln(errors.New("bucketName cannot be empty"))
	}
	/*
		init logging
	 */
	goLog.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	/*
		get config
	*/
	s3Config,err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	/*
		Create an S3 session
	 */
	s3Login := s3Core.New(s3Config,location)

	/*
		Build a List request with prefix and maxkey ( Limit)
	 */
	s3r := s3Core.S3ListRequest{
		MinioC: s3Login.MinioC,
		Bucket: bucket,
		Prefix: prefix,
		Delimiter: delimiter,
		StartAfter: after,
		// Marker: next,
		Limit: limit,
		FetchOwner:fetchOwner,
	}

	/*
		enable tracing of http requests
	*/
	s3r.Trace = false
	if s3Client.TRACE {
		s3r.Trace	= true
	}

	// s3r.S3BuildListRequest(&s3Login, bucket, prefix, false,delimiter, after, next,limit)
	runtime.GOMAXPROCS(4)
	ch := make(chan Response)
	var N, start0 = 0,time.Now()

	for {

		if results,err := s3Core.ListObjectsV1(s3r) ; err == nil {
			t,start,n := 0,time.Now(),len(results.Contents)
			if n==0 {
				goLog.Info.Printf("Bucket %s was empty",bucket)
				return
			}
			// remove n objects returned by previous list Objects
			var nextMarker string
			for _,v := range results.Contents {
				//  Objects removal are started concurrently using go routine
				go 	func(b string, k string ) {
					err := s3r.MinioC.RemoveObject(b, k)
					ch <- Response{k,err}
				}(bucket,v.Key)
				nextMarker = v.Key
			}

			// wait for the completion of the  removal of n  objects

			for {
				select {
				case  /*r:=*/ <-ch:
					{   t++
						if t == n {  // all objects are removed
							goLog.Info.Printf("%d objects have been removed from  bucket:%s in %s\n", n, bucket,time.Since(start))
							goto Next  // Exit the for select loop
						}

					}
				case <-time.After(50 * time.Millisecond):
					fmt.Printf("w")
				}
			}

		Next:
			N +=n

			goLog.Info.Printf("Is truncated  ? %v - After: %s - Next: %s\n", results.IsTruncated, results.Marker,nextMarker)

			/*
				Continue the next batch of delete if loop and list returned was truncated
			 */
			if results.IsTruncated && loop {
				s3r.SetStartAfter(nextMarker)
			} else {
				goLog.Info.Printf("%d objects were deleted from bucket '%s' in %v \n",N,bucket,time.Since(start0))
				return
			}

		} else {
			log.Fatalf("List error %v\n",err)
		}
	}

}



