package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func HandleBookAppointment(uc *UserController, s *discordgo.Session, m *discordgo.MessageCreate, details string) {
	parts := strings.SplitN(details, " ", 3)
	if len(parts) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Invalid booking format. Please use '!book date time details'.")
		return
	}

	date, timePart, detail := parts[0], parts[1], parts[2]
	appointmentTime, err := time.Parse("2006-01-02 15:04", date+" "+timePart)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed to parse date and time. Please use the format YYYY-MM-DD HH:MM.")
		fmt.Println("Error parsing date and time:", err)
		return
	}

	err = uc.BookAppointment(m.Author.ID, date, timePart, detail)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed to book appointment.")
		fmt.Println("Error booking appointment:", err)
		return
	}

	go func() {
		durationUntilAppointment := time.Until(appointmentTime)
		time.Sleep(durationUntilAppointment - 1*time.Hour) 

		dmChannel, err := s.UserChannelCreate(m.Author.ID)
		if err == nil {
			s.ChannelMessageSend(dmChannel.ID, fmt.Sprintf("Reminder: You have an appointment on %s. Details: %s", appointmentTime.Format("2006-01-02 15:04"), detail))
		}
	}()

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Appointment booked for %s at %s: %s. You will be reminded 1 hour before the appointment.", date, timePart, detail))
}

func HandleUploadDocument(uc *UserController, s *discordgo.Session, m *discordgo.MessageCreate, args string) {
    if len(m.Attachments) == 0 {
        s.ChannelMessageSend(m.ChannelID, "No document attached.")
        return
    }

    for _, attachment := range m.Attachments {
        resp, err := http.Get(attachment.URL)
        if err != nil {
            fmt.Println("Error downloading document:", err)
            continue
        }
        defer resp.Body.Close()

        savePath := fmt.Sprintf("/uploads/%s", attachment.Filename) 

        file, err := os.Create(savePath)
        if err != nil {
            fmt.Println("Error creating file:", err)
            continue
        }
        defer file.Close()

        _, err = io.Copy(file, resp.Body)
        if err != nil {
            fmt.Println("Error saving document:", err)
            continue
        }
        err = uc.UploadDocument(m.Author.ID, savePath)
        if err != nil {
            fmt.Println("Error saving document info to database:", err)
            s.ChannelMessageSend(m.ChannelID, "Failed to upload document.")
            continue
        }

        s.ChannelMessageSend(m.ChannelID, "Document uploaded successfully.")
    }
}

func (bot *Bot) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID == bot.Session.State.User.ID {
        return 
    }

    content := strings.TrimSpace(m.Content)
    if !strings.HasPrefix(content, "!") {
        return 
    }

    parts := strings.SplitN(content[len("!"):], " ", 2)
    command := parts[0]
    var args string
    if len(parts) > 1 {
        args = parts[1]
    }

    switch command {
    case "book":
        HandleBookAppointment(bot.UserController, bot.Session, m, args)
    case "upload":
        HandleUploadDocument(bot.UserController, bot.Session, m, args)
    default:
        bot.Session.ChannelMessageSend(m.ChannelID, "Unknown command. Use !help to list all commands.")
    }
}