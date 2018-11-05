
package main

import (
	"errors"
	"flag"
	"github.com/minio/minio-go"
	"github.com/s3Client/lib"
	"log"
	"time"
)

func main() {
	var (
		bucketName 		string              /* Bucket name  */
		location 		string              /* S3 location */
		objectName 		string              /* Object name */
		trace			bool
		bufsize			int
	)

	flag.StringVar(&bucketName,"b","","-b bucketName")
	flag.StringVar(&location,"s","site1","-s locationName")
	flag.StringVar(&objectName,"o","","-o objectName")
	flag.BoolVar(&trace,"trace",false,"-trace ")
	flag.IntVar(&bufsize,"size",65536,"-size of the buffer")
	// flag.StringVar(&filename,"fn","","-fn fileName")

	flag.Parse()
	if len(bucketName) == 0  || len(objectName) == 0 {
		flag.Usage()
		log.Fatalln(errors.New("bucketName or objectName cannot be empty"))
	}

	/* parse the path of the filename and Keep only the last path to form the object name
	if filename == "" {
		sl := strings.Split(objectName,delimiter)
		filename = sl[len(sl)-1]
	}
	*/


	s3Client.TRACE = trace

	s3Config,err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	s3Login := s3Client.LoginS3(s3Config,location)
	minioc := s3Login.GetS3Client()  // get minio s3Clie

	/* set transport option
	tr := &http.Transport{
		DisableCompression: true,
	}

	s3client.SetCustomTransport(tr)
	*/

	start := time.Now()

	// object reader. It implements io.Reader, io.Seeker, io.ReaderAt and io.Closer interfaces.
	//minio.GetObjectOptions.SetRange(0,2000)
	object, err := minioc.GetObject(bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	defer object.Close()

	/* Retrieve the object*/
	buffer:= make([]byte,bufsize )
	var (
		n int
		size int
	)
	for {
		n, _ = object.Read(buffer)
		size += n
		if n == 0 {
			break
		}
	}
	log.Printf("Total Duration:  %s  size: %d",time.Since(start),size)
}
