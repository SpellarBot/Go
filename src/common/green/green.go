//By botle.lwt
//2017.6.29
package green

import (
	"fmt"
	"vidmate.com/common/green/core"
	"net/http"
	"encoding/json"
	"time"
)

func GetPorn(imageUrl string, imageId string) (error, core.PornScannerResult) {
	var result core.PornScannerResult
	res, err := http.Get(imageUrl)
	if res == nil {
		fmt.Printf("Http Error: return null- %v - %v\n", imageUrl, err)
		result.Status = 1
		result.Msg    = err.Error()
		return nil, result
	}
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("Http Error - %v - %v\n", imageUrl, err)
		result.Status = 1
		result.Msg    = err.Error()
		return err, result
	}
	bucketName := core.ImgAliyunBucketName
	core.PutOSSObject(bucketName, imageId, res.Body)
	pornResult, err := core.GetPornScore(bucketName, imageId)
	if err != nil {
		result.Status = 1
		result.Msg    = err.Error()
		result.Result = pornResult
		return err, result
	} else {
		result.Status = 0
		result.Msg    = "Succeed"
		result.Result = pornResult
		return nil, result
	}
}

func PutVideoPorn(videoUrl string, videoId string) (error, bool) {
	bucketName := core.ImgAliyunBucketName

	res, err := http.Get(videoUrl)
	if res == nil {
		return fmt.Errorf("Get Video URL Return Null, msg show as %s", err.Error()), false
	}
	defer res.Body.Close()
	if err != nil {
		return fmt.Errorf("Http Error - %v - %v\n", videoUrl, err.Error()), false
	}

	isSucceed := core.PutOSSObject(bucketName, videoId, res.Body)

	if isSucceed {
		return nil, true
	} else {
		return fmt.Errorf("Put Video OSS Object Fail, Msg Show As %s", err.Error()), false
	}
}

func PutVideoPornFromFile(videoId string, filePath string) (error, bool) {
	bucketName := core.ImgAliyunBucketName

	isSucceed := core.PutOSSObjectFromFile(bucketName, videoId, filePath)

	if isSucceed {
		return nil, true
	} else {
		return fmt.Errorf("Put Video Oss Object From File Fail"), false
	}
}

func SubmitVideoCheck(bucketName string, videoId string) string {
	//提交视频检测任务
	ret := core.GetOSSObjectWithParams(bucketName, videoId, "x-oss-process=udf/green/video/scan,porn")
	//获取视频提交任务返回结果
	var result core.VideoSubmitRet
	json.Unmarshal(ret, &result)
	fmt.Printf("Put Video Check Task Result - %v\n", result.ToJSONStr())

	if result.Code != 200 {
		fmt.Printf("Put Video Check Task %s Fail, Msg Show As %s\n", videoId, result.Msg)
		return ""
	} else {
		return result.Data.TaskId
	}
}

func GetVideoCheckResult(bucketName string, videoId string, taskId string) (error, core.VideoPornResult) {
	startTime := time.Now().Unix()
	var pornResult core.VideoPornResult
	pornResult.Init()

	//根据返回的taskId查询视频检测结果 这里可能为processing
	checkRet := core.GetOSSObjectWithParams(bucketName, videoId, "x-oss-process=udf/green/video/result,"+taskId)
	var checkResult core.VideoPornOSSResult
	json.Unmarshal(checkRet, &checkResult)
	fmt.Printf("Get Video OSS Check Result - %v\n", checkResult.ToJSONStr())

	if checkResult.Code != 200 {
		pornResult.Code = checkResult.Code
		return fmt.Errorf(checkResult.Msg), pornResult
	} else if len(checkResult.Data.Results) > 0 {
		//封装视频检测结果输出
		pornResult.Ctime      = startTime
		pornResult.Code       = checkResult.Code
		pornResult.TaskId     = checkResult.RequestId
		pornResult.Scene      = checkResult.Data.Results[0].Scene
		pornResult.Label      = checkResult.Data.Results[0].Label
		pornResult.Rate       = checkResult.Data.Results[0].Rate
		pornResult.Suggestion = checkResult.Data.Results[0].Suggestion
		return nil, pornResult
	} else {
		pornResult.Code = checkResult.Code
		return fmt.Errorf("No Result"), pornResult
	}
}

