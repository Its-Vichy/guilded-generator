package guilded

import (
	"time"

	"github.com/valyala/fasthttp"
)

type GuildedSession struct {
	HttpClient  *fasthttp.Client
	HttpHeader  fasthttp.RequestHeader
	HttpCookies map[string]string

	Client LoginResponse
}

type RegisterPayload struct {
	ExtraInfo ExtraInfo `json:"extraInfo"`

	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
}

type ExtraInfo struct {
	Platform string `json:"platform,omitempty"`
}

type RegisterResponse struct {
	User struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	GetMe    bool   `json:"getMe"`
}

type LoginResponse struct {
	UpdateMessage interface{} `json:"updateMessage"`
	User          struct {
		ID                 string        `json:"id"`
		Name               string        `json:"name"`
		Subdomain          interface{}   `json:"subdomain"`
		Aliases            []interface{} `json:"aliases"`
		Email              string        `json:"email"`
		ProfilePictureSm   interface{}   `json:"profilePictureSm"`
		ProfilePicture     interface{}   `json:"profilePicture"`
		ProfilePictureLg   interface{}   `json:"profilePictureLg"`
		ProfilePictureBlur interface{}   `json:"profilePictureBlur"`
		ProfileBannerBlur  interface{}   `json:"profileBannerBlur"`
		ProfileBannerLg    interface{}   `json:"profileBannerLg"`
		ProfileBannerSm    interface{}   `json:"profileBannerSm"`
		JoinDate           time.Time     `json:"joinDate"`
		SteamID            interface{}   `json:"steamId"`

		UserStatus struct {
			Content          interface{} `json:"content"`
			CustomReactionID interface{} `json:"customReactionId"`
		} `json:"userStatus"`

		ModerationStatus           interface{}   `json:"moderationStatus"`
		AboutInfo                  interface{}   `json:"aboutInfo"`
		LastOnline                 time.Time     `json:"lastOnline"`
		BlockedUsers               []interface{} `json:"blockedUsers"`
		SocialLinks                []interface{} `json:"socialLinks"`
		UserPresenceStatus         int           `json:"userPresenceStatus"`
		Badges                     []interface{} `json:"badges"`
		HasSeenServerSubscriptions bool          `json:"hasSeenServerSubscriptions"`
		ServerSubscriptions        []interface{} `json:"serverSubscriptions"`
		IsUnrecoverable            bool          `json:"isUnrecoverable"`

		Devices []struct {
			Type       string    `json:"type"`
			ID         string    `json:"id"`
			LastOnline time.Time `json:"lastOnline"`
			IsActive   bool      `json:"isActive"`
		} `json:"devices"`

		UserChannelNotificationSettings []interface{} `json:"userChannelNotificationSettings"`

		Upsell struct {
			Type                string        `json:"type"`
			ActivationType      string        `json:"activationType"`
			Topic               string        `json:"topic"`
			IsAE                bool          `json:"isAE"`
			IncludedUpsellSpecs []interface{} `json:"includedUpsellSpecs"`

			LocalStageStats struct {
				GetDesktopApp     string `json:"getDesktopApp"`
				GetMobileApp      string `json:"getMobileApp"`
				AddProfilePicture string `json:"addProfilePicture"`
				CreateOwnServer   string `json:"createOwnServer"`
				ShowSwipeNux      string `json:"showSwipeNux"`
				ReferFriend       string `json:"referFriend"`
			} `json:"localStageStats"`

			EntityID          string        `json:"entityId"`
			IncludedUpsells   []interface{} `json:"includedUpsells"`
			CreatedAtUnixSecs string        `json:"createdAtUnixSecs"`
		} `json:"upsell"`

		Characteristics struct {
			IsChatActivated bool `json:"isChatActivated"`
		} `json:"characteristics"`
	} `json:"user"`

	Teams           []interface{} `json:"teams"`
	CustomReactions []interface{} `json:"customReactions"`
	ReactionUsages  []interface{} `json:"reactionUsages"`
	LandingURL      bool          `json:"landingUrl"`
	Friends         []interface{} `json:"friends"`
}

type MailVerificationSendResponse struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
}