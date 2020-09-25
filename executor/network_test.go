package executor

import (
	"encoding/json"
	"testing"
)

func TestUserInfo(t *testing.T) {
	userInfoRaw := `{"uploadinfo":"","user_id":31395269,"app_version":0,"app_id":0,"userkey":"8684E264F90403CBE72F173843D2D83476287E67","url_upload":"http:\/\/uplb.115.com\/2.0\/upload.php","url_resume":"http:\/\/uplb.115.com\/2.0\/resume.php","url_cancel":"http:\/\/uplb.115.com\/2.0\/cancel.php","url_speed":"http:\/\/119.147.156.144\/2.0\/bigupload","url_speed_test":{"1":"http:\/\/119.147.156.144\/ST","2":"http:\/\/58.253.94.207\/ST"},"size_limit":123480309760,"size_limit_yun":209715200,"max_dir_level":25,"max_dir_level_yun":25,"max_file_num":50000,"max_file_num_yun":10000,"upload_allowed":true,"upload_allowed_msg":"","type_limit":["doc","docx","xls","pdf","ppt","wps","dps","et","mdb","reg","txt","wri","rtf","lrc","vob","sub","srt","ass","ssa","idx","jar","umd","xlsx",".xlsm","xltx","xltm","xlam","xlsb","odt","pptx","ods","odp","chm","pot","pps","ppsx"],"file_range":{"2":"0-67108864","1":"67108864-0"},"isp_type":0,"state":true,"error":"","errno":0}`
	jsonUserInfo := new(userInfo)

	err := json.Unmarshal([]byte(userInfoRaw), jsonUserInfo)
	if err != nil {
		t.Error()
	}
}
