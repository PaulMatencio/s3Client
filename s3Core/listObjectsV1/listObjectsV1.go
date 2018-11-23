
package main

import (
	"errors"
	"flag"
	"github.com/s3Client/lib"
	"github.com/s3Client/s3Core/lib"
	"log"
	"time"
)


func main() {

	var (
		bucket 		string
		location 	string
		prefix		string
		delimiter   string
		marker      string
		limit 		int
		trace		bool
		loop        bool
	)

	/* define  input parameters */
	flag.StringVar(&bucket,"b","",s3Client.ABUCKET)
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
		log.Fatalln(errors.New("bucket name cannot be empty"))
	}




	/* get config  */
	s3Config,err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}
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
		disable trace
		enable trace  list http requests
	*/

	s3r.Trace = false
	if s3Client.TRACE || trace {
		s3r.Trace	= true
	}

	// s3r.S3BuildListRequest(&s3Login, bucket, prefix, false,delimiter, after, next,limit)

	start := time.Now()
	for {

		if results,err := s3Core.ListObjectsV1(s3r) ; err == nil {
			var  nextMarker string

			for k,v := range results.Contents {
				log.Println(k,v.Key,v.Size)
				nextMarker = v.Key  // Ugly bud needed because Scality S3 does not return  next marker as expected : Bug
			}

			log.Printf("Is truncated  ? %v - Marker: %s - Next marker: %s", results.IsTruncated, results.Marker, nextMarker)
			if results.IsTruncated && loop {
			// 	s3r.SetMarker(results.NextMarker)
				s3r.SetMarker(nextMarker)
			} else {
				log.Printf("Elapsed time: %s",time.Since(start))
				return
			}

		} else {
			log.Fatalf("List error %v",err)
		}
	}


}