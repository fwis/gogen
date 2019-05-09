package gendb

func RenderDMgo(DMTemplateFilePath string, MmPackage string, m *ModelSQLGen) ([]byte, error) {
	vars := make(map[string]interface{})

	name := m.Struct.Name()
	vars["MmPackage"] = MmPackage
	vars["SQLGen"] = m

	vars["CreateTableSQL"] = m.CreateTableSQL()
	vars["Name"] = name
	vars["Table"] = m.TableName
	vars["Cols"] = m.Cols //m.Cols[0].FieldName
	vars["PKCol"] = m.PKCol
	//vars["UpperName"] = strings.ToUpper(name)
	//vars["TableName"] = m.TableName
	//vars["StructName"] = "mm." + m.Struct.Name() //mm. hardcode

	vars["HasCreateTime"] = m.HasCreateTime()
	vars["HasModifiedTime"] = m.HasModifiedTime()

	//common.SetConflictVars(vars)
	return RenderByTplFile(DMTemplateFilePath, vars)
}
