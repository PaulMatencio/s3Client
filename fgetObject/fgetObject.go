
package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"errors"
	"github.com/minio/minio-go"
	"github.com/s3Client/lib"
	"log"
	"net/http"
)

func main() {
	var (
		bucketName string /* Bucket name  */
		location   string /* S3 location */
		filename   string /* output file */
		objectName string /* Object name */
		delimiter  string
	)

	flag.StringVar(&bucketName, "b", "", "-b bucketName")
	flag.StringVar(&location, "s", "site1", "-s locationName")
	flag.StringVar(&objectName, "o", "", "-o objectName")
	flag.StringVar(&filename, "fn", "", "-fn fileName")

	flag.Parse()
	if len(bucketName) == 0 || len(objectName) == 0 {
		flag.Usage()
		log.Fatalln(errors.New("bucketName or objectName cannot be empty"))
	}

	/* parse the path of the filename and Keep only the last path to form the object name */
	if filename == "" {
		sl := strings.Split(objectName, delimiter)
		filename = sl[len(sl)-1]
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
	start := time.Now()
	if err := s3client.FGetObject(bucketName, objectName, filename, minio.GetObjectOptions{}); err != nil {
		log.Fatalln(err)
	}
	log.Printf("Duration:  %s", time.Since(start))
}