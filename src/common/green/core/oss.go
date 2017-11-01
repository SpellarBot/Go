// @Author: Golion
// @Date: 2017.2

package core

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io/ioutil"
	"io"
)

///////////////////////////////////////////////////
// Aliyun OSS

const (
	AliyunEndpoint         = "oss-us-west-1.aliyuncs.com"
	AliyunAccessKeyId      = "UivNvn6XcKjb3FQc"
	AliyunAccessKeySecret  = "LVqK5NxgfMKbhfwfo6P9MjZZRosiYu"
	ImgAliyunBucketName    = "difoil"
)

func getAliyunConfig() (string, string, string) {
	endpoint        := AliyunEndpoint
	accessKeyId     := AliyunAccessKeyId
	accessKeySecret := AliyunAccessKeySecret
	return endpoint, accessKeyId, accessKeySecret
}

func ListOSSBucket() []string {
	ret := []string{}
	endpoint, accessKeyId, accessKeySecret := getAliyunConfig()
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Printf("Error When Init OSS Client - %v\n", err)
		return ret
	}
	lsRes, err := client.ListBuckets()
	if err != nil {
		fmt.Printf("Error When List Buckets - %v\n", err)
		return ret
	}
	cnt := 0
	for _, bucket := range lsRes.Buckets {
		fmt.Printf("List Buckets - %v\n", bucket.Name)
		ret = append(ret, fmt.Sprintf("%v", bucket.Name))
		cnt = cnt + 1
	}
	if cnt == 0 {
		fmt.Println("No Buckets")
	} else {
		fmt.Println("Finished List Buckets")
	}
	return ret
}

func CreateOSSBucket(bucketName string) bool {
	endpoint, accessKeyId, accessKeySecret := getAliyunConfig()
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Printf("Error When Init OSS Client - %v\n", err)
		return false
	}
	err = client.CreateBucket(bucketName)
	if err != nil {
		fmt.Printf("Error When Create Bucket - %v - %v\n", bucketName, err)
		return false
	} else {
		fmt.Printf("Create Bucket Succeed - %v\n", bucketName)
		return true
	}
}

func DeleteOSSBucket(bucketName string) bool {
	endpoint, accessKeyId, accessKeySecret := getAliyunConfig()
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Printf("Error When Init OSS Client - %v\n", err)
		return false
	}
	err = client.DeleteBucket(bucketName)
	if err != nil {
		fmt.Printf("Error When Delete Bucket - %v - %v\n", bucketName, err)
		return false
	} else {
		fmt.Printf("Delete Bucket Succeed - %v\n", bucketName)
		return true
	}
}

func PutOSSObject(bucketName string, objectKey string, res io.Reader) bool {
	endpoint, accessKeyId, accessKeySecret := getAliyunConfig()
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Printf("Error When Init OSS Client - %v\n", err)
		return false
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Printf("Error When Init Bucket - %v - %v\n", bucketName, err)
		return false
	}
	err = bucket.PutObject(objectKey, res)
	if err != nil {
		fmt.Printf("Error When Put Object - %v - %v - %v\n", bucketName, objectKey, err)
		return false
	} else {
		fmt.Printf("Put Object Succeed - %v - %v\n", bucketName, objectKey)
		return true
	}
}

func PutOSSObjectFromFile(bucketName string, objectKey string, filePath string) bool {
	endpoint, accessKeyId, accessKeySecret := getAliyunConfig()
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Printf("Error When Init OSS Client - %v\n", err)
		return false
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Printf("Error When Init Bucket - %v - %v\n", bucketName, err)
		return false
	}
	err = bucket.PutObjectFromFile(objectKey, filePath)
	if err != nil {
		fmt.Printf("Error When Put Object - %v - %v - %v - %v\n", bucketName, objectKey, filePath, err)
		return false
	} else {
		fmt.Printf("Put Object Succeed - %v - %v - %v\n", bucketName, objectKey, filePath)
		return true
	}
}

func DownloadOSSObject(bucketName string, objectKey string, filePath string) bool {
	endpoint, accessKeyId, accessKeySecret := getAliyunConfig()
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Printf("Error When Init OSS Client - %v\n", err)
		return false
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Printf("Error When Init Bucket - %v - %v\n", bucketName, err)
		return false
	}
	err = bucket.GetObjectToFile(objectKey, filePath)
	if err != nil {
		fmt.Printf("Error When Download Object - %v - %v - %v - %v\n", bucketName, objectKey, filePath, err)
		return false
	} else {
		fmt.Printf("Download Object Succeed - %v - %v - %v\n", bucketName, objectKey, filePath)
		return true
	}
}

func GetOSSObject(bucketName string, objectKey string) []byte {
	endpoint, accessKeyId, accessKeySecret := getAliyunConfig()
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Printf("Error When Init OSS Client - %v\n", err)
		return []byte{}
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Printf("Error When Init Bucket - %v - %v\n", bucketName, err)
		return []byte{}
	}
	body, err := bucket.GetObject(objectKey)
	if err != nil {
		fmt.Printf("Error When Get Object - %v - %v - %v\n", bucketName, objectKey, err)
		return []byte{}
	} else {
		fmt.Printf("Get Object Succeed - %v - %v\n", bucketName, objectKey)
	}
	data, err := ioutil.ReadAll(body)
	body.Close()
	if err != nil {
		fmt.Printf("Error When Read Input Stream - %v - %v - %v\n", bucketName, objectKey, err)
		return []byte{}
	}
	return data
}

func GetOSSObjectWithParams(bucketName string, objectKey string, urlParams string) []byte {
	endpoint, accessKeyId, accessKeySecret := getAliyunConfig()
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Printf("Error When Init OSS Client - %v\n", err)
		return []byte{}
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Printf("Error When Init Bucket - %v - %v\n", bucketName, err)
		return []byte{}
	}
	result, err := bucket.Client.Conn.Do("GET", bucketName, objectKey,
		urlParams, urlParams, nil, nil, 0, nil)
	if err != nil {
		fmt.Printf("Error When Get Object - %v - %v - %v\n", bucketName, objectKey, err)
		return []byte{}
	} else {
		fmt.Printf("Get Object Succeed - %v - %v\n", bucketName, objectKey)
	}
	body := result.Body
	data, err := ioutil.ReadAll(body)
	body.Close()
	if err != nil {
		fmt.Printf("Error When Read Input Stream - %v - %v - %v\n", bucketName, objectKey, err)
		return []byte{}
	}
	return data
}