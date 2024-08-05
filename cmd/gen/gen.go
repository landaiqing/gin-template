package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"strings"
)

const dbName = `schisandra-cloud-album`
const dbUser = "root"
const dbPwd = "LDQ20020618xxx"
const dbHost = "1.95.0.111"
const dbPort = "3306"
const MysqlConfig = dbUser + ":" + dbPwd + "@tcp(" + dbHost + ":" + dbPort + ")/" + string(dbName) + "?charset=utf8mb4&parseTime=True&loc=Local"

func main() {

	// 生成所有model和query
	processAllTables(initInfo())
	// 处理表关联关系
	processTableRelations(initInfo())

}

/**
 * 初始化数据库连接
 * 生成generator实例
 * 自定义字段的数据类型
 * 自定义模型结体字段的标签
 */
func initInfo() (db *gorm.DB, g *gen.Generator, fieldOpts []gen.ModelOpt) {
	// 连接数据库
	db, err := gorm.Open(mysql.Open(MysqlConfig), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(fmt.Errorf("数据库连接失败，请检查连接配置: %w", err))
	}

	// 生成实例
	g = gen.NewGenerator(gen.Config{
		// 相对执行`go run`时的路径, 会自动创建目录，相对路径为工程根目录
		OutPath: "query",

		// WithDefaultQuery 生成默认查询结构体(作为全局变量使用), 即`Q`结构体和其字段(各表模型)
		// WithoutContext 生成没有context调用限制的代码供查询
		// WithQueryInterface 生成interface形式的查询代码(可导出), 如`Where()`方法返回的就是一个可导出的接口类型
		Mode: gen.WithDefaultQuery | gen.WithQueryInterface,

		// 表字段可为 null 值时, 对应结体字段使用指针类型
		FieldNullable: true, // generate pointer when field is nullable

		// 表字段默认值与模型结构体字段零值不一致的字段, 在插入数据时需要赋值该字段值为零值的, 结构体字段须是指针类型才能成功, 即`FieldCoverable:true`配置下生成的结构体字段.
		// 因为在插入时遇到字段为零值的会被GORM赋予默认值. 如字段`age`表默认值为10, 即使你显式设置为0最后也会被GORM设为10提交.
		// 如果该字段没有上面提到的插入时赋零值的特殊需要, 则字段为非指针类型使用起来会比较方便.
		FieldCoverable: false, // generate pointer when field has default value, to fix problem zero value cannot be assign: https://gorm.io/docs/create.html#Default-Values

		// 模型结构体字段的数字类型的符号表示是否与表字段的一致, `false`指示都用有符号类型
		FieldSignable: false, // detect integer field's unsigned type, adjust generated data type
		// 生成 gorm 标签的字段索引属性
		FieldWithIndexTag: false, // generate with gorm index tag
		// 生成 gorm 标签的字段类型属性
		FieldWithTypeTag: true, // generate with gorm column type tag
	})
	// 设置目标 db
	g.UseDB(db)

	// 自定义字段的数据类型
	// 统一数字类型为int64,兼容protobuf
	dataMap := map[string]func(columnType gorm.ColumnType) (dataType string){
		"tinyint":   func(columnType gorm.ColumnType) (dataType string) { return "int64" },
		"smallint":  func(columnType gorm.ColumnType) (dataType string) { return "int64" },
		"mediumint": func(columnType gorm.ColumnType) (dataType string) { return "int64" },
		"bigint":    func(columnType gorm.ColumnType) (dataType string) { return "int64" },
		"int":       func(columnType gorm.ColumnType) (dataType string) { return "int64" },
	}
	// 要先于`ApplyBasic`执行
	g.WithDataTypeMap(dataMap)

	// 自定义模型结体字段的标签
	// 将特定字段名的 json 标签加上`string`属性,即 MarshalJSON 时该字段由数字类型转成字符串类型
	// jsonField := gen.FieldJSONTagWithNS(func(columnName string) (tagContent string) {
	// 	toStringField := `balance, `
	// 	if strings.Contains(toStringField, columnName) {
	// 		return columnName + ",string"
	// 	}
	// 	return columnName
	// })

	// 将非默认字段名的字段定义为自动时间戳和软删除字段;
	// 自动时间戳默认字段名为:`updated_at`、`created_at, 表字段数据类型为: INT 或 DATETIME
	// 软删除默认字段名为:`deleted_at`, 表字段数据类型为: DATETIME
	//autoUpdateTimeField := gen.FieldGORMTag("update_time", func(tag field.GormTag) field.GormTag {
	//	return tag.Append("autoUpdateTime")
	//})
	//autoCreateTimeField := gen.FieldGORMTag("created_time", func(tag field.GormTag) field.GormTag {
	//	return tag.Append("autoCreateTime")
	//})
	//softDeleteField := gen.FieldType("deletedAt", "gorm.DeletedAt")

	// 模型自定义选项组
	fieldOpts = []gen.ModelOpt{
		// jsonField,
		//autoCreateTimeField,
		//autoUpdateTimeField,
		//softDeleteField,
	}

	return db, g, fieldOpts
}

