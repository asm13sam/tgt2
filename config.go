package main

var (
	tablesColumnsNum    map[string]int
	tablesColumns       map[string][]TableColumn
	tablesColumnsRaw    map[string][]TableColumn
	tablesColumnsRawMap map[string](map[string]TableColumn)
)

var typesTemplates = []string{
	"\t\"%s\": %s,\n",
	"\t\"%s\": \"%s\",\n",
	"\t\"%s\": %s,\n",
	"\t\"%s\": %s,\n",
}

var operators = map[string]string{
	"eq":  "=",
	"lt":  "<",
	"gt":  ">",
	"lte": "<=",
	"gte": ">=",
}

type TableColumn struct {
	cname string
	ctype int
}

func readConfig() {
	tablesColumnsNum = map[string]int{"measure": 4, "matherial_group": 5, "matherial": 12}
	tablesColumnsRaw = map[string][]TableColumn{
		"measure": {
			TableColumn{"id", INT},
			TableColumn{"name", STRING},
			TableColumn{"full_name", STRING},
			TableColumn{"is_active", BOOL},
		},
		"matherial_group": {
			TableColumn{"id", INT},
			TableColumn{"name", STRING},
			TableColumn{"matherial_group_id", INT},
			TableColumn{"position", INT},
			TableColumn{"is_active", BOOL},
		},
		"matherial": {
			TableColumn{"id", INT},
			TableColumn{"name", STRING},
			TableColumn{"full_name", STRING},
			TableColumn{"matherial_group_id", INT},
			TableColumn{"measure_id", INT},
			TableColumn{"color_group_id", INT},
			TableColumn{"price", FLOAT},
			TableColumn{"cost", FLOAT},
			TableColumn{"total", FLOAT},
			TableColumn{"barcode", STRING},
			TableColumn{"count_type_id", INT},
			TableColumn{"is_active", BOOL}},
	}
	tablesColumnsRawMap = map[string](map[string]TableColumn){
		"measure": {
			"id":        TableColumn{"id", INT},
			"name":      TableColumn{"name", STRING},
			"full_name": TableColumn{"full_name", STRING},
			"is_active": TableColumn{"is_active", BOOL},
		},
		"matherial_group": {
			"id":                 TableColumn{"id", INT},
			"name":               TableColumn{"name", STRING},
			"matherial_group_id": TableColumn{"matherial_group_id", INT},
			"position":           TableColumn{"position", INT},
			"is_active":          TableColumn{"is_active", BOOL},
		},
		"matherial": {
			"id":                 TableColumn{"id", INT},
			"name":               TableColumn{"name", STRING},
			"full_name":          TableColumn{"full_name", STRING},
			"matherial_group_id": TableColumn{"matherial_group_id", INT},
			"measure_id":         TableColumn{"measure_id", INT},
			"color_group_id":     TableColumn{"color_group_id", INT},
			"price":              TableColumn{"price", FLOAT},
			"cost":               TableColumn{"cost", FLOAT},
			"total":              TableColumn{"total", FLOAT},
			"barcode":            TableColumn{"barcode", STRING},
			"count_type_id":      TableColumn{"count_type_id", INT},
			"is_active":          TableColumn{"is_active", BOOL}},
	}

	tablesColumns = map[string][]TableColumn{
		"measure": {
			TableColumn{"id", INT},
			TableColumn{"name", STRING},
			TableColumn{"full_name", STRING},
			TableColumn{"is_active", BOOL}},
		"matherial_group": {
			TableColumn{"id", INT},
			TableColumn{"name", STRING},
			TableColumn{"matherial_group_id", INT},
			TableColumn{"position", INT},
			TableColumn{"is_active", BOOL},
			TableColumn{"matherial_group", STRING},
		},
		"matherial": {
			TableColumn{"id", INT},
			TableColumn{"name", STRING},
			TableColumn{"full_name", STRING},
			TableColumn{"matherial_group_id", INT},
			TableColumn{"measure_id", INT},
			TableColumn{"color_group_id", INT},
			TableColumn{"price", FLOAT},
			TableColumn{"cost", FLOAT},
			TableColumn{"total", FLOAT},
			TableColumn{"barcode", STRING},
			TableColumn{"count_type_id", INT},
			TableColumn{"is_active", BOOL},
			TableColumn{"matherial_group", STRING},
			TableColumn{"measure", STRING},
			TableColumn{"color_group", STRING},
			TableColumn{"count_type", STRING},
		},
	}
}
