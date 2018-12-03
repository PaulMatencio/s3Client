package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/minio/minio-go"
	"github.com/moses/user/goLog"
	"github.com/s3Client/lib"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	site1 						s3Client.Host
	ssl,trace,overwrite 		bool
	start      					time.Time
	bucket,endpoint,location, prefix,outdir 	string
)

type Response struct {
	Key 	string
	FileName string
	Err      error
}

func main() {

	flag.StringVar(&bucket, "b", "", s3Client.ABUCKET)
	flag.StringVar(&location, "s", "site1", s3Client.ALOCATION)
	flag.StringVar(&prefix, "p", "", s3Client.APREFIX)
	flag.StringVar(&outdir,"o","",s3Client.AOUTPUTDIR)
	flag.BoolVar(&trace, "t", false, s3Client.TRACEON)
	flag.BoolVar(&overwrite,"O",false,"-O to overwrite output files")
	flag.Parse()

	// check input parameters
	if len(bucket) == 0 || len(outdir) == 0 {
		flag.Usage()
		if len(bucket) == 0 {
			log.Fatalln(errors.New("*** bucket is missing"))
		}
		if len(outdir) == 0 {
			log.Fatalln(errors.New("*** output directory is missing"))
		}
	}
	// create the oouput directory if it does not exist
	if _,err := os.Stat(outdir); err != nil {
		os.MkdirAll(outdir,0744)
	}

	/*
		init logging
	 */
	goLog.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	/* get the Config */
	s3Config, err := s3Client.GetConfig("config.json")
	if err != nil {
		log.Fatalln(err)
	}

	/* login to S3  */
	s3Login := s3Client.New(s3Config,location)
	minioc := s3Login.GetS3Client()  	// get minio client
	runtime.GOMAXPROCS(4)

	// list the Objects of a  bucket
	// Create a done channel to control 'ListObjects' go routine.
	doneCh := make(chan struct{})
	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	// List the buckets with a  prefix

	keys := []string{}
	for objInfo := range minioc.ListObjects(bucket, prefix, true, doneCh) {
		if objInfo.Err != nil {
			fmt.Println(objInfo.Err)
			return
		}
		keys = append(keys,objInfo.Key)
	}

	start0 := time.Now()
	N,T := len(keys),1
	if (N == 0) {
		log.Printf("Bucket %s is empty",bucket)
		return
	}
	ch := make(chan Response)
	runtime.GOMAXPROCS(4)

	if trace {
		minioc.TraceOn(os.Stdout)
	}

	for obj := 0; obj < N; obj++ {
		start = time.Now()
		key := keys[obj]
		filePath := filepath.Join(outdir,key)
		go func( string,string) {
			r := s3Client.S3FGetRequest{}
			// todo add getObject Options
			options := &minio.GetObjectOptions{}

			r.S3BuildFGetRequest(&s3Login,  bucket, key, filePath, options, overwrite)
			err := r.FGetObject()
			ch <- Response{key,filePath,err}
		}(key,filePath)
	}

	/*  wait  until all objects are downloaded  */
	for {

		select {
		case    r:= <-ch:
			{
				log.Printf("Download Object key: %s to Filename: %s - Error: %v",r.Key,r.FileName,r.Err)
				if T == N {
					log.Printf("Downloaded %d objects  from %s in %s\n", N, bucket,time.Since(start0))
					return
				}
				T++
			}
		case <-time.After(50 * time.Millisecond):
			fmt.Printf("w")
		}
	}

}

