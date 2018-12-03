package main

import (
	"bytes"
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
	"strings"
	"time"
)

type Response struct {
	K		 string
	Buffer  *bytes.Buffer
	ObjInfo  *minio.ObjectInfo
	Err      error
}



func main() {

	var (
		bucket 		string
		outdir		string
		location 	string
		prefix		string
		delimiter   string
		marker      string
		limit 		int
		trace		bool
		loop        bool
	)
	/*
	type Response struct {
		k    		string
		ObjInfo		minio.ObjectInfo
		Err      	error
	}
	*/

	/* define  input parameters */
	flag.StringVar(&bucket,"b","",s3Client.ABUCKET)
	flag.StringVar(&outdir,"o","","-o output directory")
	flag.StringVar(&location,"s","site1",s3Client.ALOCATION)
	flag.StringVar(&prefix,"p","",s3Client.APREFIX)
	flag.StringVar(&marker,"marker","","-marker aString")
	flag.StringVar(&delimiter,"d","",s3Client.ADELIMITER)
	flag.IntVar(&limit,"m",100,s3Client.AMAXKEY)
	flag.BoolVar(&loop,"l",false,"-l loop over to N")
	flag.BoolVar(&trace,"t",false,s3Client.TRACEON)

	flag.Parse()

	if len(bucket) == 0  {
		flag.Usage()
		log.Fatalln(errors.New("*** bucket name cannot be empty"))
	}
	if len(outdir) == 0 {
		log.Fatalln(errors.New("*** output directory is missing"))
	} else {
		if _,err:= os.Stat(outdir); os.IsNotExist(err) {
			os.MkdirAll(outdir,s3Client.DIRMODE)
		}
	}

	/* get config  */
	s3Config,err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	/*
		init logging
	 */
	goLog.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	/*
		Create an S3 session
	 */
	s3Login := s3Core.New(s3Config,location)


	/*
		Build a List request  V1
	 */
	s3r := s3Core.S3ListRequest{
		MinioC: s3Login.MinioC,
		Bucket: bucket,
		Prefix: prefix,
		Delimiter: delimiter,
		Marker: marker,
		Limit: limit,
	}

	/*
		enable tracing  list http requests
	*/

	s3r.Trace = false
	if s3Client.TRACE || trace {
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
				goLog.Info.Printf("Bucket %s is empty",bucket)
				return
			}
			// download  objects returned by previous list Objects
			var nextMarker string
			for _,v := range results.Contents {
				//  download objetcs  concurrently using go routine
				go 	func(b string, k string ) {
					/* create a request */
					r := s3Core.S3GetRequest{}
					/* to be done  add request  options */
					opts:= &minio.GetObjectOptions{}
					/* build the request */
					r.S3BuildGetRequest(&s3Login,  bucket,  k,opts)

					//  Get object and meta data
					buf,objInfo,err := s3Core.GetObject(r)

					//  return a response
					ch <- Response{K:r.Key,Buffer:buf,ObjInfo:objInfo,Err:err}
				}(bucket,v.Key)
				nextMarker = v.Key
			}

			// wait for the completion of the  removal of n  objects

			for {
				select {
				case  r:= <-ch:
					{   t++
						if r.Err == nil {

							pathname := outdir + string(os.PathSeparator) + strings.Replace(r.K,string(os.PathSeparator),"_",-1)
							if err:= ioutil.WriteFile(pathname,r.Buffer.Bytes(),s3Client.FILEMODE); err == nil {
								goLog.Trace.Printf("Key %s is downloaded %s",r.K,pathname)
							} else {
								goLog.Error.Printf("Error %v downloading %s",err,r.K)
							}

							//  extract meta data
							m1 := s3Client.ExtractUserMeta(r.ObjInfo.Metadata)
							if len(m1) > 0 {
								if usermd, err := json.Marshal(m1); err == nil {
									goLog.Trace.Printf("Object Key %s - Size %d-%d  - Usermd %s - Content type: %s",r.K,r.Buffer.Len(),r.ObjInfo.Size,usermd,r.ObjInfo.ContentType)
								} else {
									goLog.Error.Printf("Error parsing user metadata of %s - %v",r.K,err)
								}
							} else {
								goLog.Warning.Printf("Object Key %s - Size %d-%d  - Usermd %s - Content type: %s",r.K,r.Buffer.Len(),r.ObjInfo.Size,"empty",r.ObjInfo.ContentType)
							}
						} else{
							goLog.Error.Printf("Error downloading object %s  of %s - %v",r.K,err)
						}
						if t == n {
							goLog.Info.Printf("%d objects  have been downloaded from  bucket:%s in %s\n", n, bucket,time.Since(start))
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
				Continue the next batch of object stats if loop and list returned was truncated
			 */
			if results.IsTruncated && loop {
				s3r.SetMarker(nextMarker)
			} else {
				goLog.Info.Printf("Download %d objects from bucket '%s' in %v \n",N,bucket,time.Since(start0))
				return
			}

		} else {
			log.Fatalf("Download  error %v\n",err)
		}
	}

}