
package main

import (
	"flag"
	"strings"
	"time"

	"errors"
	"github.com/minio/minio-go"
	"github.com/s3Client/lib"
	"log"
)

func main() {
	var (
		bucket 	string /* Bucket name  */
		location   	string /* S3 location */
		filename   	string /* output file */
		object 	string /* Object name */
		delimiter  	string
		trace		bool
	)

	flag.StringVar(&bucket, "b", "", s3Client.ABUCKET)
	flag.StringVar(&location, "s", "site1", s3Client.ALOCATION)
	flag.StringVar(&object, "o", "", s3Client.AOBJECT)
	flag.StringVar(&filename, "f", "", s3Client.AFILE)
	flag.StringVar(&delimiter, "d", s3Client.DELIMITER, s3Client.ADELIMITER)
	flag.BoolVar(&trace,"t",false,s3Client.TRACEON)
	flag.Parse()

	if len(bucket) == 0 || len(object) == 0 {
		flag.Usage()
		log.Fatalln(errors.New("bucketName or objectName cannot be empty"))
	}

	s3Client.TRACE = trace

	/* parse the path of the filename and Keep only the last path to form the object name */
	if len(filename) == 0 {
		sl := strings.Split(object, delimiter)
		filename = sl[len(sl)-1]
	}

	/* get Config */
	s3Config, err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	s3Login := s3Client.New(s3Config,location)

	minioc := s3Login.GetS3Client()  				// get minio client

	/*
	runtime.GOMAXPROCS(4)
	minioc.SetCustomTransport(s3Client.TR)
	*/

	start := time.Now()
	if err := minioc.FGetObject(bucket, object, filename, minio.GetObjectOptions{}); err != nil {
		log.Fatalln(err)
	}
	log.Printf("Duration:  %s", time.Since(start))
}