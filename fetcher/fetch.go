package fetcher

import "sync"

type (
	ProviderFn          func() (interface{}, error)
	ProviderWithParamFn func(interface{}) (interface{}, error)
	CallbackFn          func(target sync.Locker, source interface{}) error
)

func New(target sync.Locker, provider ProviderFn, callback CallbackFn) error {
	defer target.Unlock()
	source, err := provider()
	if err != nil {
		return err
	}
	target.Lock()
	return callback(target, source)
}

func NewWithParam(target sync.Locker, param interface{}, provider ProviderWithParamFn, callback CallbackFn) error {
	return New(target, func() (interface{}, error) {
		return provider(param)
	}, callback)
}
