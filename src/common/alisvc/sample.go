package alisvc

import "fmt"
import "net/url"

func test() {
	urlstr := "http://8.37.236.234:8999/Fastdfsugcimg.php?group_name=la4group1&file_id=M05/0B/63/6QMAAFne3zGIS_4oAAAvMhmRAvIAARyoAI7s64AAC9K074.jpg"
	urlstr = url.QueryEscape(urlstr)
	content, _ := SvcPost("https://dtplus-cn-shanghai.data.aliyuncs.com/face/attribute", map[string]string{
		"type":      "0",
		"image_url": "http://img.ucweb.com/s/demo/g/logo.png?gyunoplist=,,jpeg;103," + urlstr + ";4,0x0-5x5,1,0,1",
	}, "UivNvn6XcKjb3FQc", "LVqK5NxgfMKbhfwfo6P9MjZZRosiYu")
	fmt.Println(content)
	return
	content1, _ := SvcPost_Old("http://green.aliyuncs.com", map[string]string{
		"Action":     "ImageDetection",
		"Version":    "2016-08-01",
		"Async":      "false",
		"ImageUrl.1": "http://image.uc.cn/s/demo/g/logo.png",
		"Scene.1":    "porn",
	}, "UivNvn6XcKjb3FQc", "LVqK5NxgfMKbhfwfo6P9MjZZRosiYu")
	fmt.Println(content1)
}
