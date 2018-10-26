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
		bucketName string
		location string
		filename string
		objectName string
		delimiter string
	)

	flag.StringVar(&bucketName,"b","","-b bucketName")
	flag.StringVar(&location,"s","site1","-s locationName")
	flag.StringVar(&objectName,"o","","-o objectName")
	flag.Parse()
	if len(bucketName) == 0  || len(objectName) == 0 {
		fmt.Println("bucketName or objectName cannot be empty")
		flag.Usage()
		os.Exit(100)
	}
	if filename == "" {
		sl := strings.Split(objectName,delimiter)
		filename = sl[len(sl)-1]
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
	/* Get Stat with Options */
	opts:= minio.StatObjectOptions{}
	start:= time.Now()
	objInfo, err := s3client.StatObject(bucketName, objectName,opts)
	if err != nil {
		log.Fatalln(err)
	}
	/* Extract user meta of the http.Header and convert it to json */
	m1 := s3Client.ExtractUserMeta(objInfo.Metadata)
	metadata, _ := json.Marshal(m1)

	fmt.Printf("Duration: %s\nObject key : %s\nContent Type: %s\nLast Modified: %s\nOwner %s\nSize: %d\nUser metadata: %s\n",time.Since(start),objInfo.Key, objInfo.ContentType,  objInfo.LastModified, objInfo.Owner, objInfo.Size,metadata,)


}
