package s3jsondir

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/alash3al/xyr/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// s3://s3-server-endpoint/bucket?region=region-name&ssl=true&path=true&perpage=1000
type Driver struct {
	s3      *s3.S3
	bucket  string
	perpage int64
}

// Open implements Importer#open
func (d *Driver) Open(dsn string) error {
	parsedDSN, err := url.Parse(dsn)
	if err != nil {
		return err
	}

	d.bucket = strings.Trim(parsedDSN.Path, "/")
	if strings.TrimSpace(d.bucket) == "" {
		return fmt.Errorf("you must specify the bucket as the url path")
	}

	d.perpage, _ = strconv.ParseInt(parsedDSN.Query().Get("perpage"), 10, 64)
	if d.perpage < 1 {
		d.perpage = 1000
	}

	accessKey, secretAccessKey := "", ""

	if parsedDSN.User != nil {
		accessKey = parsedDSN.User.Username()
		secretAccessKey, _ = parsedDSN.User.Password()
	}

	region := parsedDSN.Query().Get("region")
	if region == "" {
		region = "us-east-1"
	}

	s3config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretAccessKey, ""),
		Endpoint:         aws.String(parsedDSN.Hostname()),
		Region:           aws.String(region),
		DisableSSL:       aws.Bool(parsedDSN.Query().Get("ssl") == "false"),
		S3ForcePathStyle: aws.Bool(parsedDSN.Query().Get("path") == "true"),
	}

	newSession, err := session.NewSession(s3config)
	if err != nil {
		return err
	}

	d.s3 = s3.New(newSession)

	buckets, err := d.s3.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return err
	}

	bucketExists := false

	for _, bucket := range buckets.Buckets {
		if *bucket.Name == d.bucket {
			bucketExists = true
			break
		}
	}

	if !bucketExists {
		return fmt.Errorf("unknown bucket name specified (%s)", d.bucket)
	}

	return nil
}

// Import implements Importer#import
func (d *Driver) Import(s3prefix string) (<-chan map[string]interface{}, <-chan error, <-chan bool) {
	resultChan := make(chan map[string]interface{})
	errChan := make(chan error)
	doneChan := make(chan bool)

	s3ListObjectsInput := s3.ListObjectsV2Input{
		Bucket:  aws.String(d.bucket),
		Prefix:  aws.String(s3prefix),
		MaxKeys: aws.Int64(d.perpage),
	}

	walker := func(objectsList *s3.ListObjectsV2Output, isLastPage bool) bool {
		for _, item := range objectsList.Contents {
			if *item.Size < 1 {
				continue
			}

			output := &aws.WriteAtBuffer{}
			req := &s3.GetObjectInput{
				Bucket: aws.String(d.bucket),
				Key:    item.Key,
			}

			if _, err := s3manager.NewDownloaderWithClient(d.s3).Download(output, req); err != nil {
				errChan <- err
				continue
			}

			buf := bytes.NewBuffer(output.Bytes())
			decoder := json.NewDecoder(buf)

			for {
				var val interface{}

				if decoder.Decode(&val) == io.EOF {
					break
				}

				switch val := val.(type) {
				case map[string]interface{}:
					resultChan <- val
				case []interface{}:
					mSlice, err := utils.InterfaceSliceToMapStringInterfaceSlice(val)
					if err != nil {
						errChan <- err
						continue
					} else {
						for _, item := range mSlice {
							resultChan <- item
						}
					}
				default:
					errChan <- fmt.Errorf("unsupported value (%v), we only support array of objects or just objects", val)
				}
			}
		}

		return true
	}

	go (func() {
		defer (func() {
			doneChan <- true

			close(resultChan)
			close(errChan)
			close(doneChan)
		})()

		if err := d.s3.ListObjectsV2Pages(&s3ListObjectsInput, walker); err != nil {
			errChan <- err
			return
		}
	})()

	return resultChan, errChan, doneChan
}
