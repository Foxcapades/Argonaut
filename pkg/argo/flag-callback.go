package argo

type FlagCallback = func(flag Flag)

func SimpleFlagCallback(fn func()) FlagCallback {
	return func(Flag) { fn() }
}
