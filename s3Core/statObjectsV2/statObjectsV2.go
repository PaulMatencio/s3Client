package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/minio/minio-go"
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
		k    		string
		ObjInfo		minio.ObjectInfo
		Err      	error
	}

	/* check input parameters */
	flag.StringVar(&bucket,"b","",s3Client.ABUCKET)
	flag.StringVar(&location,"s","site1",s3Client.ALOCATION)
	flag.StringVar(&prefix,"p","",s3Client.APREFIX)
	flag.StringVar(&after,"a","","-after <aString>")
	flag.IntVar(&limit,"m",100,s3Client.AMAXKEY)
	flag.BoolVar(&loop,"l",false,"-l loop over to the end ")
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
		disable trace
		enable trace  list http requests
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

		if results,err := s3Core.ListObjectsV2(s3r) ; err == nil {
			t,start,n := 0,time.Now(),len(results.Contents)
			if n==0 {
				goLog.Info.Printf("Bucket %s is empty",bucket)
				return
			}
			// stats n objects returned by previous list Objects
			for _,v := range results.Contents {
				//  Objects removal are started concurrently using go routine
				go 	func(b string, k string ) {
					/* to be done */
					opts:= minio.StatObjectOptions{}
					objInfo, err := s3r.MinioC.StatObject(bucket, k,opts)
					ch <- Response{k,objInfo,err}
				}(bucket,v.Key)
			}

			// wait for the completion of the  removal of n  objects

			for {
				select {
				case  r:= <-ch:
					{   t++
						if r.Err == nil {
							m1 := s3Client.ExtractUserMeta(r.ObjInfo.Metadata)
							if len(m1) > 0 {
								if usermd, err := json.Marshal(m1); err == nil {
									goLog.Info.Printf("user metadata %s of key %s/content -type %s :",r.k,usermd,r.ObjInfo.ContentType)
								} else {
									goLog.Error.Printf("Error parsing user metadata of %s - %v",r.k,err)
								}
							} else {
								goLog.Warning.Printf("key %s - content type %s has no user metadata ",r.k,r.ObjInfo.ContentType)
							}
						} else{
							goLog.Error.Printf("Error retrieving user metadata of %s - %v",r.k,err)
						}
						if t == n {  // all objects are removed
							goLog.Info.Printf("%d user metadata  have been retrieved from  bucket:%s in %s\n", n, bucket,time.Since(start))
							goto Next  // Exit the for select loop
						}

					}
				case <-time.After(50 * time.Millisecond):
					fmt.Printf("w")
				}
			}

		Next:
			N +=n

			goLog.Info.Printf("Is truncated  ? %v - After: %s - Next: %s\n", results.IsTruncated, results.StartAfter,results.NextContinuationToken)

			/*
				Continue the next batch of object stats if loop and list returned was truncated
			 */
			if results.IsTruncated && loop {
				s3r.SetStartAfter(results.NextContinuationToken)
			} else {
				goLog.Info.Printf("Retrieve user metadata of %d objects from bucket '%s' in %v \n",N,bucket,time.Since(start0))
				return
			}

		} else {
			log.Fatalf("List error %v\n",err)
		}
	}

}
