// @Description
// @Author 小游, TNTcraftHIM
// @Date 2022/11/1
package main

import (
	"fmt"
	"img/common"
	"img/model"
	"img/sql"
	"os"
	"regexp"
	"strconv"
)

// 转换数据库
func changeData() {
	// 初始化cheverto
	sql.InitDb1()
	// 初始化lsky
	sql.InitDb2()
	prefix1 := common.GetConfigString("cheverto", "Prefix")
	prefix2 := common.GetConfigString("lsky", "Prefix")

	// 查询并转换所有用户
	if data, err := sql.Db1Dql("SELECT user_name, user_email, login_secret, user_is_admin, user_image_count, user_album_count, user_registration_ip, user_date_gmt FROM (SELECT * FROM " + prefix1 + "users JOIN " + prefix1 + "logins ON " + prefix1 + "logins.login_user_id = user_id AND " + prefix1 + "logins.login_type = 'password' ) AS p WHERE user_status = 'valid'"); err == nil {
		capacity := common.GetConfigString("config", "Capacity")
		configs := common.GetConfigString("config", "Configs")
		fmt.Printf("总计%d位用户, 开始转换\n", len(data))
		errNum := 0
		for k, v := range data {
			var lsky model.LskyUsers
			// 姓名
			lsky.Name = v[0]
			// 邮箱
			lsky.Email = v[1]
			// 密码
			lsky.Password = v[2]
			// 是否是管理员
			lsky.IsAdminer, err = strconv.Atoi(v[3])
			if err != nil {
				lsky.IsAdminer = 0
			}
			// 容量
			lsky.Capacity, err = strconv.ParseFloat(capacity, 64)
			if err != nil {
				lsky.Capacity = 0
			}
			// 配置
			lsky.Configs = configs
			// 图片数量
			lsky.ImageNum, err = strconv.Atoi(v[4])
			if err != nil {
				lsky.ImageNum = 0
			}
			// 相册数量
			lsky.AlbumNum, err = strconv.Atoi(v[5])
			if err != nil {
				lsky.AlbumNum = 0
			}
			// 注册ip
			lsky.RegisteredIp = v[6]
			// 注册时间
			lsky.EmailVerifiedAt = common.Str2time(v[7]).Unix()

			// 开始插入
			if sql.Db2Dml(`INSERT delayed INTO `+prefix2+`users
				(name,email,password,is_adminer,capacity,configs,image_num,album_num,registered_ip,email_verified_at) 
				VALUES 
				(?,?,?,?,?,?,?,?,?,?)`, lsky.Name, lsky.Email, lsky.Password, lsky.IsAdminer, lsky.Capacity, lsky.Configs, lsky.ImageNum, lsky.AlbumNum, lsky.RegisteredIp, lsky.EmailVerifiedAt) {
				fmt.Printf("第%d位用户转换成功!\n", k)
			} else {
				errNum++
			}
		}
		fmt.Printf("转换完成, 有%d转换失败!", errNum)
	} else {
		fmt.Println(err)
	}
	// 查询并转换所有相册
	if data, err := sql.Db1Dql("SELECT album_user_id, album_name, album_description, album_image_count, album_date_gmt FROM " + prefix1 + "albums"); err == nil {
		fmt.Printf("总计%d个相册, 开始转换\n", len(data))
		errNum := 0
		for k, v := range data {
			var lsky model.LskyAlbums
			// 用户id
			lsky.UserID = v[0]
			// 相册名
			lsky.Name = v[1]
			// 相册描述
			lsky.Intro = v[2]
			// 图片数量
			lsky.ImageNum, err = strconv.Atoi(v[3])
			if err != nil {
				lsky.ImageNum = 0
			}
			// 创建时间
			lsky.CreatedAt = common.Str2time(v[4]).Unix()

			// 开始插入
			if sql.Db2Dml(`INSERT delayed INTO `+prefix2+`users
				(user_id,name,intro,image_num,created_at) 
				VALUES 
				(?,?,?,?,?)`, lsky.UserID, lsky.Name, lsky.Intro, lsky.ImageNum, lsky.CreatedAt) {
				fmt.Printf("第%d个相册转换成功!\n", k)
			} else {
				errNum++
			}
		}
		fmt.Printf("转换完成, 有%d转换失败!", errNum)
	} else {
		fmt.Println(err)
	}
	// 查询并转换所有的图片
	if data, err := sql.Db1Dql("SELECT image_user_id, image_album_id, image_date, image_name, image_original_filename, image_size, image_extension, image_md5, image_width, image_height, image_nsfw, image_uploader_ip, image_date_gmt FROM " + prefix1 + "images"); err == nil {
		fmt.Printf("总计%d张图片, 开始转换\n", len(data))
		errNum := 0
		for k, v := range data {
			var lsky model.LskyImages
			// 用户id
			lsky.UserID = v[0]
			// 相册id
			lsky.AlbumID = v[1]
			// 策略id
			lsky.StrategyID = 1
			// key
			lsky.Key = common.RandString(6)
			// 路径
			create := common.Str2time(v[2])
			lsky.Path = common.Time2String(create, false)
			// 图片名
			lsky.Name = v[3]
			// 原始文件名
			lsky.OriginName = v[4]
			// 大小
			lsky.Size = v[5]
			// 类型
			switch v[6] {
			case "png":
				lsky.MimeType = "image/png"
			case "jpg":
				lsky.MimeType = "image/jpeg"
			case "jpeg":
				lsky.MimeType = "image/jpeg"
			case "gif":
				lsky.MimeType = "image/gif"
			case "ico":
				lsky.MimeType = "image/x-icon"
			default:
				lsky.MimeType = "image/png"
			}
			// 扩展名
			lsky.Extension = v[6]
			// md5加密值
			lsky.Md5 = v[7]
			// 宽度
			lsky.Width = v[8]
			// 高度
			lsky.Height = v[9]
			// 是否为不健康
			lsky.IsUnhealthy = v[10]
			// 上传ip
			lsky.UploadedIP = v[11]
			// 上传时间
			lsky.CreatedAt = common.Str2time(v[12]).Unix()
			// 开始插入
			if sql.Db2Dml(`INSERT delayed INTO `+prefix2+`images
			(user_id,album_id,strategy_id,key,path,name,origin_name,size,mimetype,extension,md5,sha1,width,height,is_unhealthy,uploaded_ip,created_at) 
			VALUES 
			(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, lsky.UserID, lsky.AlbumID, lsky.StrategyID, lsky.Key, lsky.Path, lsky.Name, lsky.OriginName, lsky.Size, lsky.MimeType, lsky.Extension, lsky.Md5, lsky.Md5, lsky.Width, lsky.Height, lsky.IsUnhealthy, lsky.UploadedIP, lsky.CreatedAt) {
				fmt.Printf("第%d张图片转换成功!\n", k)
			} else {
				errNum++
			}
		}
		fmt.Printf("图片转换完成！！！有%d转换失败", errNum)
	} else {
		fmt.Println(err)
	}
}

// 删除多余图片
func deleteMoreImage() {
	path := common.GetConfigString("img", "path")
	list := new([]string)
	// 获取所有的文件
	if common.GetAllFile(path, list) != nil {
		fmt.Println("获取文件失败")
	}
	re, _ := regexp.Compile(`\.[th|md].*?$`)
	total := 0
	errNUm := 0
	// 遍历删除文件夹
	for _, v := range *list {
		if re.MatchString(v) {
			if os.Remove(v) == nil {
				fmt.Println("删除" + v)
				total++
			} else {
				errNUm++
			}
		}
	}
	// 打印成功
	fmt.Printf("已为你删除%d个文件, 删除失败%d个", total, errNUm)
}

// 转换函数
func main() {
	var input int
	// 是否需要清空数据库
	fmt.Printf("欢迎使用图床转换工具\n请选择操作(1转换数据库 2删除重复文件):")
	if _, err := fmt.Scan(&input); err == nil && input == 1 {
		changeData()
	} else if err == nil && input == 2 {
		deleteMoreImage()
	}
}
