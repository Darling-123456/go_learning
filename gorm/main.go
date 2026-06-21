package main

import (
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	//创建一个名为test.db的数据库实例
	//db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db, err := gorm.Open(mysql.Open("root:040725ge@tcp(localhost:3306)/test"))
	if err != nil {
		panic("连接数据库失败")
	}

	db = db.Debug()

	// 迁移schema值  自动建表
	db.AutoMigrate(&Product{})

	// 创建
	db.Create(&Product{Code: "D42", Price: 100})

	// 查询
	var product Product
	db.First(&product, 1)                 // 查找对应主键的产品
	db.First(&product, "code = ?", "D42") // 查找 code 为 D42 的所有产品

	// 更新 - 将产品价格更新为 200
	db.Model(&product).Update("Price", 200)
	// 更新 - 更新多个字段
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零字段
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// 删除 - 删除产品
	db.Delete(&product, 1)
}
