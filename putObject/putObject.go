
package main

import (
	"flag"
	"fmt"
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
	)

	flag.StringVar(&bucketName,"b","","-b bucketName")
	flag.StringVar(&location,"s","site1","-s locationName")
	flag.StringVar(&filename,"o","","-on objectName")
	flag.Parse()
	if len(bucketName) == 0  ||  len(filename) == 0 {
		fmt.Println("bucketName or objectName cannot be empty")
		flag.Usage()
		os.Exit(100)
	}

	/* get Config */

	s3Config,err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	/* create an S3 session */
	s3:= s3Client.SetS3Session(s3Config,location)
	s3client, err := minio.New(s3.Endpoint, s3.AccessKeyID, s3.SecretKey,s3.SSL)
	if err != nil {
		log.Fatalln(err)
	}
	/* */
	object, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}

	defer object.Close()
	objectStat, err := object.Stat()

	if err != nil {
		log.Fatalln(err)
	}
	sl := strings.Split(filename,"/")
	objectName =  sl[len(sl)-1]


	opts:= minio.PutObjectOptions{}

	opts.ContentType = "application/octet-stream"
	opts.StorageClass= "STANDARD"

	usermd := make(map[string]string)
	usermd["lastName"] = "Matencio"
	usermd["fisrtName"]= "Paul"
	usermd["address"] = "Regentesselaan 14"
	opts.UserMetadata = usermd

	/* */
	s3client.TraceOn(os.Stdout)
	start = time.Now()
	n,err := s3client.PutObject(bucketName, objectName,object,
		objectStat.Size(),opts)

	if err != nil {
		log.Fatalln(err)
	}

	printOk(objectName,n)
}