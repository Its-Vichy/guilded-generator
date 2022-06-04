package main

import (
	"fmt"
	"math/rand"
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
)

var (
	ThreadNumber = 5 // okey, only 5 thread but this gen was OP AS FUCK, you can do 5/s with only 5threads so don't worry :p

	MailAddr     = "the name here"
	MailPassword = "you care.."
	MailDomain   = "@knowledgemd.com"

	MailBox      = map[string]string{}

	InviteCode = "EoeQbWPk"
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

		time.Sleep(1 * time.Second)
	}
}

func UpdateTitle() {
	for {
		exec.Command("cmd", "/C", "title", fmt.Sprintf("GuildeadGenerator - Generated: %d Verified: %d", Generated, Verified)).Run()
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	go FetchMailBox()
	go UpdateTitle()

	for i := 1; i <= ThreadNumber; i++ {
		go func() {
			for {
				Session := guilded.CreateSession(utils.GetNexProxie())

				Email := MailAddr + "+" + utils.RandHexString(5) + MailDomain
				Pass := utils.RandHexString(5)
				Username := utils.GetNexUsername()

				Session.CreateAccount(Email, Pass, Username)
				Me := Session.Login(Email, Pass)

				color.Yellow("%d | Email: %s Pass: %s | Name: %s ID: %s | #%d\n", Verified, Me.User.Email, Pass, Me.User.Name, Me.User.ID, Generated)
				Generated++

				if Session.SpoofEvent() {
					go func() {
						Session.SentVerificationMail()
						IsVerified := false

						for IsVerified == false {
							for key, value := range MailBox {
								if key == Email {
									if Session.VerifyEmail(value) {
										color.Green("%d | Email: %s Pass: %s | Name: %s ID: %s | #%d\n", Verified, Me.User.Email, Pass, Me.User.Name, Me.User.ID, Verified)
										utils.AppendLine("tokens.txt", fmt.Sprintf("%s:%s:%s:%s", Email, Pass, Session.HttpCookies["hmac_signed_session"], Me.User.ID))
										
										delete(MailBox, key)
										
										IsVerified = true
										Verified++
										
										if Session.SetAvatar(utils.GetNexPfP()) {
											go color.Magenta("%d | Email: %s Pass: %s | Name: %s ID: %s | Avatar set\n", Verified, Me.User.Email, Pass, Me.User.Name, Me.User.ID)
										}

										/*if Session.SetBio(utils.GetNexBio()) {
											color.HiMagenta("Email: %s Pass: %s | Name: %s ID: %s | Bio set\n", Me.User.Email, Pass, Me.User.Name, Me.User.ID)
										}*/

										if Session.SetActivity(1 + rand.Intn(3-1)) {
											go color.Blue("%d | Email: %s Pass: %s | Name: %s ID: %s | Set activity\n", Verified, Me.User.Email, Pass, Me.User.Name, Me.User.ID)
										}

										if Session.SetPlay(utils.GetNexStatus(), 90002200 + rand.Intn(90002539-90002200)) {
											go color.HiBlue("%d | Email: %s Pass: %s | Name: %s ID: %s | Set game\n", Verified, Me.User.Email, Pass, Me.User.Name, Me.User.ID)
										}

										if Session.Ping() {
											go color.HiBlue("%d | Email: %s Pass: %s | Name: %s ID: %s | Ping sent\n", Verified, Me.User.Email, Pass, Me.User.Name, Me.User.ID)
										}

										if Session.JoinGuild(InviteCode) {
											go color.Cyan("%d | Email: %s Pass: %s | Name: %s ID: %s | Joined server\n", Verified, Me.User.Email, Pass, Me.User.Name, Me.User.ID)
										}
									}
								}
							}

							time.Sleep(1 * time.Second)
						}
					}()
				}
			}
		}()
	}

	Sc := make(chan os.Signal, 1)
	signal.Notify(Sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-Sc
}
