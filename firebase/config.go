package firebase

import (
	"strconv"

	"github.com/desmos-labs/caerus/utils"
)

type Config struct {
	CredentialsFilePath string               `json:"credentials_file_path" yaml:"credentials_file_path" toml:"credentials_file_path"`
	Notifications       *NotificationsConfig `json:"notifications" yaml:"notifications" toml:"notifications"`
	Links               *LinksConfig         `json:"links" yaml:"links" toml:"links"`
}

func ReadConfigFromEnvVariables() (*Config, error) {
	notificationsConfig, err := ReadNotificationsConfigFromEnvVariables()
	if err != nil {
		return nil, err
	}

	linksConfig, err := ReadLinksConfigFromEnvVariables()
	if err != nil {
		return nil, err
	}

	return &Config{
		CredentialsFilePath: utils.GetEnvOr(EnvCredentialsFilePath, ""),
		Notifications:       notificationsConfig,
		Links:               linksConfig,
	}, nil
}

// -------------------------------------------------------------------------------------------------------------------

type NotificationsConfig struct {
	// SendNotificationField tells whether the notification field should be included or not.
	// If this field is set to false, the notification field will not be included in the payload and all its data
	// will instead be included inside the data field.
	SendNotificationField bool `json:"send_notification_field" yaml:"send_notification_field" toml:"send_notification_field"`
}

func ReadNotificationsConfigFromEnvVariables() (*NotificationsConfig, error) {
	sendNotificationsField, err := strconv.ParseBool(utils.GetEnvOr(EnvNotificationsSendNotificationField, "false"))
	if err != nil {
		return nil, err
	}

	return &NotificationsConfig{
		SendNotificationField: sendNotificationsField,
	}, nil
}

// -------------------------------------------------------------------------------------------------------------------

type LinksConfig struct {
	Domain  string      `json:"domain" yaml:"domain" toml:"domain"`
	Desktop DesktopInfo `json:"desktop" yaml:"desktop" toml:"desktop"`
	Android AndroidInfo `json:"android" yaml:"android" toml:"android"`
	Ios     IosInfo     `json:"ios" yaml:"ios" toml:"ios"`
}

func ReadLinksConfigFromEnvVariables() (*LinksConfig, error) {
	desktopInfo, err := ReadDesktopInfoFromEnvVariables()
	if err != nil {
		return nil, err
	}

	androidInfo, err := ReadAndroidInfoFromEnvVariables()
	if err != nil {
		return nil, err
	}

	iosInfo, err := ReadIosInfoFromEnvVariables()
	if err != nil {
		return nil, err
	}

	return &LinksConfig{
		Domain:  utils.GetEnvOr(EnvLinksDomain, ""),
		Desktop: desktopInfo,
		Android: androidInfo,
		Ios:     iosInfo,
	}, nil
}

type DesktopInfo struct {
	FallbackLink string `json:"fallback_link" yaml:"fallback_link" toml:"fallback_link"`
}

func ReadDesktopInfoFromEnvVariables() (DesktopInfo, error) {
	return DesktopInfo{
		FallbackLink: utils.GetEnvOr(EnvLinksDesktopFallbackLink, ""),
	}, nil
}

type AndroidInfo struct {
	PackageName           string `json:"package_name" yaml:"package_name" toml:"package_name"`
	MinPackageVersionCode string `json:"min_package_version_code" yaml:"min_package_version_code" toml:"min_package_version_code"`
}

func ReadAndroidInfoFromEnvVariables() (AndroidInfo, error) {
	return AndroidInfo{
		PackageName:           utils.GetEnvOr(EnvLinksAndroidPackageName, ""),
		MinPackageVersionCode: utils.GetEnvOr(EnvLinksAndroidMinPackageVersionCode, ""),
	}, nil
}

type IosInfo struct {
	BundleID       string `json:"bundle_id" yaml:"bundle_id" toml:"bundle_id"`
	MinimumVersion string `json:"minimum_version" yaml:"minimum_version" toml:"minimum_version"`
	AppStoreID     string `json:"app_store_id" yaml:"app_store_id" toml:"app_store_id"`
}

func ReadIosInfoFromEnvVariables() (IosInfo, error) {
	return IosInfo{
		BundleID:       utils.GetEnvOr(EnvLinksIOSBundleID, ""),
		MinimumVersion: utils.GetEnvOr(EnvLinksIOSMinimumVersion, ""),
		AppStoreID:     utils.GetEnvOr(EnvLinksIOSAppStoreID, ""),
	}, nil
}
