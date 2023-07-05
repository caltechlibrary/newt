package newt

import (
	"fmt"
	"strings"
	"time"
)


var (
	// SQLTypes is a map of dsl types to SQL types.
	SQLTypes = map[string]string{
		"String":   "VARCHAR(256) DEFAULT ''",
		"Text":     "TEXT DEFAULT ''",
		"Integer":  "INTEGER DEFAULT 0",
		"Real":     "REAL DEFAULT 0.0",
		"Boolean":  "BOOLEAN",
		"Date":     "DATE",
		"Year":     "INTEGER DEFAULT 0",
		"Month":    "INTEGER DEFAULT 0",
		"Day":      "INTEGER DEFAULT 0",
		"Basename": "VARCHAR(256) DEFAULT ''",
		"Extname":  "VARCHAR(6) DEFAULT ''",
		"ISBN10":   "VARCHAR(10) DEFAULT ''",
		"ISBN13":   "VARCHAR(13) DEFAULT ''",
		"ISBN":     "VARCHAR(13) DEFAULT ''",
		"ISSN":     "VARCHAR(9) DEFAULT ''",
		"DOI":      "VARCHAR(256) DEFAULT ''",
		"ISNI":     "VARCHAR(16) DEFAULT ''",
		"ORCID":    "VARCHAR(19) DEFAULT ''",
		"ArXiv":    "VARCHAR(10) DEFAULT ''",
		"Markdown": "TEXT DEFAULT ''",
	}
)

func commentSection(configFName string, modelName string) []byte {
    now := time.Now()
	parts := []string{ 
		"--", 
	}
	if modelName != "" {
	    parts = append(parts, fmt.Sprintf("-- Model: %s", modelName))
	}
	if configFName != "" {
		parts = append(parts, fmt.Sprintf("-- Based on %s, %s", configFName, now.Format("2006-01-02")))
	}
	if len(parts) > 2 {
		parts = append(parts, "--")
	}
	parts = append(parts, "")
	return []byte(strings.Join(parts, "\n"))
}

// dslToBaseType will parse a DataType string and return the
// base type (e.g. String) and a boolean indiciting if it is the
// primary key.
func dslToBaseType(s string) (string, bool) {
	parts := strings.Split(strings.TrimSpace(s), " ")
	if len(parts) == 0 {
		return "", false
	}
	tName := strings.TrimSuffix(parts[0], "*")
	return tName, (parts[0] != tName) 
}

// dslToSQLType converts a string describing a type into a SQL
// type definition.
func dslToSQLType(s string) (string, error) {
    parts := strings.Split(strings.TrimSpace(s), " ")
    if len(parts) == 0 {
    	return "", fmt.Errorf("missing type definition")
    }
    // NOTE: I'm trying out a suffix of asterix for indicating that the
    // type is going to be used as a primary key in a SQL table.
    tName := strings.TrimSuffix(parts[0], "*")
    isPrimaryKey := (tName != parts[0])
    if sqlType, ok := SQLTypes[tName]; ok {
       if isPrimaryKey {
		   if strings.HasPrefix(sqlType, "INTEGER ") {
			   return fmt.Sprintf("SERIAL PRIMARY KEY"), nil
		   }
       	   return fmt.Sprintf("%s PRIMARY KEY", sqlType), nil
       } else {
       	   return sqlType, nil
       }
    }
    return "", fmt.Errorf("Cannot map type %q to SQL", parts[0])
}

// createStatement generates a SQL CREATE statement from a model.
func createStatement(model *ModelDSL) ([]byte, error) {
   	parts := []string{}
	prefix := fmt.Sprintf("CREATE TABLE %s (\n    ", model.Name)
    suffix := "\n);\n\n"
	for k, v := range model.Var {
	    t, err := dslToSQLType(v)
	    if err != nil {
	    	return nil, fmt.Errorf("erorr in defining %q, %s", k, err)
	    }
		parts = append(parts, fmt.Sprintf("%s %s", k, t))
	}
	return []byte(fmt.Sprintf("%s%s%s", prefix, strings.Join(parts, ",\n    "), suffix)), nil
}


// createListView generates a SQL view statement from a model.
func createListView(model *ModelDSL) ([]byte, error) {
   	parts := []string{}
	prefix := fmt.Sprintf(`--
-- LIST VIEW: %s 
-- FIXME: You probably want to customized this statement 
-- (e.g. add WHERE CLAUSE, ORDER BY, GROUP BY).
--
CREATE OR REPLACE VIEW %s_list_view AS
    SELECT `, model.Name, model.Name)
    suffix := fmt.Sprintf("\n    FROM %s;\n\n", model.Name)
    //FIXME: need to code up table attributes.
	for k, _ := range model.Var {
		parts = append(parts, k)
	}
	return []byte(fmt.Sprintf("%s%s%s", prefix, strings.Join(parts, ", "), suffix)), nil
}

// PgSetupSQL generate PostgREST setup SQL for models given a configuration name and model.
func PgSetupSQL(configFName string, model *ModelDSL) ([]byte, error) {
	if model == nil {
		return nil, fmt.Errorf("No model found in %s", configFName)
	}
	return nil, fmt.Errorf("PgSetupSQL() not implemented")
}

// PgModelSQL generate Postgres SQL from a model given a configuration name
// and model. 
func PgModelSQL(configFName string, model *ModelDSL) ([]byte, error) {
	if model == nil {
		return nil, fmt.Errorf("No model found in %s", configFName)
	}
    src := []byte{}
    src = append(src, commentSection(configFName, model.Name)...)
    s, err := createStatement(model)
    if err != nil {
    	return nil, err
    }
    src = append(src, s...)
	s, err = createListView(model)
	if err != nil {
		return nil, err
	}
    src = append(src, s...)
	//FIXME: Need to figure out a function/procedure to return a single
	// record (row) given a primary key value.
	return src, nil
}

// PgModelTestSQL generate Postgres SQL from a model given a 
// configuration name and model. 
func PgModelTestSQL(configFName string, model *ModelDSL) ([]byte, error) {
	if model == nil {
		return nil, fmt.Errorf("No model found in %s", configFName)
	}
    src := []byte{}
    src = append(src, commentSection(configFName, model.Name)...)
    s, err := createStatement(model)
    if err != nil {
    	return nil, err
    }
    src = append(src, s...)
	s, err = createListView(model)
	if err != nil {
		return nil, err
	}
    src = append(src, s...)
	//FIXME: Need to figure out a function/procedure to return a single
	// record (row) given a primary key value.
	return src, nil
}