/**
 *创建全部模型文件，生成所有model和query
 *将生成的query目录下的gen.go文件移动到当前目录tmp文件夹下
 */
func processAllTables(db *gorm.DB, g *gen.Generator, fieldOpts []gen.ModelOpt) {
	allModel := g.GenerateAllTable(fieldOpts...)
	g.ApplyBasic(allModel...)
	g.Execute()

	// 将生成的query目录下的gen.go文件移动到当前目录tmp文件夹下
	moveGenFile()
}

/**
 * 将生成的query目录下的gen.go根文件移动到当前目录tmp文件夹下，
 * gen.go文件中保存的是所有表的模型的引用，
 * gen在生成query文件时，只会将ApplyBasic方法参数中的模型写入query中的根文件gen.go中，
 * 而我们在后续调用processTableRelations方法处理关联关系的时候，只处理有关联关系的表，
 * 方法中生成的gen.go中只会有有关联关系的表的模型的引用，因此需要将保存了所有表的模型的引用的gen.go文件
 * 移动到tmp文件夹下，然后再调用processTableRelations方法处理关联关系，处理完关联关系后，
 * 再将tmp文件夹下的gen.go文件移动到query目录下。
 */
func moveGenFile() {
	workDir, _ := os.Getwd()
	err := os.MkdirAll(workDir+"/tmp", 0777)
	if err != nil {
		fmt.Println("创建文件夹logs失败!", err)
		return
	}
	genFile := workDir + "/query/gen.go"
	if _, err := os.Stat(genFile); err != nil {
		fmt.Println("gen.go文件不存在!")
		return
	}
	fmt.Println("gen.go文件存在:", genFile)
	er := os.Rename(genFile, workDir+"/tmp/gen.go")
	if er != nil {
		return
	}
}

type Results struct {
	TABLE_NAME             string //子表名
	COLUMN_NAME            string //子表列名
	CONSTRAINT_NAME        string //约束名
	REFERENCED_TABLE_NAME  string //关联表名
	REFERENCED_COLUMN_NAME string //关联列名
}

/**
 * 处理表关联关系
 */
