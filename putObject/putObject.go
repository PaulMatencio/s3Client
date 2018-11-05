
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
	bucketName 	string
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
		objectName string
		separator string
	)

	flag.StringVar(&bucketName,"b","","-b bucketName")
	flag.StringVar(&location,"s","site1","-s locationName")
	flag.StringVar(&filename,"fn","","-fn File name")
	flag.StringVar(&objectName,"o","","-o object name")
	flag.StringVar(&separator,"separator","/","-sep  <aSeparator>")
	flag.BoolVar(&trace,"trace",false,"-trace")
	flag.Parse()
	if len(bucketName) == 0  ||  len(filename) == 0 {
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
	s3Login := s3Client.LoginS3(s3Config,location)
	minioc := s3Login.GetS3Client() // get minio Client
	/* read the file */
	object, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}

	defer object.Close()
	objectStat,err := object.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	if len(objectName) == 0 {
		sl := strings.Split(filename, separator)
		objectName = sl[len(sl)-1]
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
	n,err := minioc.PutObject(bucketName, objectName,object, objectStat.Size(),opts)
	if err != nil {
		log.Fatalln(err)
	}

	printOk(objectName,n)
}