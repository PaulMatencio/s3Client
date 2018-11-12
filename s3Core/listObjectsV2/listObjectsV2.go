package main

import (
	"errors"
	"flag"
	"github.com/s3Client/lib"
	"github.com/s3Client/s3Core/lib"
	"log"
)


func main() {

	var (
		bucket 		string
		location 	string
		prefix		string
		// next        string
		delimiter   string
		after       string
		limit 		int
		fetchOwner  bool
		trace		bool
		loop        bool
	)

	/* define  input parameters */
	flag.StringVar(&bucket,"b","",s3Client.ABUCKET)
	flag.StringVar(&location,"s","site1",s3Client.ALOCATION)
	flag.StringVar(&prefix,"p","",s3Client.APREFIX)
	//flag.StringVar(&next,"next","","-next <Next continuation Token>")
	flag.StringVar(&after,"after","","-after <aString>")
	flag.StringVar(&delimiter,"d","",s3Client.ADELIMITER)
	flag.IntVar(&limit,"m",100,s3Client.AMAXKEY)
	flag.BoolVar(&fetchOwner,"fetchOwner",false,"-fo fetchOwner")
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
		Build a List request
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


	for {

		if results,err := s3Core.ListObjectsV2(s3r) ; err == nil {

			for k,v := range results.Contents {
				log.Println(k,v.Key,v.Size)
			}

			log.Printf("Is truncated  ? %v - After: %s - Next: %s", results.IsTruncated, results.StartAfter,results.NextContinuationToken)
			if results.IsTruncated && loop {
				s3r.SetStartAfter(results.NextContinuationToken)
			} else {
				return
			}

		} else {
			log.Fatalf("List error %v",err)
		}
	}

}


