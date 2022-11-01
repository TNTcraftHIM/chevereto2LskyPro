// @Description lsky数据库模型
// @Author 小游, TNTcraftHIM
// @Date 2022/11/1
package model

type LskyUsers struct {
	Name            string  // 姓名
	Email           string  // 邮箱
	Password        string  // 密码
	IsAdminer       int     // 是否是管理员
	Capacity        float64 // 容量
	Configs         string  // 配置
	ImageNum        int     // 图片数量
	AlbumNum        int     // 相册数量
	RegisteredIp    string  // 注册ip
	EmailVerifiedAt int64   // 激活时间
}

type LskyAlbums struct {
	UserID    string // 用户
	Name      string // 名称
	Intro     string // 简介
	ImageNum  int    // 图片数量
	CreatedAt int64  // 创建时间
}

type LskyImages struct {
	UserID      string // 所属用户
	AlbumID     string // 所属相册
	StrategyID  int    // 存储策略
	Key         string // key
	Path        string // 保存路径
	Name        string // 保存名称
	OriginName  string // 原始名称
	Size        string // 图片大小
	MimeType    string // 图片类型
	Extension   string // 图片后缀
	Md5         string // 图片 hash md5加密
	Sha1        string // 图片 hash sha1加密
	Width       string // 图片宽度
	Height      string // 图片高度
	IsUnhealthy string // 是否不健康
	UploadedIP  string // 上传者ip
	CreatedAt   int64  // 创建时间
}
