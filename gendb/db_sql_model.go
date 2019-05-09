package gendb

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type DataType int

const (
	BIGINT = DataType(1)
	INT    = DataType(2)
	BOOL   = DataType(3)
	STRING = DataType(4)
	FLOAT  = DataType(5)
	STRUCT = DataType(100)
)

var IGNORE_UPDATE_COLNAMES = []string{"CreateTime", "ModifiedTime"}

type Col struct {
	ColName     string
	ColDataType DataType
	GoTypeName  string
	NotAutoInc  bool //默认主键是自增长的
	FieldName   string
	Comment     string //字段注释
}

func (m *Col) IsNum() bool {
	return (m.ColDataType == BIGINT || m.ColDataType == INT)
}

type ModelSQLGen struct {
	Struct            reflect.Type
	TableName         string
	PKCol             *Col
	Cols              []*Col
	UpdateCols        []*Col //no PK, createtime, modifiedtime
	_max_col_name_len int
	_colsSQL          string
	_updateColsSQL    string
	_placeholdSQL     string
	_createTableSQL   string
	_insertCode       string
	_queryByPK_SQL    string
	_deleteByPK_SQL   string
	_has_createtime   int //0:还未check,1:有,2:没有
	_has_modifiedtime int //0:还未check,1:有,2:没有
}

func (m *ModelSQLGen) HasCreateTime() bool {
	if m._has_createtime == 0 {
		for _, c := range m.Cols {
			if c.FieldName == "CreateTime" && c.ColDataType == BIGINT {
				m._has_createtime = 1
				break
			}
		}
		if m._has_createtime == 0 {
			m._has_createtime = 2
		}
	}
	return (m._has_createtime == 1)
}

func (m *ModelSQLGen) HasModifiedTime() bool {
	if m._has_modifiedtime == 0 {
		for _, c := range m.Cols {
			if c.FieldName == "ModifiedTime" && c.ColDataType == BIGINT {
				m._has_modifiedtime = 1
				break
			}
		}
		if m._has_modifiedtime == 0 {
			m._has_modifiedtime = 2
		}
	}
	return (m._has_modifiedtime == 1)
}

func (m *ModelSQLGen) JoinCols(alias string, hasPK, hasTimes bool) string {
	s := ""

	if alias != "" && !strings.HasSuffix(alias, ".") {
		alias = alias + "."
	}

	for _, col := range m.Cols {
		if !hasPK && col.ColName == m.PKCol.ColName {
			continue
		}

		if !hasTimes && (col.ColName == "CreateTIme" || col.ColName == "ModifiedTime") {
			continue
		}

		if s == "" {
			s += alias + col.ColName
		} else {
			s += "," + alias + col.ColName
		}
	}
	return s
}

func (m *ModelSQLGen) ColsSQL() string {
	if m._colsSQL == "" {
		s := ""
		for _, col := range m.Cols {
			if s == "" {
				s += col.ColName
			} else {
				s += "," + col.ColName
			}
		}
		m._colsSQL = s
	}

	return m._colsSQL
}

func isIgnoreUpdateCol(col *Col) bool {
	for _, ignore_update_colname := range IGNORE_UPDATE_COLNAMES {
		if col.ColName == ignore_update_colname {
			return true
		}
	}
	return false
}

func (m *ModelSQLGen) UpdateColsSQL() string {
	if m._updateColsSQL == "" {
		s := ""
		for _, col := range m.UpdateCols {
			if col.ColName == m.PKCol.ColName || isIgnoreUpdateCol(col) {
				continue
			}

			if s == "" {
				s += col.ColName + "=?"
			} else {
				s += "," + col.ColName + "=?"
			}
		}
		m._updateColsSQL = s
	}

	return m._updateColsSQL
}

func (m *ModelSQLGen) PlaceholdSQL() string {
	if m._placeholdSQL == "" {
		n := len(m.Cols)
		s := strings.Repeat("?,", n)
		m._placeholdSQL = s[:(len(s) - 1)]
	}
	return m._placeholdSQL
}

func (m *ModelSQLGen) QueryByPK_SQL() string {
	if m._queryByPK_SQL == "" {
		s := "select " + m.ColsSQL() + " from " + m.TableName + " where " + m.PKCol.ColName + "=?"
		m._queryByPK_SQL = s
	}
	return m._queryByPK_SQL
}

func (m *ModelSQLGen) DeleteByPK_SQL() string {
	if m._deleteByPK_SQL == "" {
		s := "delete from " + m.TableName + " where " + m.PKCol.ColName + "=?"
		m._deleteByPK_SQL = s
	}
	return m._deleteByPK_SQL
}

