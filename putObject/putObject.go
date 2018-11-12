
package main

import (
	"errors"
	"flag"
	"github.com/minio/minio-go"
	"github.com/s3Client/lib"
	"log"
	"os"
	"strings"
	"time"
)

var (
	bucket	 	string
	location 	string
	endpoint 	string
	site1 		s3Client.Host
	ssl 		bool
	start       time.Time
	trace       bool
)

func printOk(object string, size int64) {
	duration := float64(time.Since(start) / 1000000.0)
	speed :=  float64(size /1024.0) / duration
	log.Printf("uploaded  %s of size %d successfully . Duration %.3f msec. Speed  %.3f  MBps",object,size,duration,speed)
}


func main() {

	var (
		filename string
		object string
		separator string
	)

	flag.StringVar(&bucket,"b","",s3Client.ABUCKET)
	flag.StringVar(&location,"s","site1",s3Client.ALOCATION)
	flag.StringVar(&filename,"f","",s3Client.AFILE)
	flag.StringVar(&object,"o","",s3Client.AOBJECT)
	flag.StringVar(&separator,"d",s3Client.DELIMITER,s3Client.ADELIMITER)
	flag.BoolVar(&trace,"t",false,s3Client.TRACEON)
	flag.Parse()
	if len(bucket) == 0  ||  len(filename) == 0 {
		flag.Usage()
		log.Fatalln(errors.New("bucketName or filename cannot be empty"))
	}

	if trace {
		s3Client.TRACE = true
	}

	/* get S3 Config */
	s3Config,err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}
	/* login to s3 */
	s3Login := s3Client.New(s3Config,location)
	minioc := s3Login.GetS3Client() // get minio Client
	/* read the file */
	fi, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}

	defer fi.Close()
	objectStat,err := fi.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	if len(object) == 0 {
		sl := strings.Split(filename, separator)
		object = sl[len(sl)-1]
	}

	opts:= minio.PutObjectOptions{}
	opts.ContentType = "application/octet-stream"
	opts.StorageClass= "STANDARD"

	usermd := map[string]string {
		"lastName": "Matencio",
		"firstname": "Paul",
		"address": "Regentesselaan 14",
	}
	opts.UserMetadata = usermd



	if s3Client.TRACE {
		minioc.TraceOn(os.Stdout)
	}
	start = time.Now()
	n,err := minioc.PutObject(bucket, object,fi, objectStat.Size(),opts)
	if err != nil {
		log.Fatalln(err)
	}

	printOk(object,n)
}