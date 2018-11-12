package s3Core

import (
	"github.com/minio/minio-go"
	"github.com/s3Client/lib"
	"os"
)

type S3ListRequest struct {
	MinioC 		*minio.Core
	Bucket 		string
	Prefix  	string
	Delimiter   string
	StartAfter  string
	Marker		string
	Limit       int
	FetchOwner	bool
	Trace 		bool
	// Results     *[]minio.ObjectInfo
}

func (r *S3ListRequest) SetS3Core(s3c *minio.Core) {
	r.MinioC = s3c
}

func (r *S3ListRequest) SetBucketName(bucket string) {
	r.Bucket = bucket
}

func (r *S3ListRequest) SetPrefix(prefix string) {
	r.Prefix = prefix
}

func (r *S3ListRequest) SetStartAfter(startAfter string) {
	r.StartAfter = startAfter
}

func (r *S3ListRequest) SetMarker(marker string) {
	r.Marker = marker
}

func (r *S3ListRequest) SetTrace(trace bool) {
	r.Trace = trace
}

func (r *S3ListRequest) S3BuildListRequest(login *S3Login, bucket string, prefix string, fetchOwner bool,delimiter string, startAfter string, Marker string, limit int){

	r.MinioC 		=  login.MinioC
	r.Bucket 		=  bucket
	r.Prefix		=  prefix
	r.StartAfter 	=  startAfter
	r.Marker		=  Marker
	r.Limit         = limit
	r.Trace = false
	if s3Client.TRACE {
		r.Trace	= true
	}
	r.FetchOwner = fetchOwner


}

/* https://godoc.org/github.com/minio/minio-go#ListBucketResult */


func ListObjectsV2(r S3ListRequest) (minio.ListBucketV2Result, error) {
	if r.Trace {
		r.MinioC.TraceOn(os.Stdout)
	}
	return r.MinioC.ListObjectsV2(r.Bucket, r.Prefix, r.Marker, r.FetchOwner, r.Delimiter, r.Limit,r.StartAfter )
}