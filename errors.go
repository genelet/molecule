package molecule

import (
	"fmt"
)

func errorActionNotDefined(name string) error {
	return fmt.Errorf("action %s not defined", name)
}

func errorActionNil(name string) error {
	return fmt.Errorf("actions or action %s is nil", name)
}

func errorInputDataType(v interface{}) error {
	return fmt.Errorf("wrong input data type %T", v)
}

func errorExtraDataType(v interface{}) error {
	return fmt.Errorf("wrong extra data type %T", v)
}

func errorAtomNotFound(name string) error {
	return fmt.Errorf("atom %s not found in molecule", name)
}

func errorActionNotFound(action, atom string) error {
	return fmt.Errorf("action %s not found in atom %s", action, atom)
}

func errorRollback(err, roolbackErr error) error {
	return fmt.Errorf("error original: %v, rollback: %v", err, roolbackErr)
}

func errorMissingKeys(name string) error {
	return fmt.Errorf("no pk nor fk is found in table %s", name)
}

func errorMissingPk(name string) error {
	return fmt.Errorf("no pk is found in table %s", name)
}

func errorDeleteWhole(name string) error {
	return fmt.Errorf("delete whole table %s not allowed", name)
}

func errorEmptyInput(name string) error {
	return fmt.Errorf("no input data to %s", name)
}

func errorNoSuchColumn(name string) error {
	return fmt.Errorf("column '%s' not found in input", name)
}

func errorNoUniqueKey(name string) error {
	return fmt.Errorf("unique key not defined in %s", name)
}

func errorNotUnique(name string) error {
	return fmt.Errorf("multiple records found for unique key in %s", name)
}
