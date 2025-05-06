package gotests

type StackCancelFunc[T cancelFuncType] []T

type cancelFuncType interface {
	CancelFunc | CancelFuncWithErr
}

type CancelFunc func()

type CancelFuncWithErr func() error

func (s *StackCancelFunc[T]) Push(fn T) {
	*s = append(*s, fn)
}

func (s *StackCancelFunc[T]) Clear() error {
	for i := len(*s) - 1; i >= 0; i-- {
		switch fn := any((*s)[i]).(type) {
		case CancelFunc:
			fn()
		case CancelFuncWithErr:
			if err := fn(); err != nil {
				return err
			}
		}
	}

	return nil
}
