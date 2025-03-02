package autowire

import (
	"context"
	"reflect"
)

// Build implementation of Container interface
func (c *container) Build(targetType reflect.Type, opts ...ContextOption) (value reflect.Value, err error) {
	provider, err := c.providerSet.GetFor(targetType)
	if err != nil {
		return value, err
	}

	ctx := &Context{
		sharedMode:     c.sharedMode,
		providerSet:    c.providerSet.shallowClone(),
		objectMap:      c.objectMap,
		resolvingTypes: make(map[reflect.Type]struct{}, 10), //nolint:gomnd
	}
	for _, opt := range opts {
		opt(ctx)
	}

	value, err = provider.Build(ctx, targetType)
	if err != nil {
		return value, err
	}

	return value, nil
}

// BuildWithCtx implementation of Container interface
func (c *container) BuildWithCtx(ctx context.Context, targetType reflect.Type, opts ...ContextOption) (
	value reflect.Value, err error,
) {
	return c.Build(targetType, append(opts, ProviderOverwrite(ctx))...)
}

// BuildAll creates a values and all other required values.
func (c *container) BuildAll(opts ...ContextOption) error {
	for _, provider := range c.providerSet.GetAll() {
		ctx := &Context{
			sharedMode:     c.sharedMode,
			providerSet:    c.providerSet.shallowClone(),
			objectMap:      c.objectMap,
			resolvingTypes: make(map[reflect.Type]struct{}, 10), //nolint:gomnd
		}
		for _, opt := range opts {
			opt(ctx)
		}

		targetType := provider.TargetTypes()[0]

		_, err := provider.Build(ctx, targetType)
		if err != nil {
			return err
		}
	}
	return nil
}
