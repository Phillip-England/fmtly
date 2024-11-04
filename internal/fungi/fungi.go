package fungi

func ProcessErrFuncs(funcs ...func() error) error {
	for _, fn := range funcs {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func Process(funcs ...func() error) error {
	for _, fn := range funcs {
		if err := fn(); err != nil {
			return err // Return the first error encountered
		}
	}
	return nil // No errors found, return nil
}
