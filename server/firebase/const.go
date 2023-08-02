//nolint:gosec

package firebase

const (
	EnvCredentialsFilePath = "FIREBASE_CREDENTIALS_FILE_PATH"

	EnvNotificationsSendNotificationField = "FIREBASE_NOTIFICATIONS_SEND_NOTIFICATION_FIELD"

	EnvLinksDomain                       = "FIREBASE_LINKS_DOMAIN"
	EnvLinksDesktopFallbackLink          = "FIREBASE_LINKS_DESKTOP_FALLBACK_LINK"
	EnvLinksAndroidPackageName           = "FIREBASE_LINKS_ANDROID_PACKAGE_NAME"
	EnvLinksAndroidMinPackageVersionCode = "FIREBASE_LINKS_ANDROID_MIN_PACKAGE_VERSION_CODE"
	EnvLinksIOSBundleID                  = "FIREBASE_LINKS_IOS_BUNDLE_ID"
	EnvLinksIOSMinimumVersion            = "FIREBASE_LINKS_IOS_MINIMUM_VERSION"
	EnvLinksIOSAppStoreID                = "FIREBASE_LINKS_IOS_APP_STORE_ID"
)

const (
	NotificationTitleKey = "notification_title"
	NotificationBodyKey  = "notification_body"
	NotificationImageKey = "notification_image"
)
