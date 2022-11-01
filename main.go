// @Description
// @Author 小游, TNTcraftHIM
// @Date 2022/11/1
package main

import (
	"chevereto2LskyPro/common"
	"chevereto2LskyPro/model"
	"chevereto2LskyPro/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"regexp"
	"strconv"
)

// 转换数据库
func changeData(startAt int) {
	// 初始化cheverto
	if !sql.InitDb1() {
		log.Fatalln("chevereto数据库初始化失败")
	}
	// 初始化lsky
	if !sql.InitDb2() {
		log.Fatalln("lsky数据库初始化失败")
	}
	prefix1 := common.GetConfigString("cheverto", "Prefix")
	prefix2 := common.GetConfigString("lsky", "Prefix")
	log.Println("开始转换数据库")

	// 查询并转换所有用户
	if data, err := sql.Db1Dql("SELECT user_id, COALESCE(user_name, ''), user_email, login_secret, user_is_admin, user_image_count, user_album_count, user_registration_ip, user_date FROM (SELECT * FROM " + prefix1 + "users JOIN " + prefix1 + "logins ON " + prefix1 + "logins.login_user_id = user_id AND " + prefix1 + "logins.login_type = 'password' ) AS p WHERE user_status = 'valid'"); err == nil && (startAt == 1 || startAt == 2) {
		capacity := common.GetConfigString("config", "Capacity")
		configs := common.GetConfigString("config", "Configs")
		groupID := common.GetConfigString("config", "GroupID")
		log.Printf("总计%d位用户, 开始转换\n", len(data))
		errNum := 0
		for k, v := range data {
			var lsky model.LskyUsers
			// 用户ID
			lsky.ID, err = strconv.Atoi(v[0])
			if err != nil {
				log.Printf("第%d位用户转换失败, 错误信息: %s\n", k, err.Error())
				errNum++
				continue
			}
			// 用户组
			lsky.GroupID, err = strconv.Atoi(groupID)
			if err != nil {
				lsky.GroupID = 1
			}
			// 姓名
			lsky.Name = v[1]
			// 邮箱
			lsky.Email = v[2]
			// 密码
			lsky.Password = v[3]
			// 是否是管理员
			lsky.IsAdminer, err = strconv.Atoi(v[4])
			if err != nil {
				lsky.IsAdminer = 0
			}
			// 容量
			lsky.Capacity, err = strconv.ParseFloat(capacity, 64)
			if err != nil {
				lsky.Capacity = 524288
			}
			// 配置
			lsky.Configs = configs
			// 图片数量
			lsky.ImageNum, err = strconv.Atoi(v[5])
			if err != nil {
				lsky.ImageNum = 0
			}
			// 相册数量
			lsky.AlbumNum, err = strconv.Atoi(v[6])
			if err != nil {
				lsky.AlbumNum = 0
			}
			// 注册ip
			lsky.RegisteredIp = v[7]
			// 注册时间
			lsky.EmailVerifiedAt = v[8]
			// 创建时间
			lsky.CreatedAt = v[8]
			// 更新时间
			lsky.UpdatedAt = v[8]

			// 开始插入
			success, err := sql.Db2Dml(`INSERT delayed INTO `+prefix2+`users
			(id,group_id,name,email,password,is_adminer,capacity,configs,image_num,album_num,registered_ip,email_verified_at,created_at,updated_at) 
			VALUES 
			(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, lsky.ID, lsky.GroupID, lsky.Name, lsky.Email, lsky.Password, lsky.IsAdminer, lsky.Capacity, lsky.Configs, lsky.ImageNum, lsky.AlbumNum, lsky.RegisteredIp, lsky.EmailVerifiedAt, lsky.CreatedAt, lsky.UpdatedAt)
			if success && err == nil {
				log.Printf("第%d位用户转换成功!\n", k)
			} else {
				errNum++
				log.Printf("第%d位用户转换失败, 错误信息: %s\n", k, err.Error())
				continue
			}
		}
		log.Printf("转换完成, 有%d转换失败!\n", errNum)
	} else {
		if err != nil {
			log.Println(err)
		}
	}
	// 查询并转换所有相册
	if data, err := sql.Db1Dql("SELECT album_id, album_user_id, album_name, COALESCE(album_description, ''), album_image_count, album_date FROM " + prefix1 + "albums WHERE album_user_id IS NOT NULL"); err == nil && startAt == 1 {
		log.Printf("总计%d个相册, 开始转换\n", len(data))
		errNum := 0
		for k, v := range data {
			var lsky model.LskyAlbums
			// 相册ID
			lsky.ID, err = strconv.Atoi(v[0])
			if err != nil {
				log.Printf("第%d个相册转换失败, 错误信息: %s\n", k, err.Error())
				errNum++
				continue
			}
			// 用户id
			lsky.UserID = v[1]
			// 相册名
			lsky.Name = v[2]
			// 相册描述
			lsky.Intro = v[3]
			// 图片数量
			lsky.ImageNum, err = strconv.Atoi(v[4])
			if err != nil {
				lsky.ImageNum = 0
			}
			// 创建时间
			lsky.CreatedAt = v[5]
			// 更新时间
			lsky.UpdatedAt = v[5]

			// 开始插入
			success, err := sql.Db2Dml(`INSERT delayed INTO `+prefix2+`albums
			(id, user_id,name,intro,image_num,created_at) 
			VALUES 
			(?,?,?,?,?,?)`, lsky.ID, lsky.UserID, lsky.Name, lsky.Intro, lsky.ImageNum, lsky.CreatedAt)
			if success && err == nil {
				log.Printf("第%d个相册转换成功!\n", k)
			} else {
				errNum++
				log.Printf("第%d个相册转换失败, 错误信息: %s\n", k, err.Error())
				continue
			}
		}
		log.Printf("转换完成, 有%d转换失败!\n", errNum)
	} else {
		if err != nil {
			log.Println(err)
		}
	}
	// 查询并转换所有的图片
	if data, err := sql.Db1Dql("SELECT COALESCE(image_user_id, ''), COALESCE(image_album_id, ''), image_date, image_name, image_original_filename, image_size, image_extension, image_md5, image_width, image_height, image_nsfw, image_uploader_ip FROM " + prefix1 + "images"); err == nil && startAt != 4 {
		log.Printf("总计%d张图片, 开始转换\n", len(data))
		errNum := 0
		for k, v := range data {
			var lsky model.LskyImages
			// 用户id
			if startAt == 3 {
				v[0] = ""
			}
			lsky.UserID = common.NewNullString(v[0])
			if lsky.UserID.Valid {
				lsky.UserID.String = v[0]
			}
			// 相册id
			if startAt == 3 {
				v[1] = ""
			}
			lsky.AlbumID = common.NewNullString(v[1])
			if lsky.AlbumID.Valid {
				lsky.AlbumID.String = v[1]
			}
			// 策略id
			lsky.StrategyID = 1
			// key
			lsky.Key = common.RandString(6)
			// 路径
			create := common.Str2time(v[2])
			lsky.Path = common.Time2String(create, false)
			// 图片名
			lsky.Name = v[3] + "." + v[6]
			// 原始文件名
			lsky.OriginName = v[4]
			// 大小
			var size float64
			size, err := strconv.ParseFloat(v[5], 64)
			if err != nil {
				size = 0.00
			}
			lsky.Size = strconv.FormatFloat(size/1024.00, 'f', 2, 64)
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
			lsky.CreatedAt = v[2]
			// 更新时间
			lsky.UpdatedAt = v[2]
			// 开始插入
			success, err := sql.Db2Dml(`INSERT delayed INTO `+prefix2+`images
			(user_id,album_id,strategy_id,`+prefix2+`images.key,path,name,origin_name,size,mimetype,extension,md5,sha1,width,height,is_unhealthy,uploaded_ip,created_at,updated_at) 
			VALUES 
			(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, lsky.UserID, lsky.AlbumID, lsky.StrategyID, lsky.Key, lsky.Path, lsky.Name, lsky.OriginName, lsky.Size, lsky.MimeType, lsky.Extension, lsky.Md5, "", lsky.Width, lsky.Height, lsky.IsUnhealthy, lsky.UploadedIP, lsky.CreatedAt, lsky.UpdatedAt)
			if success && err == nil {
				log.Printf("第%d张图片转换成功!\n", k)
			} else {
				errNum++
				log.Printf("第%d张图片转换失败, 错误信息: %s\n", k, err.Error())
				continue
			}
		}
		log.Printf("图片转换完成, 有%d转换失败\n", errNum)
	} else {
		if err != nil {
			log.Println(err)
		}
	}
	// 更新策略
	{
		log.Println("开始更新策略")
		config, err := sql.Db2Dql("SELECT configs from " + prefix2 + "strategies WHERE id = 1")
		if err != nil || config[0][0] == "" {
			log.Fatalln("策略更新失败, 错误信息: ", err)
		}
		// 解析json
		var data map[string]interface{}
		err = json.Unmarshal([]byte(config[0][0]), &data)
		if err != nil {
			log.Fatalln("策略更新失败, 错误信息: ", err)
		}
		host := data["url"].(string)
		pattern := regexp.MustCompile("^http(s)?://.+?/")
		host = pattern.FindString(host)
		if host == "" {
			log.Fatalln("策略更新失败, 错误信息: ", "无法找到域名")
		}
		targetPath := data["root"].(string)
		path := common.GetConfigString("img", "path") + "/public/images"
		data["url"] = host + "images"

		symErr := os.Symlink(targetPath, path)
		if symErr == nil {
			log.Println("软链接创建成功")
		} else {
			log.Println("软链接创建失败, 错误信息: ", symErr)
		}

		// 更新数据库
		json, DbErr := json.Marshal(data)
		if DbErr != nil {
			log.Fatalln("数据库更新失败, 错误信息: ", err)
		}
		_, DbErr = sql.Db2Dml("UPDATE "+prefix2+"strategies SET configs = ? WHERE id = 1", json)
		if DbErr == nil {
			log.Println("数据库更新成功")
		} else {
			log.Println("数据库更新失败, 错误信息: ", err)
		}

		if symErr == nil && DbErr == nil {
			log.Println("策略更新成功")
			return
		}
		log.Println("策略更新部分失败")
	}
	log.Println("转换全部完成!!!")
}

// 删除多余图片
func deleteMoreImage() {
	path, err := os.Readlink(common.GetConfigString("img", "path") + "/public/images")
	if err != nil {
		log.Fatalln("读取文件夹失败, 错误信息: ", err)
	}
	path += "/"
	list := new([]string)
	// 获取所有的文件
	if common.GetAllFile(path, list) != nil {
		log.Println("获取文件失败")
	}
	re := regexp.MustCompile(`\.[th|md].*?$`)
	log.Println("开始删除重复文件")
	total := 0
	errNUm := 0
	// 遍历删除文件夹
	for _, v := range *list {
		if re.MatchString(v) {
			if os.Remove(v) == nil {
				log.Println("删除" + v)
				total++
			} else {
				errNUm++
				log.Printf("第%d个文件删除失败, 错误信息: %s\n", total, err.Error())
			}
		}
	}
	// 打印成功
	log.Printf("已为你删除%d个文件, 删除失败%d个", total, errNUm)
}

// 修改文件权限
func changePermissions() {
	path, err := os.Readlink(common.GetConfigString("img", "path") + "/public/images")
	if err != nil {
		log.Fatalln("读取文件夹失败, 错误信息: ", err)
	}
	path += "/"
	list := new([]string)
	// 获取所有的文件
	if common.GetAllFile(path, list) != nil {
		log.Println("获取文件失败")
	}
	log.Println("开始修改文件权限")
	total := 0
	errNUm := 0
	// 遍历删除文件夹
	for _, v := range *list {
		if os.Chmod(v, fs.FileMode(0755)) == nil {
			log.Println("修改" + v)
			total++
		} else {
			errNUm++
			log.Printf("第%d个文件修改失败, 错误信息: %s\n", total, err.Error())
		}
	}
	// 打印成功
	log.Printf("已为你修改%d个文件的权限, 修改失败%d个", total, errNUm)
}

// 转换函数
func main() {
	// 读取配置文件中log的配置
	logConfig := common.GetConfigString("log", "path")
	var logOutput io.Writer
	// 设置日志文件
	if logConfig != "" {
		logFile, err := os.OpenFile(logConfig, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Error:", err)
		}
		defer logFile.Close()
		logOutput = io.MultiWriter(os.Stdout, logFile)
	} else {
		logConfig = "不保存日志文件"
		logOutput = os.Stdout
	}
	log.SetOutput(logOutput)
	fmt.Println("日志输出配置为: ", logConfig)
	var input int
	// 是否需要删除重复文件
	fmt.Printf("欢迎使用图床转换工具\n请选择操作(1转换数据库 2修改文件权限 3删除重复文件):")
	if _, err := fmt.Scan(&input); err == nil && input == 1 {
		fmt.Printf("请选择从哪个步骤开始(1从转移用户&相册开始 2从转移用户开始 3从转移图片开始 4从更新策略开始):")
		if _, err := fmt.Scan(&input); err == nil && (input == 1 || input == 2 || input == 3 || input == 4) {
			changeData(input)
		} else {
			log.Fatalln("输入错误")
		}
	} else if err == nil && input == 2 {
		changePermissions()
	} else if err == nil && input == 3 {
		deleteMoreImage()
	} else if err != nil {
		log.Fatalln(err)
	} else {
		log.Fatalln("输入错误")
	}
}
