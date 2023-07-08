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
func dslToSQLType(s string, useSimpleType bool) (string, error) {
	parts := strings.Split(strings.TrimSpace(s), " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("missing type definition")
	}
	// NOTE: I'm trying out a suffix of asterix for indicating that the
	// type is going to be used as a primary key in a SQL table.
	tName := strings.TrimSuffix(parts[0], "*")
	isPrimaryKey := (tName != parts[0])
	if sqlType, ok := SQLTypes[tName]; ok {
		if useSimpleType {
			parts := strings.SplitN(sqlType, " ", 2)
			if len(parts) > 0 {
				sqlType = parts[0]
			}
		}
		if isPrimaryKey && ! useSimpleType {
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

// getNamespaceFlatName takes a ModelDSL.Name and returns
// a namespace and flattened name strings
func getNamespaceFlatName(modelName string) (string, string) {
	if strings.Contains(modelName, ".") {
		parts := strings.SplitN(modelName, ".", 2)
		if len(parts) == 2 {
			return parts[0], parts[1]
		}
	}
	return modelName, modelName
}
// createStatement generates a SQL CREATE statement from a model.
func createStatement(model *ModelDSL) ([]byte, error) {
	parts := []string{}
	prefix := fmt.Sprintf(`
DROP TABLE IF EXISTS %s CASCADE;
CREATE TABLE %s (
    `, model.Name, model.Name)
	suffix := "\n);\n\n"
	for k, v := range model.Var {
		t, err := dslToSQLType(v, false)
		if err != nil {
			return nil, fmt.Errorf("erorr in defining %q, %s", k, err)
		}
		parts = append(parts, fmt.Sprintf("%s %s", k, t))
	}
	return []byte(fmt.Sprintf("%s%s%s", prefix, strings.Join(parts, ",\n    "), suffix)), nil
}

// createListView generates a SQL view statement from a model.
func createListView(model *ModelDSL) ([]byte, error) {
	namespace, flatName := getNamespaceFlatName(model.Name)
	parts := []string{}
	prefix := fmt.Sprintf(`--
-- LIST VIEW: %s 
-- FIXME: You probably want to customized this statement 
-- (e.g. add WHERE CLAUSE, ORDER BY, GROUP BY).
--
CREATE OR REPLACE VIEW %s.%s_view AS
    SELECT `, model.Name, namespace, flatName)
	suffix := fmt.Sprintf("\n    FROM %s;\n\n", model.Name)
	//FIXME: need to code up table attributes.
	for k, _ := range model.Var {
		parts = append(parts, k)
	}
	return []byte(fmt.Sprintf("%s%s%s", prefix, strings.Join(parts, ", "), suffix)), nil
}

// createGetRecord generates a SQL function for retrieving record
// by primary key.
func createGetRecord(model *ModelDSL) ([]byte, error) {
	namespace, flatName := getNamespaceFlatName(model.Name)
	fnName := "get_" + flatName
	sqlColDefs := []string{}
	colNames := []string{}
	primaryKey, primaryKeyDef := "", ""
	for varName, tVal := range model.Var {
		sqlType, err := dslToSQLType(tVal, true)
		if err != nil {
			return nil, fmt.Errorf("-- Failed to generate function %q for model %q, %s", fnName, model.Name, err)
		}
		if strings.HasSuffix(tVal, "*") {
			primaryKey = varName
			primaryKeyDef = sqlType
		}
		sqlColDefs = append(sqlColDefs, fmt.Sprintf("%s %s", varName, sqlType))
		colNames = append(colNames, varName)
	}
	if primaryKey == "" {
		return nil, fmt.Errorf("-- Failed to find primary key to generate function %q for model %q", fnName, model.Name)
	}

	txt := fmt.Sprintf(`--
-- {func_name} provides a 'get record' for model %s
--
DROP FUNCTION IF EXISTS {namespace}.{func_name}({primary_key} {primary_key_def});
CREATE FUNCTION {namespace}.{func_name}({primary_key} {primary_key_def})
RETURNS TABLE ({sql_col_defs}) AS $$
	SELECT {col_names} FROM {model_name} WHERE {primary_key} = {primary_key} LIMIT 1
$$ LANGUAGE SQL;
`, model.Name)
	src := []byte(
		strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
			txt, "{func_name}", fnName),
			"{namespace}", namespace),
			"{primary_key_def}", primaryKeyDef),
			"{sql_col_defs}", strings.Join(sqlColDefs, ", ")),
			"{col_names}", strings.Join(colNames, ", ")),
			"{model_name}", model.Name),
			"{primary_key}", primaryKey))
	return []byte(src), nil
}

