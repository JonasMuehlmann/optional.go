package optional_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/JonasMuehlmann/optional.go"
	"github.com/gocarina/gocsv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// GetDB opens a copy of the test DB in memory.
func GetDB(t *testing.T) (*sql.DB, error) {
	// Connect to new temporary database
	db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=1")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestUnMarshallFromWhole(t *testing.T) {
	marshaled := []byte(`{"wrapee": 123, "has_value": true}`)

	var optional optional.Optional[int]

	err := json.Unmarshal(marshaled, &optional)
	assert.NoError(t, err)
	assert.Equal(t, 123, optional.Wrappee)
	assert.Equal(t, true, optional.HasValue)
}

func TestUnMarshallFromWholeNoValue(t *testing.T) {
	marshaled := []byte(`{"wrapee": 123, "has_value": false}`)

	var optional optional.Optional[int]

	err := json.Unmarshal(marshaled, &optional)
	assert.NoError(t, err)
	assert.Equal(t, 123, optional.Wrappee)
	assert.Equal(t, false, optional.HasValue)
}

func TestUnMarshallFromWrappee(t *testing.T) {
	targetStruct := struct {
		Foo        string                 `json:"foo"`
		MyOptional optional.Optional[int] `json:"my_optional"`
	}{}
	marshaled, err := json.Marshal(map[string]any{"foo": "bar", "my_optional": 123})
	assert.NoError(t, err)

	err = json.Unmarshal(marshaled, &targetStruct)
	assert.NoError(t, err)
	assert.Equal(t, 123, targetStruct.MyOptional.Wrappee)
	assert.Equal(t, true, targetStruct.MyOptional.HasValue)
}

// ******************************************************************//
//                           fmt.Stringer                           //
// ******************************************************************//.
func TestOptionalToStringHasValue(t *testing.T) {
	optional := optional.Optional[int]{Wrappee: 123, HasValue: true}

	stringified := fmt.Sprint(optional)
	assert.Equal(t, "123", stringified)
}

func TestOptionalToStringHasNoValue(t *testing.T) {
	optional := optional.Optional[int]{Wrappee: 123, HasValue: false}

	stringified := fmt.Sprint(optional)
	assert.Equal(t, "empty optional", stringified)
}

// ******************************************************************//
//        database.sql.Scanner and database.sql.driver.Valuer       //
// ******************************************************************//.
func TestSQLScanAndValueHasValue(t *testing.T) {
	db, err := GetDB(t)
	assert.NoError(t, err)

	_, err = db.Exec("CREATE TABLE FOO(BAR TEXT)")
	assert.NoError(t, err)

	input := optional.Optional[string]{Wrappee: "bar", HasValue: true}
	_, err = db.Exec("INSERT INTO FOO VALUES(?)", input)
	assert.NoError(t, err)

	var output optional.Optional[string]
	err = db.QueryRow("SELECT * FROM Foo").Scan(&output)
	assert.NoError(t, err)
	assert.Equal(t, input, output)
}

func TestSQLScanAndValueHasNoValue(t *testing.T) {
	db, err := GetDB(t)
	assert.NoError(t, err)

	_, err = db.Exec("CREATE TABLE FOO(BAR TEXT)")
	assert.NoError(t, err)

	input := optional.Optional[string]{Wrappee: "bar", HasValue: false}
	_, err = db.Exec("INSERT INTO FOO VALUES(?)", input)
	assert.NoError(t, err)

	var output optional.Optional[string]
	err = db.QueryRow("SELECT * FROM Foo").Scan(&output)
	assert.NoError(t, err)
	// If HasValue is false, we don't care about the Wrappee
	assert.Equal(t, input.HasValue, output.HasValue)
}

//******************************************************************//
// gocarina/gocsv.TypeUnmarshaller and gocarina/gocsv.TypeMarshaller//
//******************************************************************//
func TestCSVUnmarshalHasValue(t *testing.T) {
	inputCSV := `Foo,Bar
1,2
`
	outputStructs := []struct {
		Foo optional.Optional[int] `csv:"Foo"`
		Bar optional.Optional[int] `csv:"Bar"`
	}{}

	err := gocsv.UnmarshalString(inputCSV, &outputStructs)
	assert.NoError(t, err)
	assert.Len(t, outputStructs, 1)

	assert.True(t, outputStructs[0].Foo.HasValue)
	assert.True(t, outputStructs[0].Bar.HasValue)

	assert.Equal(t, 1, outputStructs[0].Foo.Wrappee)
	assert.Equal(t, 2, outputStructs[0].Bar.Wrappee)
}

func TestCSVUnmarshalHasNoValue(t *testing.T) {
	inputCSV := `Foo,Bar
,
`
	outputStructs := []struct {
		Foo optional.Optional[int] `csv:"Foo"`
		Bar optional.Optional[int] `csv:"Bar"`
	}{}

	err := gocsv.UnmarshalString(inputCSV, &outputStructs)
	assert.NoError(t, err)
	assert.Len(t, outputStructs, 1)

	assert.False(t, outputStructs[0].Foo.HasValue)
	assert.False(t, outputStructs[0].Bar.HasValue)
}

func TestCSVMarshalHasValue(t *testing.T) {
	outputCSV := `Foo,Bar
1,2
`
	inputStruct := []struct {
		Foo optional.Optional[int] `csv:"Foo"`
		Bar optional.Optional[int] `csv:"Bar"`
	}{
		{
			Foo: optional.Make(1),
			Bar: optional.Make(2),
		},
	}

	out, err := gocsv.MarshalString(inputStruct)
	assert.NoError(t, err)
	assert.Equal(t, outputCSV, out)
}

func TestCSVMarshalHasNoValue(t *testing.T) {
	outputCSV := `Foo,Bar
,
`
	inputStruct := []struct {
		Foo optional.Optional[int] `csv:"Foo"`
		Bar optional.Optional[int] `csv:"Bar"`
	}{
		{
			Foo: optional.Optional[int]{},
			Bar: optional.Optional[int]{},
		},
	}

	out, err := gocsv.MarshalString(inputStruct)
	assert.NoError(t, err)
	assert.Equal(t, outputCSV, out)
}
