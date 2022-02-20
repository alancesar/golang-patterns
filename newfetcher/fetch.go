package newfetcher

import "sync"

type (
	ProviderFn[T any]                func() (T, error)
	ProviderWithParamFn[T, K any]    func(T) (K, error)
	CallbackFn[T sync.Locker, S any] func(target T, source S) error
)

func New[T sync.Locker, K any](target T, provider ProviderFn[K], callback CallbackFn[T, K]) error {
	defer target.Unlock()
	source, err := provider()
	if err != nil {
		return err
	}
	target.Lock()
	return callback(target, source)
}

func NewWithParam[T sync.Locker, P any, S any](target T, param P, provider ProviderWithParamFn[P, S], callback CallbackFn[T, S]) error {
	return New(target, func() (S, error) {
		return provider(param)
	}, callback)
}