// createAddRecord generates a SQL function for create record
// return created record with primary key.
func createAddRecord(model *ModelDSL) ([]byte, error) {
	namespace, flatName := getNamespaceFlatName(model.Name)
	fnName := "add_" + flatName
	sqlColDefs := []string{}
	colNames := []string{}
	primaryKey, primaryKeyDef := "", ""
	for varName, tVal := range model.Var {
		sqlType, err := dslToSQLType(tVal, true)
		if err != nil {
			return nil, fmt.Errorf("-- Failed to generate function %q for model %q, %s", fnName, model.Name, err)
		}
		if strings.HasSuffix(tVal, "*") {
			primaryKey = varName
			primaryKeyDef = sqlType
		}
		sqlColDefs = append(sqlColDefs, fmt.Sprintf("%s %s", varName, sqlType))
		colNames = append(colNames, varName)
	}
	if primaryKey == "" {
		return nil, fmt.Errorf("-- Failed to find primary key to generate function %q for model %q", fnName, model.Name)
	}

	txt := fmt.Sprintf(`--
-- {func_name} provides an 'add record' for model %s
-- It returns the row inserted including the primary key
DROP FUNCTION IF EXISTS {namespace}.{func_name}({sql_col_defs});
CREATE FUNCTION {namespace}.{func_name}({sql_col_defs})
RETURNS {primary_key_def} AS $$
    INSERT INTO {model_name} 
               ({col_names}) 
        VALUES ({col_names})
    RETURNING {primary_key}
$$ LANGUAGE SQL;
`, model.Name)
	src := []byte(
		strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(
			txt, "{func_name}", fnName),
			"{namespace}", namespace),
			"{primary_key_def}", primaryKeyDef),
			"{sql_col_defs}", strings.Join(sqlColDefs, ", ")),
			"{col_names}", strings.Join(colNames, ", ")),
			"{model_name}", model.Name),
			"{primary_key}", primaryKey))
	return []byte(src), nil
}

// createUpdateRecord generates a SQL function for updating record
// returning the updated row.
func createUpdateRecord(model *ModelDSL) ([]byte, error) {
	namespace, flatName := getNamespaceFlatName(model.Name)
	txt := fmt.Sprintf(`--
-- FIXME: Need to generate a 'update record' for model %s
-- FIXME: namespace %q, flat name %q

`, model.Name, namespace, flatName)

	return []byte(txt), nil
}

// PgSetupSQL generate Postgres+PostgREST setup SQL for roles in a
// given namespace
func PgSetupSQL(configFName string, namespace string, password string) ([]byte, error) {
	now := time.Now()
	txt := fmt.Sprintf(`--
-- Setup new empty database and schema for {namespace} based on {configFName}, %s
--
DROP DATABASE IF EXISTS {namespace};
CREATE DATABASE {namespace};
\c {namespace}
DROP SCHEMA IF EXISTS {namespace} CASCADE;
CREATE SCHEMA {namespace};

--
-- Create role "{namespace}_anonymous"
--
DROP ROLE IF EXISTS {namespace}_anonymous;
CREATE ROLE {namespace}_anonymous nologin;

--
-- Create role "{namespace}_authenticated"
--
DROP ROLE IF EXISTS {namespace}_authenticated;
-- WARNING: This "CREATE ROLE" statement sets a password!!!!
-- Don't make this publically available!!!!
CREATE ROLE {namespace}_authenticated NOINHERIT LOGIN PASSWORD '{password_goes_here}';

`, now.Format("2006-01-02"))
	if namespace == "" {
		return nil, fmt.Errorf("No namespace found in %s", configFName)
	}
	if password == "" {
		password = "<PASSWORD_GOES_HERE>"
	}
	return []byte(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(txt, "{configFName}", configFName), "{namespace}", namespace), "{password_goes_here}", password)), nil
}