func (m *ModelSQLGen) CreateTableSQL() string {
	if m._createTableSQL == "" {
		islastcol := false
		cols_count := len(m.Cols)
		s := "CREATE TABLE " + m.TableName + " (\n"
		for i, col := range m.Cols {
			if (i + 1) == cols_count {
				islastcol = true
			}
			col_format := "%-" + strconv.Itoa(m._max_col_name_len+1) + "s"
			s += "    " + fmt.Sprintf(col_format, col.ColName)
			if col.ColDataType == 0 {
				continue
			} else if col.ColDataType == INT {
				s += "INT NOT NULL"
				if i == 0 {
					s += " PRIMARY KEY"
					if !col.NotAutoInc {
						s += " AUTO_INCREMENT"
					}
				} else {
					s += " DEFAULT 0"
				}
			} else if col.ColDataType == BIGINT {
				s += "BIGINT NOT NULL"
				if i == 0 {
					s += " PRIMARY KEY"
					if !col.NotAutoInc {
						s += " AUTO_INCREMENT"
					}
				} else {
					s += " DEFAULT 0"
				}
			} else if col.ColDataType == BOOL {
				s += "BOOL NOT NULL"
				if i == 0 {
					s += " PRIMARY KEY"
				} else {
					s += " DEFAULT FALSE"
				}
			} else if col.ColDataType == FLOAT {
				s += "FLOAT NOT NULL"
				if i == 0 {
					s += " PRIMARY KEY"
				} else {
					s += " DEFAULT 0"
				}
			} else if col.ColDataType == STRING {
				s += "VARCHAR(255?) NOT NULL"
				if i == 0 {
					s += " PRIMARY KEY"
				} else {
					s += " DEFAULT ''"
				}
			}
			if !islastcol {
				s += ",\n"
			} else {
				s += "\n"
			}
		}
		s += ") engine=InnoDB charset=UTF8;\n"

		m._createTableSQL = s
	}

	return m._createTableSQL
}

func (m *ModelSQLGen) InsertCode() string {
	into_sql := ""
	values_sql := ""
	field_code := ""

	if m._insertCode == "" {
		for i, col := range m.Cols {
			if i == 0 {
				into_sql += col.ColName
				values_sql += "?"
				field_code += "    obj." + col.FieldName + ",\n"
			} else {
				into_sql += "," + col.ColName
				values_sql += ",?"
				field_code += "    obj." + col.FieldName + ",\n"
			}
		}
		insert_sql := "insert into " + m.TableName
		insert_sql += "(" + into_sql + ") "
		insert_sql += "values(" + values_sql + ")"
		m._insertCode = fmt.Sprintf(`
    insert_sql := "%s"
    m.db.Exec(insert_sql,
%s
    )
`, insert_sql, field_code)
	}
	return m._insertCode
}

func parseField_PrimitiveDataType(rf_field reflect.StructField) *Col {
	if rf_field.PkgPath != "" { // unexported
		return nil
	}

	colDataType := GetDBTypeByKind(rf_field.Type.Kind())

	col := &Col{}
	col.FieldName = rf_field.Name
	col.ColName = col.FieldName
	col.GoTypeName = rf_field.Type.Name()

	dbtag_str := rf_field.Tag.Get("db")
	dbtags := strings.Split(dbtag_str, ",")

	for i, dbtag := range dbtags {
		if i == 0 {
			if dbtag == "-" {
				continue
			} else if dbtag != "" {
				col.ColName = dbtag
			}
		} else {
			if dbtag == "notautoinc" {
				col.NotAutoInc = true
			} else if dbtag == "pk" {
				//
			} else {
				panic("Unknow db tag attr: " + dbtag)
			}
		}
	}

	col.ColDataType = colDataType

	return col

}

func parseField(t reflect.Type, afterParseField func(col *Col)) {
	for i := 0; i < t.NumField(); i++ {
		rf_field := t.Field(i) //rf means reflect

		if rf_field.PkgPath != "" { // unexported
			continue
		}

		ColDataType := GetDBTypeByKind(rf_field.Type.Kind())
		if ColDataType == STRUCT {
			if rf_field.Name == rf_field.Type.Name() {
				parseField(rf_field.Type, afterParseField)
			} else {
				panic("不能嵌套带名的Struct")
			}
		} else {
			col := parseField_PrimitiveDataType(rf_field)
			if col != nil {
				afterParseField(col)
			}
		}
	}
}

func CreateSQLGenByObj(obj interface{}, tableName string, pk_col_name string) *ModelSQLGen {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		panic("obj必须是struct才能生成ModelSQLGen")
	}

	cols := []*Col{}
	update_cols := []*Col{}
	max_col_name_len := 0

	parseField(t, func(col *Col) {
		if len(col.ColName) > max_col_name_len {
			max_col_name_len = len(col.ColName)
		}
		cols = append(cols, col)
	})

	for _, col := range cols {
		if !isIgnoreUpdateCol(col) {
			update_cols = append(update_cols, col)
		}
	}

	pkcol := cols[0]

	if pkcol.ColDataType != INT && pkcol.ColDataType != BIGINT {
		//panic("Struct first field is PK, must be integer")
		pkcol.NotAutoInc = true
	}

	if pk_col_name != "" && pk_col_name != pkcol.ColName {
		panic("obj struct 首成员必须是表的主键")
	}

	//排除 update_cols 中的 PKCol
	if update_cols[0].ColName == pkcol.ColName {
		update_cols = update_cols[1:]
	}

	return &ModelSQLGen{
		Struct:            t,
		TableName:         tableName,
		PKCol:             pkcol,
		Cols:              cols,
		UpdateCols:        update_cols,
		_max_col_name_len: max_col_name_len,
	}
}

func GetDBTypeByKind(kind reflect.Kind) DataType {
	if kind == 0 {
		return 0
	}

	switch kind {
	case reflect.Invalid:
		return 0
	case reflect.Bool:
		return BOOL
	case reflect.Int:
		return INT
	case reflect.Int8:
		return INT
	case reflect.Int16:
		return INT
	case reflect.Int32:
		return INT
	case reflect.Int64:
		return BIGINT
	case reflect.Uint8:
		return INT
	case reflect.Uint16:
		return INT
	case reflect.Uint32:
		return INT
	case reflect.Uint64:
		return BIGINT
	case reflect.Float32:
		return FLOAT
	case reflect.Float64:
		return FLOAT
	case reflect.String:
		return STRING
	case reflect.Struct:
		return STRUCT
	default:
		return 0
	}
}
