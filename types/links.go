package types

import (
	"google.golang.org/api/firebasedynamiclinks/v1"
)

type LinkInfo struct {
	// Link represents the link your app will open, You can specify any URL your app can handle.
	// This link must be a well-formatted URL, be properly URL-encoded, and use the HTTP or
	// HTTPS scheme. See 'link' parameters in the documentation
	// (https://firebase.google.com/docs/dynamic-links/create-manually).
	Link string

	// Suffix represents the Short Dynamic Link suffix
	Suffix *firebasedynamiclinks.Suffix

	// SocialMetaTagInfo contains the parameters for social meta tag params. Used to set
	// meta tag data for link previews on social sites.
	SocialMetaTagInfo *firebasedynamiclinks.SocialMetaTagInfo
}

// ShortLinkSuffix returns a new Suffix instance that can be used to create a short link
func ShortLinkSuffix() *firebasedynamiclinks.Suffix {
	return &firebasedynamiclinks.Suffix{
		Option: "SHORT",
	}
}
