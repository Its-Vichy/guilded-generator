package guilded

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/dgrr/cookiejar"

	"github.com/google/uuid"
	"github.com/its-vichy/guildedGen/package/utils"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

func CreateSession(Proxy string) *GuildedSession {
	session := &GuildedSession{
		HttpClient: &fasthttp.Client{
			Dial: fasthttpproxy.FasthttpHTTPDialer(Proxy),
		},
		HttpCookies: make(map[string]string),
	}

	header := map[string]string{
		"Sec-Ch-Ua":           `" Not;A Brand";v="99", "Google Chrome";v="97", "Chromium";v="97"`,
		"Sec-Ch-Ua-Mobile":    `?0`,
		"User-Agent":          `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36`,
		"Content-Type":        `application/json`,
		"Accept":              `application/json, text/javascript, **; q=0.01`,
		"X-Requested-With":    `XMLHttpRequest`,
		"Sec-Ch-Ua-Platform":  `"macOS"`,
		"Origin":              `https://www.guilded.gg`,
		"Sec-Fetch-Site":      `same-origin`,
		"Sec-Fetch-Mode":      `cors`,
		"Sec-Fetch-Dest":      `empty`,
		"Referer":             "https://www.guilded.gg/",
		"guilded-client-id":   uuid.New().String(),
		"guilded-device-id":   utils.RandHexString(64),
		"guilded-device-type": "desktop",
		"Accept-Language":     `fr-FR,fr;q=0.9`,
	}

	for key, value := range header {
		session.HttpHeader.Add(key, value)
	}

	return session
}

func (Session GuildedSession) PostRequest(url string, payload []byte, method string) *fasthttp.Response {
	for {
		request := fasthttp.AcquireRequest()
		request.Header = Session.HttpHeader

		cj := cookiejar.AcquireCookieJar()
		defer cookiejar.ReleaseCookieJar(cj)

		request.Header.Add("Content-Length", fmt.Sprint(len(payload)))
		request.SetRequestURI(url)
		request.Header.SetMethod(method)
		request.SetBodyRaw(payload)

		for key, value := range Session.HttpCookies {
			request.Header.SetCookie(key, value)
		}

		response := fasthttp.AcquireResponse()
		err := Session.HttpClient.Do(request, response)

		fasthttp.ReleaseRequest(request)

		if err != nil {
			//fmt.Printf("error: %s\n", err)
			break
		}

		cj.ReadResponse(response)

		for {
			c := cj.Get()
			if c == nil {
				break
			}

			Session.HttpCookies[string(c.Key())] = string(c.Value())
			fasthttp.ReleaseCookie(c)
		}

		return response
	}

	return &fasthttp.Response{}
}

func (Session GuildedSession) CreateAccount(Email string, Password string, Name string) *RegisterResponse {
	payload, _ := json.Marshal(&RegisterPayload{
		ExtraInfo: ExtraInfo{
			Platform: "desktop",
		},

		Email:    Email,
		Name:     Name,
		FullName: Name,
		Password: Password,
	})

	rdata := Session.PostRequest("https://www.guilded.gg/api/users?type=email", payload, fasthttp.MethodPost)
	response := &RegisterResponse{}

	json.Unmarshal(rdata.Body(), response)
	return response
}

func (Session GuildedSession) Login(Email string, Password string) *LoginResponse {
	payload, _ := json.Marshal(&LoginPayload{
		Email:    Email,
		Password: Password,
		GetMe:    true,
	})

	rdata := Session.PostRequest("https://www.guilded.gg/api/login", payload, fasthttp.MethodPost)
	response := &LoginResponse{}

	json.Unmarshal(rdata.Body(), response)
	Session.Client = *response

	return response
}

func (Session GuildedSession) SentVerificationMail() *MailVerificationSendResponse {
	rdata := Session.PostRequest("https://www.guilded.gg/api/email/verify", nil, fasthttp.MethodPost)
	response := &MailVerificationSendResponse{}

	json.Unmarshal(rdata.Body(), response)
	return response
}

