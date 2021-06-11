package bot

import (
    "regexp"
)

func Socialize(b *Bot, m *Message, u *User) {
    const mediaURL = "https://i.imgur.com/PQH5x0L.jpg"
    const meetURL = "https://meet.google.com/xxx-xxxx-xxx"
    const mediaAltText = "SOCIALIZE SOCIALIZE"

    if !regexp.MustCompile(`(?i)socialize socialize`).MatchString(m.Message) {
        return
    }

    var allowedRooms = map[string]bool {
        "just-yelling-training": true,
    }
    var Channel, _ = b.GetChannel(m.Channel)
    if !allowedRooms[Channel.Name] {
        return
    }

    m.Message = mediaAltText
    attachments := []MessageAttachment {
        MessageAttachment{
            Title: mediaAltText,
            TitleLink: meetURL,
            Fallback: meetURL,
            ImageURL: mediaURL,
            TS: m.TS,
            Actions: []MessageAttachmentAction {
                MessageAttachmentAction{
                    Name: mediaAltText,
                    Text: mediaAltText,
                    Type: "button",
                    Style: "primary",
                    Url: meetURL,
                },
            },
        },
    }

    b.SendAttachments(m, attachments)
}
