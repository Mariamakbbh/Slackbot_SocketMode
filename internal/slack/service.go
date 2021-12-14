package slack

import (
	"fmt"
	"strings"
	"time"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/zap"
)

type Service struct {
	client       *slack.Client
	socketClient *socketmode.Client

	logger *zap.SugaredLogger
}

func (s *Service) Flow(event socketmode.Event) error {
	switch event.Type {
	case socketmode.EventTypeEventsAPI:
		// The Event sent on the channel is not the same as the EventAPI events so we need to type cast it
		eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)

		if !ok {
			s.logger.Errorf("could not type cast the event to the EventsAPIEvent: %v", event)
			return fmt.Errorf("could not type cast the event to the EventsAPIEvent: %v", event)
		}

		// We need to send an Acknowledge to the slack server
		s.socketClient.Ack(*event.Request)
		// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
		// log.Println(eventsAPIEvent)
		// Now we have an Events API event, but this event type can in turn be many types, so we actually need another type switch
		s.logger.Info("* WebSocketBOT received an event")

		if err := s.handleEventMessage(eventsAPIEvent); err != nil {
			// Replace with actual err handling
			s.logger.Error(err)
		}

	default:
		s.logger.Infof("unknown event type %s", event.Type)
	}

	return nil
}

// handleEventMessage will take an event and handle it properly based on the type of event
func (s *Service) handleEventMessage(event slackevents.EventsAPIEvent) error {
	switch event.Type {
	// First we check if this is an CallbackEvent
	case slackevents.CallbackEvent:
		s.logger.Infof("* It's a callback event")
		innerEvent := event.InnerEvent

		// Yet Another Type switch on the actual Data to see if its an AppMentionEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			// The application has been mentioned since this Event is a Mention event
			s.logger.Infof("* WebSocketBOT did they mention you in Slack?")

			if err := s.handleAppMentionEvent(ev); err != nil {
				return err
			}

		default:
			s.logger.Infof("* It's callback but not App Mention Event")
		}
	default:
		return fmt.Errorf("unsupported event type")
	}

	return nil
}

// handleAppMentionEvent is used to take care of the AppMentionEvent when the bot is mentioned
func (s *Service) handleAppMentionEvent(event *slackevents.AppMentionEvent) error {
	s.logger.Infof("* WebSocketBOT is this even running?")
	// Grab the user name based on the ID of the one who mentioned the bot
	user, err := s.client.GetUserInfo(event.User)
	if err != nil {
		return err
	}
	// Check if the user said Hello to the bot
	text := strings.ToLower(event.Text)
	s.logger.Infof("* WebSocketBOT who was the user that needs your help and what they want: ", user, text)
	// Create the attachment and assigned based on the message
	attachment := slack.Attachment{}
	// Add Some default context like user who mentioned the bot
	attachment.Fields = []slack.AttachmentField{
		{
			Title: "Date",
			Value: time.Now().String(),
		}, {
			Title: "Initializer",
			Value: user.Name,
		},
	}
	if strings.Contains(text, "hello") {
		// Greet the user
		attachment.Text = fmt.Sprintf("Hello %s", user.Name)
		attachment.Pretext = "Greetings"
		attachment.Color = "#4af030"
	} else {
		// Send a message to the user
		attachment.Text = fmt.Sprintf("How can I help you %s?", user.Name)
		attachment.Pretext = "How can I be of service"
		attachment.Color = "#3d3d3d"
	}
	// Send the message to the channel
	// The Channel is available in the event message
	if _, _, err = s.client.PostMessage(event.Channel, slack.MsgOptionAttachments(attachment)); err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}

	return nil
}

func NewService(client *slack.Client, socketClient *socketmode.Client, logger *zap.SugaredLogger) ServiceInterface {
	svc := new(Service)

	svc.client = client
	svc.socketClient = socketClient
	svc.logger = logger

	return svc
}
