package week02

import (
	"github.com/pkg/errors"
)

func dao() (int, error) {
	result, err := sql.Select()
	if err != nil {
		// oops, something went wrong in basic package
		// wrap it!
		if x, ok := err.(interface{ Is(error) bool }); ok && x.Is(sql.ErrNoRows) {
			// extra work for specific error type
			// do stuff

			return 0, errors.Wrap(err, "this is a sql.ErrNoRows")
		}
		return 0, errors.Wrap(err, "something wrong with sql")
	}
	// do stuff

	return result, nil
}

func service() (int, error) {
	value, err := dao()
	if err != nil {
		// inside the same package, do not Wrap error
		return 0, err
		// to add more context
		// return errors.WithMessage(err, "the sql is %s", sqlStr)
	}
	// do stuff

	return value, nil
}

func api() (int, error) {
	data, err := service()
	if err != nil {
		// only handle error once
		logs.Error(err)
		// downgrading, maintain the integrity of data
		data = 0xdead
		return data, nil
	}
	// do stuff

	return data, nil
}