// PgModelSQL generate Postgres SQL from a model given a configuration name
// and model.
func PgModelSQL(configFName string, model *ModelDSL) ([]byte, error) {
	if model == nil {
		return nil, fmt.Errorf("No model found in %s", configFName)
	}
	namespace := model.Name
	if strings.Contains(namespace, ".") {
		parts := strings.SplitN(namespace, ".", 2)
		if len(parts) > 0 {
			namespace = parts[0]
		}
	}
	now := time.Now()
	src := []byte(fmt.Sprintf(`--
-- Define Models for %s in %s, %s
--
\c %s
SET search_path TO %s, public;

`, model.Name, configFName, now.Format("2006-01-02"), namespace, namespace))
	src = append(src, commentSection(configFName, model.Name)...)
	s, err := createStatement(model)
	if err != nil {
		return nil, err
	}
	src = append(src, s...)

	// Create list view, PostgREST path `/{model_name}_view`
	s, err = createListView(model)
	if err != nil {
		return nil, err
	}
	src = append(src, s...)

	// Create a get record function, `PostgREST path `/rpc/get_{model_name}`
	s, err = createGetRecord(model)
	if err != nil {
		return nil, err
	}
	src = append(src, s...)

	// Create an add record function. PostgREST path `/rpc/add_{model_name}`
	s, err = createAddRecord(model)
	if err != nil {
		return nil, err
	}
	src = append(src, s...)

/*
	// Create an update record function. PostgREST path `/rpc/update_{model_name}`
	s, err = createUpdateRecord(model)
	if err != nil {
		return nil, err
	}
	src = append(src, s...)
*/
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

// PgModelPermissions generates a SQL GRANT statements for
// a Schema's anonymous and authenticated roles.
func PgModelPermissions(configFName, namespace string, modelNames []string) ([]byte, error) {
	anonRole := namespace + "_anonymous"
	authRole := namespace + "_authenticated"
	src := []byte(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(`--
-- PostgREST access and controls.
--
-- GRANT or REVOKE role permissions here to match our models.
--

-- Give access to the Schema to PostgREST for each role.
GRANT USAGE ON SCHEMA {namespace} TO {anon_role};
GRANT USAGE ON SCHEMA {namespace} TO {auth_role};

`, "{namespace}", namespace), "{anon_role}", anonRole), "{auth_role}", authRole))

	txt := `--
-- Permissions for model {model_name}
--

-- Access for our anonymous role {anon_role}
GRANT SELECT ON {model_name} TO {anon_role};
GRANT SELECT ON {namespace}.{flat_name}_view TO {anon_role};
GRANT EXECUTE ON FUNCTION {namespace}.get_{flat_name} TO {anon_role};

-- Access for our authenticated role {auth_role}
GRANT SELECT, INSERT, UPDATE, DELETE ON {model_name} TO {auth_role};
GRANT SELECT ON {namespace}.{flat_name}_view TO {auth_role};
GRANT EXECUTE ON FUNCTION {namespace}.get_{flat_name} TO {auth_role};
-- GRANT EXECUTE ON FUNCTION {namespace}.add_{flat_name} TO {auth_role};
-- GRANT EXECUTE ON FUNCTION {namespace}.update_{flat_name} TO {auth_role};

`
	for _, modelName := range modelNames {
		_, flatName := getNamespaceFlatName(modelName)
		src = append(src, []byte(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(txt, "{namespace}", namespace), "{anon_role}", anonRole), "{model_name}", modelName), "{auth_role}", authRole), "{flat_name}", flatName))...)
	}
	return src, nil
}
