package main

import (
	"encoding/json"
	"flag"
	"github.com/s3Client/lib"
	"log"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go"
)

func main() {
	var (
		bucket		string
		location 	string
		filename 	string
		signature   string
		object 	string
		delimiter 	string
		trace		bool
		s3Login     s3Client.S3Login
	)

	flag.StringVar(&bucket,"b","","-b bucketName")
	flag.StringVar(&location,"s","site1","-s locationName")
	flag.StringVar(&object,"o","","-o objectName")
	flag.StringVar(&signature,"S","V4","-S signature")
	flag.BoolVar(&trace,"trace",false,"-trace")
	flag.Parse()
	if len(bucket) == 0  || len(object) == 0 {
		log.Println("bucket name  or object name cannot be empty")
		flag.Usage()
		os.Exit(100)
	}
	s3Client.TRACE = trace
	if filename == "" {
		sl := strings.Split(object,delimiter)
		filename = sl[len(sl)-1]
	}

	/* get Config */
	s3Config,err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	/* login to s3  NewV2 signature V2*/
	/* login to s3  New signature V4*/

	if signature == "V2" {
		s3Login = s3Client.NewV2(s3Config, location)
	} else {
		s3Login = s3Client.New(s3Config, location)
	}
	minioc   := s3Login.GetS3Client() // get minio Client

	if trace {
		minioc.TraceOn(os.Stdout)
	}
	/* Get Stat with Options */
	opts:= minio.StatObjectOptions{}
	start:= time.Now()
	objInfo, err := minioc.StatObject(bucket, object,opts)
	if err != nil {
		log.Fatalln(err)
	}
	/* Extract user meta of the http.Header and convert it to json */
	m1 := s3Client.ExtractUserMeta(objInfo.Metadata)
	metadata, _ := json.Marshal(m1)

	log.Printf("Duration: %s\nObject key : %s\nContent Type: %s\nLast Modified: %s\nOwner %s\nSize: %d\nUser metadata: %s\n",time.Since(start),objInfo.Key, objInfo.ContentType,  objInfo.LastModified, objInfo.Owner, objInfo.Size,metadata,)


}