func processTableRelations(db *gorm.DB, g *gen.Generator, fieldOpts []gen.ModelOpt) {
	relationList := []Results{}
	// 执行这条sql语句，获取当前数据库中所有表之间的外键关联关系
	// 执行结果保存到relationList中
	db.Raw("SELECT TABLE_NAME,COLUMN_NAME,CONSTRAINT_NAME,REFERENCED_TABLE_NAME,REFERENCED_COLUMN_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE TABLE_SCHEMA = ? AND REFERENCED_TABLE_SCHEMA IS NOT NULL;", dbName).Scan(&relationList)

	type subTable struct {
		TABLE_NAME     string //子表名
		TABLE_NAME_UP  string //子表名首字母大写
		COLUMN_NAME    string //子表列名
		COLUMN_NAME_UP string //子表列名首字母大写
	}

	masterTableMap := make(map[string][]subTable)
	// 将relationList中的数据按照关联表名进行分组，将关联了父表名的所有子表数据放到一个切片中，然后将切片放到map中，map的key为父表名，value为子表切片
	for _, sub := range relationList {
		st := subTable{
			TABLE_NAME:     sub.TABLE_NAME,                         //子表名
			COLUMN_NAME:    sub.COLUMN_NAME,                        //子表列名
			TABLE_NAME_UP:  Case2Camel(sub.TABLE_NAME),             //将子表名下划线去掉，转换成首字母大写
			COLUMN_NAME_UP: Case2Camel(ProcessID(sub.COLUMN_NAME)), //将子表列名中以id结尾的字段中的id转换成ID格式，再将子表列名下划线去掉，转换成首字母大写
		}
		masterTableMap[sub.REFERENCED_TABLE_NAME] = append(masterTableMap[sub.REFERENCED_TABLE_NAME], st)
	}

	fmt.Println("主表 Map:::", ToJson(masterTableMap))

	// 生成新的generator实例，用于通过数据库子表名称，创建子表的模型基本结构体（BaseStruct）
	newGenerator := gen.NewGenerator(gen.Config{})
	newGenerator.UseDB(db)

	relationModels := []any{}
	// 遍历map，将map中的数据取出来，生成对应的关联关系模型文件
	for masterTable, subTables := range masterTableMap {
		subModels := []gen.ModelOpt{}
		// 遍历子表切片，将子表切片中的数据取出来，生成对应的关联关系模型文件
		for _, subTable := range subTables {
			// 目前只支持一对多关联关系，即：HasMany
			// 但是也能覆盖has_one和belongs_to的关联关系，只不过在生成的model中会多出一个切片字段，该切片中只有一个值
			// 对于多对多关联关系(many2many)，请先设计中间连接表，连接表中定义两个主键，即：复合主键，每个主键关联一张主表，
			// 这样就能生成两个一对多的关联关系，再运行本程序，就能实现多对多的关联关系了
			subModels = append(subModels, gen.FieldRelate(field.HasMany, subTable.TABLE_NAME_UP, newGenerator.GenerateModel(subTable.TABLE_NAME),
				&field.RelateConfig{
					// RelateSlice配置为true，那么在主表生成model的时候会生成关联表的切片
					RelateSlice: true,
					// 配置关联关系的外键字段，并且将外键字段的gorm标签中的foreignKey属性设置为关联表的列名
					GORMTag: field.GormTag{"foreignKey": {subTable.COLUMN_NAME_UP}},
				}))
		}
		relationModels = append(relationModels, g.GenerateModel(masterTable, append(fieldOpts, subModels...)...))
	}

	g.ApplyBasic(relationModels...)
	g.Execute()

	// 将当前目录tmp文件夹下的gen.go文件移动到query目录下
	moveGenFileBack()
}

/**
 * 将当前目录tmp文件夹下的gen.go文件移动到query目录下
 */
func moveGenFileBack() {
	workDir, _ := os.Getwd()
	genFile := workDir + "/query/gen.go"
	if _, err := os.Stat(genFile); err != nil {
		fmt.Println("gen.go文件不存在!")
		return
	}
	err := os.Rename(workDir+"/tmp/gen.go", genFile)
	if err != nil {
		return
	}
}

// Case2Camel 下划线写法转为驼峰写法
func Case2Camel(name string) string {
	words := strings.Split(name, "_")
	var result string
	for _, word := range words {
		result += strings.ToUpper(string(word[0])) + word[1:]
	}
	return result
}

func ProcessID(str string) string {
	if strings.HasSuffix(str, "id") {
		str, _ = strings.CutSuffix(str, "id")
		str = str + "ID"
	}
	return str
}

func ToJson(result interface{}) string {
	jsonBytes, _ := json.Marshal(result)
	return string(jsonBytes)
}
