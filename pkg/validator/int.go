package validator

type intValidator struct {
	NotZeroError error
}

func NewIntValidator() intValidator {
	return intValidator{
		NotZeroError: ErrInvalidValue,
	}
}

func ValidateInt(key string, value int, validators ...func(int) error) error {
	errors := make(Errors)
	for _, validator := range validators {
		if err := validator(value); err != nil {
			errors[key] = err
			return errors
		}
	}

	return nil
}

func (cfg intValidator) IntNotZero(i int) error {
	if i == 0 {
		return cfg.NotZeroError
	}

	return nil
}
