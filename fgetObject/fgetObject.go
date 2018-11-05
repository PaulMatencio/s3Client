
package main

import (
	"flag"
	"runtime"
	"strings"
	"time"

	"errors"
	"github.com/minio/minio-go"
	"github.com/s3Client/lib"
	"log"
)

func main() {
	var (
		bucketName 	string /* Bucket name  */
		location   	string /* S3 location */
		filename   	string /* output file */
		objectName 	string /* Object name */
		delimiter  	string
		trace		bool
	)

	flag.StringVar(&bucketName, "b", "", "-b bucketName")
	flag.StringVar(&location, "s", "site1", "-s locationName")
	flag.StringVar(&objectName, "o", "", "-o objectName")
	flag.StringVar(&filename, "fn", "", "-fn filename")
	flag.StringVar(&delimiter, "delimiter", "/", "-delimiter /")
	flag.BoolVar(&trace,"trace",false,"-trace ")
	flag.Parse()

	if len(bucketName) == 0 || len(objectName) == 0 {
		flag.Usage()
		log.Fatalln(errors.New("bucketName or objectName cannot be empty"))
	}

	s3Client.TRACE = trace

	/* parse the path of the filename and Keep only the last path to form the object name */
	if len(filename) == 0 {
		sl := strings.Split(objectName, delimiter)
		filename = sl[len(sl)-1]
	}

	/* get Config */
	s3Config, err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	s3Login := s3Client.LoginS3(s3Config,location)
	minioc := s3Login.GetS3Client()  				// get minio s3Clie

	runtime.GOMAXPROCS(4)
	minioc.SetCustomTransport(s3Client.TR)

	start := time.Now()
	if err := minioc.FGetObject(bucketName, objectName, filename, minio.GetObjectOptions{}); err != nil {
		log.Fatalln(err)
	}
	log.Printf("Duration:  %s", time.Since(start))
}