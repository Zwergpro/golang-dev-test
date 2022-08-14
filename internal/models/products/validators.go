package products

import (
	"errors"
)

func ValidateName(name string) error {
	if len(name) == 0 {
		return errors.New("name length must be greater than 0")
	}
	return nil
}

func ValidatePrice(price uint64) error {
	if price == 0 {
		return errors.New("price must be greater than 0")
	}
	return nil
}

func ValidateQuantity(quantity uint64) error {
	if quantity == 0 {
		return errors.New("quantity must be greater than 0")
	}
	return nil
}

func ValidateProductFields(name string, price, quantity uint64) []error {
	validationErrors := make([]error, 0, 3)

	if err := ValidateName(name); err != nil {
		validationErrors = append(validationErrors, err)
	}

	if err := ValidatePrice(price); err != nil {
		validationErrors = append(validationErrors, err)
	}

	if err := ValidateQuantity(quantity); err != nil {
		validationErrors = append(validationErrors, err)
	}

	return validationErrors
}
