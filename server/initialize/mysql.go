package initialize

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"github.com/go-spring/go-spring/spring-boot"
	_ "github.com/go-spring/go-spring/starter-mysql-gorm"
	"github.com/jinzhu/gorm"
)

// 初始化数据库并产生数据库全局变量
func Mysql() {
	SpringBoot.Config(func(db *gorm.DB, maxIdleConns, maxOpenConns int, logMode bool) {
		db.DB().SetMaxIdleConns(maxIdleConns)
		db.DB().SetMaxOpenConns(maxOpenConns)
		db.LogMode(logMode)
		DBTables(db)
	}, "", "${db.max-idle-conns}", "${db.max-open-conns}", "${db.log-mode}")
}

// 注册数据库表专用
func DBTables(db *gorm.DB) {
	db.AutoMigrate(model.SysUser{},
		model.SysAuthority{},
		model.SysApi{},
		model.SysBaseMenu{},
		model.JwtBlacklist{},
		model.SysWorkflow{},
		model.SysWorkflowStepInfo{},
		model.ExaFileUploadAndDownload{},
		model.ExaFile{},
		model.ExaFileChunk{},
		model.ExaCustomer{},
	)
	global.GVA_LOG.Debug("register table success")
}
