package emstrigger

import (
	"context"

	"github.com/mmussett/ems"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/trigger"
	"github.com/prometheus/common/log"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &Output{})

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

type Factory struct {
}

// Metadata implements trigger.Factory.Metadata
func (*Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// New implements trigger.Factory.New
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)

	if err != nil {
		return nil, err
	}

	return &Trigger{id: config.Id, settings: s}, nil
}

type Trigger struct {
	settings *Settings
	id       string
	Client   *ems.Client
	Handlers []trigger.Handler
}

func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	ops := ems.NewClientOptions().SetServerUrl(t.settings.ServerURL).SetUsername(t.settings.Username).SetPassword(t.settings.Password)

	client := ems.NewClient(ops).(*ems.Client)

	err := client.Connect()

	if err != nil {
		ctx.Logger().Error(err)
		return err
	}
	t.Client = client

	for _, handler := range ctx.GetHandlers() {
		t.Handlers = append(t.Handlers, handler)
	}

	return nil
}

func (t *Trigger) Start() error {

	return t.startHandlers()
}

// Stop implements util.Managed.Stop
func (t *Trigger) Stop() error {
	t.Client.Disconnect()

	return nil
}

func (t *Trigger) startHandlers() error {

	for {

		msgText, err := t.Client.Receive(t.settings.Destination)
		if err != nil {

		}
		out := &Output{}

		out.Data = msgText

		for _, handler := range t.Handlers {

			_, err := handler.Handle(context.Background(), out)
			if err != nil {
				log.Error("Error starting action: %v", err)
			}
		}

	}
	return nil
}
