package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/minio/minio-go"
	"github.com/s3Client/lib"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"
)

var (
	bucketName 	string
	location 	string
	endpoint 	string
	site1 		s3Client.Host
	ssl 		bool
	start       time.Time
	prefix      string
)


func main() {
	var (
		directory   string
	)
	type Response struct {
		Filename  string
		Size      int64
		Duration  time.Duration
		Err       error
	}

	flag.StringVar(&bucketName, "b", "", "-b bucketName")
	flag.StringVar(&location, "s", "site1", "-s locationName")
	flag.StringVar(&directory, "d", "", "-d directory")
	flag.StringVar(&prefix, "prefix", "images/jpeg/", "-prefix <prefix>")
	flag.Parse()
	if len(bucketName) == 0 || len(directory) == 0 {
		log.Fatalln(errors.New("bucketName or objectName cannot be empty"))
	}

	/* read the directory */
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	filenames := []string{}
	for _, f := range files {
		if !f.IsDir() {
			filenames= append(filenames, f.Name())
		}
	}

	/* get the Config */

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
	opts := minio.PutObjectOptions{}
	opts.ContentType = "application/octet-stream"
	opts.StorageClass = "STANDARD"
	/* add some user metadata  */
	s3Client.SetMetadata(opts)
	/* create a goroutine chanel */
	messages := make(chan Response)
	runtime.GOMAXPROCS(4)
	N:= len(filenames)
	T:=1
	/* */
	// s3client.TraceOn(os.Stdout)
	start0 := time.Now()
	for  obj:=0; obj < N; obj++ {
		start = time.Now()
		var n int64 = 0;
		filename := filenames[obj]
		file,err  := os.Open(filename)
		defer file.Close()
		if err == nil {
			go func(*os.File,string) {
				var (
				 size int64
				 start time.Time
				)
				if err == nil {
					fileStat, _ := file.Stat()
					size   = fileStat.Size()
					start  = time.Now()
					objectName := prefix + filename
					n, err = s3client.PutObject(bucketName, objectName, file,
						size, opts)
					if err != nil {
						s3Client.PrintNotOk(filename, err)
					}
				} else {
					s3Client.PrintNotOk(filename, err)
				}
				s3Client.PrintOk(filename,n,start)  /* ok or not */
				messages <- Response{filename,size, time.Since(start), err}
			}(file,filename)
		}
	}
	var (
		totalSize int64 = 0
		totalDuration time.Duration = 0
	)
	for {
		select {
		case   r:=  <-messages:
			{
				totalSize += r.Size
				totalDuration += r.Duration
				if T ==  N  {
					elapsedTm := time.Since(start0)
					MBps :=  1000* float64(totalSize)/float64(elapsedTm)
					fmt.Printf("Uplooad %d objects[ Total size (byte): %d - MBps: %0.4f] in Total Duration:%s/Total Elapsed:%s\n",N, totalSize, MBps, totalDuration, elapsedTm)
					return
			     }
				T++
			}
		case <-time.After(50 * time.Millisecond):
			fmt.Printf("w")
		}
	}


}