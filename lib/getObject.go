package s3Client

import (
	"bytes"
	"errors"
	"github.com/minio/minio-go"
	"os"
)


type S3Request struct {
	MinioC 	*minio.Client
	// Buf 	*bytes.Buffer
	Bucket 	string
	Key 	string
	Opts 	*minio.GetObjectOptions
	BufSize int
	Trace 	bool
}



func (r *S3Request) SetS3Client(s3c *minio.Client) {
	r.MinioC = s3c
}


func (r *S3Request) SetBucketName(bucket string) {
	r.Bucket = bucket
}

func (r *S3Request) SetKey(key string) {
	r.Bucket = key
}

func (r *S3Request) SetPutOpts(options *minio.GetObjectOptions) {
	r.Opts= options
}

func (r *S3Request) SetTrace(trace bool) {
	r.Trace = trace
}

func (r *S3Request) S3BuildGetRequest(login *S3Login, bucket string, key string, options *minio.GetObjectOptions){
	r.MinioC,r.Bucket,r.Key,r.Opts,r.BufSize,r.Trace	=  login.MinioC,bucket,key,options,BUFFERSIZE,false
	if TRACE {
		r.Trace	= true
	}
}


/*

 */

func GetObject(r S3Request) (*bytes.Buffer, error) {

	var (
		object *minio.Object
		err error
		n, size int
	)
	if r.Trace {
		r.MinioC.TraceOn(os.Stdout)
	}

	// minio,getObject returns the address of a object stream
	/* Abvailable methods are
	  object.Read( [] byte)  ->  return the number of bytes read
	  object.Stat() ->  return minio.ObjectInfo
	  object.ReadAt([]byte,offset int64)   -> return the number of bytes
	  object.Seek(offset int64, whence int)
	  object.Close()   Close the object stream
	*/


	// it is possible to check if the Object exist by using object.Stat(). a HTTP HEAD request will be sent to S3
	// However it is faster no to do it when the S3 server is distant
	// Just  getObject ( HTTP GET)

	if object,err = r.MinioC.GetObject(r.Bucket,r.Key,*r.Opts); err != nil {
		return nil,err
	}

	// the stream will be closed when the function is exited

	defer object.Close()

	/* Reading  the object*/
	buffer:= make([]byte,r.BufSize)
	buf := new(bytes.Buffer)

	for {
		n, err = object.Read(buffer)
		if err != nil {
			size += n
			buf.Write(buffer[0:n])
			if n == 0 {
				break
			}
		} else {
			return buf,err
		}
	}

	if buf.Len() == 0 {
		err = errors.New(r.Key + " is empty or not found")
	}
	return buf,err

}
