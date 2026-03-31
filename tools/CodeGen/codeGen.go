package main

/*
 * @Desc: 生成 GORM 表结构
 * @author: 福狼
 * @version: v1.0.0
 */

import (
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func main() {
	//dbName := "flf"
	//dbUser := "postgres"
	//dbPassword := "GFD@Postgresql#176!"
	//dbHost := "148.70.18.111"
	//dbPort := "45432"

	dbName := "gfg"
	dbUser := "postgres"
	dbPassword := "123456"
	dbHost := "192.168.153.121"
	dbPort := "5432"
	myDBConf := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	generate(myDBConf, "gfg_prize_member")
}

func generate(sqlConf string, tables ...string) {
	// 使用 PGSQL 驱动
	db, err := gorm.Open(postgres.Open(sqlConf))
	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "./query",
		// WithoutContext 生成没有context调用限制的代码供查询
		Mode: gen.WithoutContext,
		// 表字段可为 null 时对应结构体字段使用指针类型
		FieldNullable:  true,
		FieldCoverable: false,
		FieldSignable:  false,
		// 生成 gorm 标签的字段索引属性
		FieldWithIndexTag: false,
		// 生成 gorm 标签的字段类型属性
		FieldWithTypeTag: true,
	})

	g.UseDB(db)

	jsonField := gen.FieldJSONTagWithNS(func(columnName string) (tagContent string) {
		appendString := false
		if strings.Contains(columnName, "_id") {
			appendString = true
		}
		camelNameArray := strings.Split(columnName, "_")
		for i, _ := range camelNameArray {
			if i == 0 {
				continue
			}
			camelNameArray[i] = strings.Title(camelNameArray[i])
		}
		camelName := strings.Join(camelNameArray, "")
		if appendString {
			return camelName + ",string"
		} else {
			return camelName
		}
	})

	autoUpdateTimeField := gen.FieldGORMTag("update_time", func(tag field.GormTag) field.GormTag {
		tag.Set("column", "update_time")
		tag.Set("type", "int", "unsigned")
		tag.Set("autoUpdateTime")
		return tag
	})

	autoCreateTimeField := gen.FieldGORMTag("create_time", func(tag field.GormTag) field.GormTag {
		tag.Set("column", "create_time")
		tag.Set("type", "int", "unsigned")
		tag.Set("autoCreateTime")
		return tag
	})

	softDeleteField := gen.FieldType("delete_time", "soft_delete.DeletedAt")

	fieldOpts := []gen.ModelOpt{jsonField, autoCreateTimeField, autoUpdateTimeField, softDeleteField}

	if tables != nil {
		for _, table := range tables {
			model := g.GenerateModel(table, fieldOpts...)
			g.ApplyBasic(model)
		}
	} else {
		allModel := g.GenerateAllTable(fieldOpts...)
		g.ApplyBasic(allModel...)
	}

	g.Execute()
}
