package argo

type FlagCallback = func(flag FlagSpec)

func SimpleFlagCallback(fn func()) FlagCallback {
	return func(FlagSpec) { fn() }
}
