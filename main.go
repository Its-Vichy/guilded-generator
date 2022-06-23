package main

import (
	"fmt"
	"math/rand"
	//"os"
	"os/exec"
	//"os/signal"
	"strings"
	//"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/felixstrobel/mailtm"
	"github.com/its-vichy/guildedGen/package/guilded"
	"github.com/its-vichy/guildedGen/package/utils"
	"github.com/zenthangplus/goccm"
)

var (
	ThreadNumber = 50 // okey, only 5 thread but this gen was OP AS FUCK (with good proxies), you can do 5/s with only 5threads so don't worry :p

	MailAddr     = "lmfaonotamailtemp"
	MailPassword = "lollol"
	MailDomain   = "@knowledgemd.com"

	MailBox = map[string]string{}

	InviteCode = "ElQv91P2"
)

// couters
var (
	Generated int
	Verified  int
)

func FetchMailBox() {
	Client, _ := mailtm.NewMailClient()
	Client.GetAuthTokenCredentials(MailAddr+MailDomain, MailPassword)

	for {
		Messages, err := Client.GetMessages(1)

		if err != nil {
			fmt.Println("get mess: " + string(err.Error()))
			continue
		}

		for _, Message := range Messages {
			Mess, err := Client.GetMessageByID(Message.ID)

			if err != nil {
				//fmt.Println("get mess by id: " + string(err.Error()))
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
		exec.Command("cmd", "/C", "title", fmt.Sprintf("GuildeadGenerator - Generated: %d Verified: %d - Proxies: %d", Generated, Verified, len(utils.Proxies))).Run()
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	go FetchMailBox()
	go UpdateTitle()

	c := goccm.New(ThreadNumber)

	for {
		c.Wait()

		go func() {
			Proxy := utils.GetNexProxie()
			Session := guilded.CreateSession(Proxy)

			Email := MailAddr + "+" + utils.RandHexString(5) + MailDomain
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

								/*if Session.SetBio(utils.GetNexBio()) {
									color.HiMagenta("Email: %s Pass: %s | Name: %s ID: %s | Bio set\n", Me.User.Email, Pass, Me.User.Name, Me.User.ID)
								}*/

								Session.SetAvatar(utils.GetNexPfP())
								Session.SetActivity(1 + rand.Intn(3-1))
								Session.SetPlay(utils.GetNexStatus(), 90002200+rand.Intn(90002539-90002200))
								Session.Ping()

								//Session.JoinGuild(InviteCode)

								color.Green("%d | %s:%s\n", Verified, Me.User.Email, Pass)

								/*if Session.JoinGuild(InviteCode) {
									go color.Cyan("%d | Email: %s Pass: %s | Name: %s ID: %s | Joined server\n", Verified, Me.User.Email, Pass, Me.User.Name, Me.User.ID)
								}*/
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

	/*c.WaitAllDone()

	Sc := make(chan os.Signal, 1)
	signal.Notify(Sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-Sc*/
}
