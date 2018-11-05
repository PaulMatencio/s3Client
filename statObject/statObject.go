package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/s3Client/lib"
	"log"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go"
)

func main() {
	var (
		bucketName 	string
		location 	string
		filename 	string
		objectName 	string
		delimiter 	string
		trace		bool
	)

	flag.StringVar(&bucketName,"b","","-b bucketName")
	flag.StringVar(&location,"s","site1","-s locationName")
	flag.StringVar(&objectName,"o","","-o objectName")
	flag.BoolVar(&trace,"trace",false,"-trace")
	flag.Parse()
	if len(bucketName) == 0  || len(objectName) == 0 {
		fmt.Println("bucketName or objectName cannot be empty")
		flag.Usage()
		os.Exit(100)
	}
	s3Client.TRACE = trace
	if filename == "" {
		sl := strings.Split(objectName,delimiter)
		filename = sl[len(sl)-1]
	}

	/* get Config */
	s3Config,err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	/* login to s3 */
	s3Login := s3Client.LoginS3(s3Config,location)
	minioc   := s3Login.GetS3Client() // get minio Client

	/* Get Stat with Options */
	opts:= minio.StatObjectOptions{}
	start:= time.Now()
	objInfo, err := minioc.StatObject(bucketName, objectName,opts)
	if err != nil {
		log.Fatalln(err)
	}
	/* Extract user meta of the http.Header and convert it to json */
	m1 := s3Client.ExtractUserMeta(objInfo.Metadata)
	metadata, _ := json.Marshal(m1)

	fmt.Printf("Duration: %s\nObject key : %s\nContent Type: %s\nLast Modified: %s\nOwner %s\nSize: %d\nUser metadata: %s\n",time.Since(start),objInfo.Key, objInfo.ContentType,  objInfo.LastModified, objInfo.Owner, objInfo.Size,metadata,)


}
