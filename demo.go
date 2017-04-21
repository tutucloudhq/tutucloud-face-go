package main

import (
	"fmt"
	"log"
	"github.com/tutucloudhq/tutucloud-face-go/tutucloud"
)

func main() {
	keys := tutucloud.Keys{
		// 公有 key
		Api_key   : "",
		// 私有 key
		Api_secret: "",
	}

	face      := &tutucloud.FaceApi{Keys: keys}

	image     := map[string]string{"url": "https://files.tusdk.com/img/faces/f-dd1.jpg"}
	// image     := map[string]string{"file": "test.webp"}
	params    := map[string]string{}
	data, err := face.Request("analyze/detection", image, params)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(data)
	}
}
