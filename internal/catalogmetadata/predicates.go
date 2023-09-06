package catalogmetadata

import (
	mmsemver "github.com/Masterminds/semver/v3"
	bsemver "github.com/blang/semver/v4"

	olmentity "github.com/operator-framework/operator-controller/internal/resolution/entities"
)

// TODO: Move somewhere nice. Probably into a ./predicates package

func WithPackageName(packageName string) Predicate[Bundle] {
	return func(bundle *Bundle) bool {
		return bundle.Package == packageName
	}
}

func InMastermindsSemverRange(semverRange *mmsemver.Constraints) Predicate[Bundle] {
	return func(bundle *Bundle) bool {
		bVersion, err := bundle.Version()
		if err != nil {
			return false
		}
		// No error should occur here because the simple version was successfully parsed by blang
		// We are unaware of any tests cases that would cause one to fail but not the other
		// This will cause code coverage to drop for this line. We don't ignore the error because
		// there might be that one extreme edge case that might cause one to fail but not the other
		mVersion, err := mmsemver.NewVersion(bVersion.String())
		if err != nil {
			return false
		}
		return semverRange.Check(mVersion)
	}
}

func InBlangSemverRange(semverRange bsemver.Range) Predicate[Bundle] {
	return func(bundle *Bundle) bool {
		bundleVersion, err := bundle.Version()
		if err != nil {
			return false
		}
		return semverRange(*bundleVersion)
	}
}

func InChannel(channelName string) Predicate[Bundle] {
	return func(bundle *Bundle) bool {
		for _, ch := range bundle.InChannels {
			if ch.Name == channelName {
				return true
			}
		}
		return false
	}
}

func ProvidesGVK(gvk *olmentity.GVK) Predicate[Bundle] {
	return func(bundle *Bundle) bool {
		providedGVKs, err := bundle.ProvidedGVKs()
		if err != nil {
			return false
		}
		for i := 0; i < len(providedGVKs); i++ {
			providedGVK := &providedGVKs[i]
			if providedGVK.String() == gvk.String() {
				return true
			}
		}
		return false
	}
}

func WithBundleImage(bundleImage string) Predicate[Bundle] {
	return func(bundle *Bundle) bool {
		return bundle.Image == bundleImage
	}
}

func Replaces(bundleName string) Predicate[Bundle] {
	return func(bundle *Bundle) bool {
		for _, ch := range bundle.InChannels {
			for _, chEntry := range ch.Entries {
				if bundle.Name == chEntry.Name && chEntry.Replaces == bundleName {
					return true
				}
			}
		}
		return false
	}
}