// Magic payload lol
func (Session GuildedSession) SpoofEvent() bool {
	payload, err := ioutil.ReadAll(strings.NewReader(fmt.Sprintf(`{"data":[{"registrationType":"email","success":true,"durationInMs":774,"userId": "%s","loginType":"email","time":1654089986341,"eventSource":"Client","eventName":"Login","viewerPlatform":"desktop","viewerAppType":"null","viewerSystemName":"Win32","browser":"Firefox","device":null,"attributionSource":null,"gitHash":"ec484bed"},{"name":"SignUpLogInOverlay","confirmed":true,"time":1654089986353,"eventSource":"Client","eventName":"OverlayClosed","viewerPlatform":"desktop","viewerAppType":"null","viewerSystemName":"Win32","browser":"Firefox","device":null,"attributionSource":null,"gitHash":"ec484bed"},{"name":"SignUpLogInOverlay","confirmed":false,"time":1654089986522,"eventSource":"Client","eventName":"OverlayClosed","viewerPlatform":"desktop","viewerAppType":"null","viewerSystemName":"Win32","browser":"Firefox","device":null,"attributionSource":null,"gitHash":"ec484bed"},{"notificationPermission":"granted","time":1654089986795,"eventSource":"Client","eventName":"ServiceWorkerRegistrationSuccess","viewerPlatform":"desktop","viewerAppType":"null","viewerSystemName":"Win32","browser":"Firefox","device":null,"attributionSource":null,"gitHash":"ec484bed"}]}`, Session.Client.User.ID)))

	if err != nil {
		fmt.Println(err)
		return false
	}

	response := Session.PostRequest("https://www.guilded.gg/api/data/event", payload, fasthttp.MethodPut)

	if response.StatusCode() == 200 {
		return true
	} else {
		return false
	}
}

func (Session GuildedSession) VerifyEmail(VerificationToken string) bool {
	response := Session.PostRequest(fmt.Sprintf("https://www.guilded.gg/api/email/verify?token=%s", VerificationToken), nil, fasthttp.MethodGet)

	if response.StatusCode() == 302 {
		return true
	} else {
		return false
	}
}

func (Session GuildedSession) JoinGuild(InviteCode string) bool {
	response := Session.PostRequest(fmt.Sprintf("https://www.guilded.gg/api/invites/%s", InviteCode), nil, fasthttp.MethodPut)

	if response.StatusCode() == 200 {
		return true
	} else {
		return false
	}
}

func (Session GuildedSession) SetAvatar(AvatarUrl string) bool {
	payload, err := ioutil.ReadAll(strings.NewReader(fmt.Sprintf(`{"imageUrl": "%s"}`, AvatarUrl)))

	if err != nil {
		fmt.Println(err)
		return false
	}

	response := Session.PostRequest("https://www.guilded.gg/api/users/me/profile/images", payload, fasthttp.MethodPost)

	if response.StatusCode() == 200 {
		return true
	} else {
		return false
	}
}

func (Session GuildedSession) SetBio(Content string) bool {
	payload, err := ioutil.ReadAll(strings.NewReader(fmt.Sprintf(`{"userId":"%s","aboutInfo":{"bio":"%s","tagLine":"%s"}}`, Session.Client.User.ID, Content, Content)))

	if err != nil {
		fmt.Println(err)
		return false
	}

	response := Session.PostRequest(fmt.Sprintf("https://www.guilded.gg/api/users/%s/profilev2", Session.Client.User.ID), payload, fasthttp.MethodPut)

	if response.StatusCode() == 200 {
		return true
	} else {
		return false
	}
}

func (Session GuildedSession) SetActivity(Activity int) bool {
	payload, err := ioutil.ReadAll(strings.NewReader(fmt.Sprintf(`{"status": %d}`, Activity)))

	if err != nil {
		fmt.Println(err)
		return false
	}

	response := Session.PostRequest("https://www.guilded.gg/api/users/me/presence", payload, fasthttp.MethodPost)

	if response.StatusCode() == 200 {
		return true
	} else {
		return false
	}
}

func (Session GuildedSession) Ping() bool {
	response := Session.PostRequest("https://www.guilded.gg/api/users/me/ping", nil, fasthttp.MethodPut)

	if response.StatusCode() == 200 {
		return true
	} else {
		return false
	}
}

func (Session GuildedSession) SetPlay(Content string, CustomReactionId int) bool {
	payload, err := ioutil.ReadAll(strings.NewReader(fmt.Sprintf(`{"content":{"object":"value","document":{"object":"document","data":{},"nodes":[{"object":"block","type":"paragraph","data":{},"nodes":[{"object":"text","leaves":[{"object":"leaf","text":"%s","marks":[]}]}]}]}},"customReactionId":%d,"expireInMs":0}`, Content, CustomReactionId)))

	if err != nil {
		fmt.Println(err)
		return false
	}

	response := Session.PostRequest("https://www.guilded.gg/api/users/me/status", payload, fasthttp.MethodPost)

	if response.StatusCode() == 200 {
		return true
	} else {
		return false
	}
}