package main

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/felixstrobel/mailtm"
	"github.com/its-vichy/guildedGen/package/guilded"
	"github.com/its-vichy/guildedGen/package/utils"
	"github.com/zenthangplus/goccm"
)

var (
	ThreadNumber = 50 // okey, only 5 thread but this gen was OP AS FUCK (with good proxies), you can do 5/s with only 5threads so don't worry :p

	MailAddr     = "gerdgrezgergredsfs"
	MailPassword = "gerdgrezgergredsfs"
	MailDomain   = "@knowledgemd.com"

	MailBox = map[string]string{}

	InviteCode = "pYr9V8YE"
)

// couters
var (
	Generated int
	Verified  int
	MailErrors int
)

func FetchMailBox(username string, password string) {
	Client, _ := mailtm.NewMailClient()
	
	url_i := url.URL{}
	url_proxy, _ := url_i.Parse(fmt.Sprintf("http://%s", utils.GetNexProxie()))

	transport := http.Transport{}
	transport.Proxy = http.ProxyURL(url_proxy)
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	http_client := http.Client{Transport: &transport}

	Client.HttpClient = &http_client
	Client.GetAuthTokenCredentials(username, password)

	for {
		Messages, err := Client.GetMessages(1)

		if err != nil {
			//fmt.Println("get mess: " + string(err.Error()))
			MailErrors++
			continue
		}

		for _, Message := range Messages {
			Mess, err := Client.GetMessageByID(Message.ID)

			if err != nil {
				continue
			}

			if strings.Contains(Mess.Subject, "Welcome to Guilded") {
				go Client.DeleteMessageByID(Message.ID)
				continue
			}

			if Mess.Subject == "Verify your email on Guilded" {
				VerificationToken := strings.Split(strings.Split(Mess.Html[0], "https://www.guilded.gg/api/email/verify?token=")[1], `"`)[0]
				go Client.DeleteMessageByID(Message.ID)

				MailBox[Mess.To[0].Address] = VerificationToken
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func UpdateTitle() {
	for {
		exec.Command("cmd", "/C", "title", fmt.Sprintf("GuildeadGenerator - Generated: %d Verified: %d - Proxies: %d - MailErrors: %d", Generated, Verified, len(utils.Proxies), MailErrors)).Run()
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	for _, mail := range utils.Emails {
		go func(mail string) {
			splitted := strings.Split(mail, ":")

			go FetchMailBox(splitted[0], splitted[1])
			go UpdateTitle()

			c := goccm.New(ThreadNumber)

			for {
				c.Wait()

				go func() {
					Proxy := utils.GetNexProxie()
					Session := guilded.CreateSession(Proxy)

					splt := strings.Split(splitted[0], "@")

					Email := splt[0] + "+" + utils.RandHexString(5) + "@" + splt[1]
					Pass := utils.RandHexString(5)
					Username := utils.GetNexUsername()

					r := Session.CreateAccount(Email, Pass, Username)

					if r.User.Email == "" {
						utils.Proxies = utils.RemoveIProxy(Proxy, utils.Proxies)
						c.Done()
						return
					}

					Me := Session.Login(Email, Pass)

					if Me.User.Email == "" {
						c.Done()
						return
					}

					color.Yellow("%d | %s:%s\n", Verified, Me.User.Email, Pass)
					Generated++

					if Session.SpoofEvent() {
						Session.SentVerificationMail()
						IsVerified := false

						for IsVerified == false {
							for key, value := range MailBox {
								if key == Email {
									if Session.VerifyEmail(value) {
										utils.AppendLine("./data/tokens.txt", fmt.Sprintf("%s:%s:%s:%s", Email, Pass, Session.HttpCookies["hmac_signed_session"], Me.User.ID))

										delete(MailBox, key)
										IsVerified = true
										Verified++

										Session.SetAvatar(utils.GetNexPfP())
										Session.SetActivity(1 + rand.Intn(3-1))
										Session.SetPlay(utils.GetNexStatus(), 90002200+rand.Intn(90002539-90002200))
										Session.Ping()

										go Session.JoinGuild(InviteCode)
										//go Session.JoinTeam(InviteCode)

										color.Green("%d | %s:%s\n", Verified, Me.User.Email, Pass)
									} else {
										IsVerified = true
									}
								}
							}

							time.Sleep(500 * time.Millisecond)
						}
					}

					c.Done()
				}()
			}
		}(mail)
	}

	Sc := make(chan os.Signal, 1)
	signal.Notify(Sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-Sc
}
