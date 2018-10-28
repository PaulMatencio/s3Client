
package main

import (
	"errors"
	"flag"
	"github.com/minio/minio-go"
	"github.com/s3Client/lib"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// "time"
)

func main() {
	var (
		bucketName string /* Bucket name  */
		location   string /* S3 location */
		filename   string /* output file */
		objectName string /* Object name */

	)

	flag.StringVar(&bucketName, "b", "", "-b bucketName")
	flag.StringVar(&location, "s", "site1", "-s locationName")
	flag.StringVar(&objectName, "o", "", "-o objectName")
	flag.StringVar(&filename, "fn", "", "-fn filename")

	flag.Parse()
	if len(bucketName) == 0 || len(filename) == 0 {
		flag.Usage()
		log.Fatalln(errors.New("bucketName or filename cannot be empty"))
	}

	if len(objectName) == 0 {
		sl := strings.Split(filename, "/")
		objectName = sl[len(sl)-1]
	}


	/* get Config */
	s3Config, err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	/* create an S3 session */
	s3 := s3Client.SetS3Session(s3Config, location)
	s3client, err := minio.New(s3.Endpoint, s3.AccessKeyID, s3.SecretKey, s3.SSL)
	if err != nil {
		log.Fatalln(err)
	}

	/* set transport option */
	tr := &http.Transport{
		DisableCompression: true,
	}

	s3client.SetCustomTransport(tr)
	s3client.TraceOn(os.Stdout)
	opts := minio.PutObjectOptions{}
	opts.ContentType="application/octet-stream"
	start := time.Now()
	if  _,err := s3client.FPutObject(bucketName, objectName, filename,opts); err != nil {
		log.Fatalln(err)
	}
	log.Printf("Duration:  %s", time.Since(start))
}