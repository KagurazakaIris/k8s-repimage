package utils

import (
	"strings"
)

const (
	defaultDomain    = "docker.io"
	officialRepoName = "library"
)

// ReplaceImageName adds a mirror prefix or mapping to container image names, but skips domains in ignoreDomains.
// mappings maps original domains to replacement hostnames. If mappings is empty, falls back to the legacy prefix behavior.
func ReplaceImageName(prefix string, ignoreDomains []string, mappings map[string]string, name string) string {
	parts := strings.SplitN(name, "/", 3)
	// if the image already uses the prefix, leave as-is
	if parts[0] == prefix {
		return name
	}

	// helper to build with a replacement host
	buildWithHost := func(host string, rest ...string) string {
		all := append([]string{host}, rest...)
		return strings.Join(all, "/")
	}

	// If mapping exists for a domain, use it (unless domain is in ignore list)
	switch len(parts) {
	case 1:
		// image like "nginx"
		if mapped, ok := mappings[defaultDomain]; ok && !shouldIgnoreDomain(defaultDomain, ignoreDomains) {
			return buildWithHost(mapped, officialRepoName, parts[0])
		}

		if shouldIgnoreDomain(defaultDomain, ignoreDomains) {
			return buildWithHost(defaultDomain, officialRepoName, parts[0])
		}
		return buildWithHost(prefix, defaultDomain, officialRepoName, parts[0])
	case 2:
		// either user/image or domain/image
		if !isDomain(parts[0]) {
			if mapped, ok := mappings[defaultDomain]; ok && !shouldIgnoreDomain(defaultDomain, ignoreDomains) {
				return buildWithHost(mapped, parts[0], parts[1])
			}

			if shouldIgnoreDomain(defaultDomain, ignoreDomains) {
				return buildWithHost(defaultDomain, parts[0], parts[1])
			}

			return buildWithHost(prefix, defaultDomain, parts[0], parts[1])
		}

		// parts[0] is a domain
		domain := parts[0]
		if isLegacyDefaultDomain(domain) {
			domain = defaultDomain
		}

		if shouldIgnoreDomain(domain, ignoreDomains) {
			return buildWithHost(domain, parts[1])
		}

		if mapped, ok := mappings[domain]; ok {
			return buildWithHost(mapped, parts[1])
		}

		return buildWithHost(prefix, parts[0], parts[1])
	case 3:
		// domain/org/repo or potentially non-domain/org/repo
		if !isDomain(parts[0]) {
			if mapped, ok := mappings[defaultDomain]; ok && !shouldIgnoreDomain(defaultDomain, ignoreDomains) {
				return buildWithHost(mapped, parts[0], parts[1], parts[2])
			}

			if shouldIgnoreDomain(defaultDomain, ignoreDomains) {
				return buildWithHost(defaultDomain, parts[0], parts[1], parts[2])
			}

			return buildWithHost(prefix, defaultDomain, parts[0], parts[1], parts[2])
		}

		domain := parts[0]
		if isLegacyDefaultDomain(domain) {
			domain = defaultDomain
		}

		if shouldIgnoreDomain(domain, ignoreDomains) {
			return buildWithHost(domain, parts[1], parts[2])
		}

		if mapped, ok := mappings[domain]; ok {
			return buildWithHost(mapped, parts[1], parts[2])
		}

		return buildWithHost(prefix, parts[0], parts[1], parts[2])
	}
	return name
}

// shouldIgnoreDomain checks if the image domain should be ignored
func shouldIgnoreDomain(domain string, ignoreDomains []string) bool {
	for _, ignoreDomain := range ignoreDomains {
		if domain == ignoreDomain {
			return true
		}
	}
	return false
}

func isDomain(name string) bool {
	return strings.Contains(name, ".")
}

var (
	legacyDefaultDomain = map[string]struct{}{
		"index.docker.io":      {},
		"registry-1.docker.io": {},
	}
)

func isLegacyDefaultDomain(name string) bool {
	_, ok := legacyDefaultDomain[name]
	return ok
}
