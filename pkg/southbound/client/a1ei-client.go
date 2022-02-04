// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package sbclient

import (
	"context"
	"github.com/google/uuid"
	"github.com/onosproject/onos-a1t/pkg/stream"
	"github.com/onosproject/onos-api/go/onos/a1t/a1"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"io"
)

var a1eiLog = logging.GetLogger("southbound", "a1ei-client")

func NewA1EIClient(ctx context.Context, targetXAppID string, ipAddress string, port uint32, streamBroker stream.Broker) (Client, error) {
	// create gRPC client
	conn, err := createGRPCConn(ipAddress, port)
	if err != nil {
		return nil, err
	}

	// create broker session
	createStream(ctx, targetXAppID, stream.EnrichmentInformation, streamBroker)

	return &a1eiClient{
		targetXAppID: targetXAppID,
		ipAddress:    ipAddress,
		port:         port,
		grpcClient:   a1.NewEIServiceClient(conn),
		streamBroker: streamBroker,
		sessions:     make(map[stream.A1SBIRPCType]interface{}),
	}, nil
}

type a1eiClient struct {
	targetXAppID string
	ipAddress    string
	port         uint32
	grpcClient   a1.EIServiceClient
	streamBroker stream.Broker
	sessions     map[stream.A1SBIRPCType]interface{}
}

func (a *a1eiClient) Run(ctx context.Context) error {
	err := a.createSessions(ctx)
	if err != nil {
		return err
	}

	a.runIncomingMsgForwarder(ctx)

	err = a.runOutgoingMsgDispatcher(ctx)
	if err != nil {
		a1eiLog.Warn(err)
		a.Close()
		return err
	}

	return nil
}

func (a *a1eiClient) createSessions(ctx context.Context) error {
	var err error

	err = a.createEIQuerySession(ctx)
	if err != nil {
		return err
	}

	err = a.createEIJobSetupSession(ctx)
	if err != nil {
		return err
	}

	err = a.createEIJobUpdateSession(ctx)
	if err != nil {
		return err
	}

	err = a.createEIJobDeleteSession(ctx)
	if err != nil {
		return err
	}

	err = a.createEIJobStatusQuerySession(ctx)
	if err != nil {
		return err
	}

	return err
}

func (a *a1eiClient) createEIQuerySession(ctx context.Context) error {
	s, err := a.grpcClient.EIQuery(ctx)
	if err != nil {
		a1eiLog.Warn(err)
		return err
	}

	a.sessions[stream.EIQuery] = s
	return nil
}

func (a *a1eiClient) createEIJobSetupSession(ctx context.Context) error {
	s, err := a.grpcClient.EIJobSetup(ctx)
	if err != nil {
		a1eiLog.Warn(err)
		return err
	}

	a.sessions[stream.EIJobSetup] = s
	return nil
}

func (a *a1eiClient) createEIJobUpdateSession(ctx context.Context) error {
	s, err := a.grpcClient.EIJobUpdate(ctx)
	if err != nil {
		a1eiLog.Warn(err)
		return err
	}

	a.sessions[stream.EIJobUpdate] = s
	return nil
}

func (a *a1eiClient) createEIJobDeleteSession(ctx context.Context) error {
	s, err := a.grpcClient.EIJobDelete(ctx)
	if err != nil {
		a1eiLog.Warn(err)
		return err
	}

	a.sessions[stream.EIJobDelete] = s
	return nil
}

func (a *a1eiClient) createEIJobStatusQuerySession(ctx context.Context) error {
	s, err := a.grpcClient.EIJobStatusQuery(ctx)
	if err != nil {
		a1eiLog.Warn(err)
		return err
	}

	a.sessions[stream.EIJobStatusQuery] = s
	return nil
}

func (a *a1eiClient) runIncomingMsgForwarder(ctx context.Context) {
	go a.incomingEIQueryForwarder(ctx)
	go a.incomingEIJobSetupForwarder(ctx)
	go a.incomingEIJobUpdateForwarder(ctx)
	go a.incomingEIJobDeleteForwarder(ctx)
	go a.incomingEIJobStatusQueryForwarder(ctx)
}

func (a *a1eiClient) incomingEIQueryForwarder(ctx context.Context) {
	defer a.Close()
	for {
		select {
		case <-ctx.Done():
			a1eiLog.Warn("A1EI SBI client incoming forwarder for EI Query service is just closed")
			return
		default:
			if _, ok := a.sessions[stream.EIQuery]; !ok {
				a1eiLog.Warn("A1EI SBI client incoming forwarder for EI Query service is just closed")
				return
			}
			msg, err := a.sessions[stream.EIQuery].(a1.EIService_EIQueryClient).Recv()
			if err == io.EOF || err == context.Canceled {
				a1eiLog.Warn("A1EI SBI client incoming forwarder for EI Query service is just closed")
				return
			}
			if err != nil {
				a1eiLog.Warn(err)
				return
			}
			sbMessage := &stream.SBStreamMessage{
				TargetXAppID:     a.targetXAppID,
				A1SBIMessageType: stream.EIRequestMessage,
				A1Service:        stream.EnrichmentInformation,
				A1SBIRPCType:     stream.EIQuery,
				Payload:          msg,
			}
			nbID := stream.ID{
				SrcEndpointID:  stream.GetEndpointIDWithTargetXAppID(a.targetXAppID, stream.EnrichmentInformation),
				DestEndpointID: "a1p-controller",
			}
			err = a.streamBroker.Send(nbID, sbMessage)
			if err != nil {
				a1eiLog.Warn(err)
			}
		}
	}
}

func (a *a1eiClient) incomingEIJobSetupForwarder(ctx context.Context) {
	defer a.Close()
	for {
		select {
		case <-ctx.Done():
			a1eiLog.Warn("A1EI SBI client incoming forwarder for EI Job Setup service is just closed")
			return
		default:
			if _, ok := a.sessions[stream.EIJobSetup]; !ok {
				a1eiLog.Warn("A1EI SBI client incoming forwarder for EI Job Setup service is just closed")
				return
			}
			msg, err := a.sessions[stream.EIJobSetup].(a1.EIService_EIJobSetupClient).Recv()
			if err == io.EOF || err == context.Canceled {
				a1eiLog.Warn("A1EI SBI client incoming forwarder for EI Job Setup service is just closed")
				return
			}
			if err != nil {
				a1eiLog.Warn(err)
				return
			}
			sbMessage := &stream.SBStreamMessage{
				TargetXAppID:     a.targetXAppID,
				A1SBIMessageType: stream.EIRequestMessage,
				A1Service:        stream.EnrichmentInformation,
				A1SBIRPCType:     stream.EIJobSetup,
				Payload:          msg,
			}
			nbID := stream.ID{
				SrcEndpointID:  stream.GetEndpointIDWithTargetXAppID(a.targetXAppID, stream.EnrichmentInformation),
				DestEndpointID: "a1p-controller",
			}
			err = a.streamBroker.Send(nbID, sbMessage)
			if err != nil {
				a1eiLog.Warn(err)
			}
		}
	}
}

func (a *a1eiClient) incomingEIJobUpdateForwarder(ctx context.Context) {
	defer a.Close()
	for {
		select {
		case <-ctx.Done():
			a1eiLog.Warn("A1EI SBI client incoming forwarder for EI Job Update service is just closed")
			return
		default:
			if _, ok := a.sessions[stream.EIJobUpdate]; !ok {
				a1eiLog.Warn("A1EI SBI client incoming forwarder for EI Job Update service is just closed")
				return
			}
			msg, err := a.sessions[stream.EIJobUpdate].(a1.EIService_EIJobUpdateClient).Recv()
			if err == io.EOF || err == context.Canceled {
				a1eiLog.Warn("A1EI SBI client incoming forwarder for EI Job Update service is just closed")
				return
			}
			if err != nil {
				a1eiLog.Warn(err)
				return
			}
			sbMessage := &stream.SBStreamMessage{
				TargetXAppID:     a.targetXAppID,
				A1SBIMessageType: stream.EIRequestMessage,
				A1Service:        stream.EnrichmentInformation,
				A1SBIRPCType:     stream.EIJobUpdate,
				Payload:          msg,
			}
			nbID := stream.ID{
				SrcEndpointID:  stream.GetEndpointIDWithTargetXAppID(a.targetXAppID, stream.EnrichmentInformation),
				DestEndpointID: "a1p-controller",
			}
			err = a.streamBroker.Send(nbID, sbMessage)
			if err != nil {
				a1eiLog.Warn(err)
			}
		}
	}
}

func (a *a1eiClient) incomingEIJobDeleteForwarder(ctx context.Context) {
	defer a.Close()
	for {
		select {
		case <-ctx.Done():
			a1eiLog.Warn("A1EI SBI client incoming forwarder for EI Job Delete service is just closed")
			return
		default:
			if _, ok := a.sessions[stream.EIJobDelete]; !ok {
				a1eiLog.Warn("A1EI SBI client incoming forwarder for EI Job Delete service is just closed")
				return
			}
			msg, err := a.sessions[stream.EIJobDelete].(a1.EIService_EIJobDeleteClient).Recv()
			if err == io.EOF || err == context.Canceled {
				a1eiLog.Warn("A1EI SBI client incoming forwarder for EI Job Delete service is just closed")
				return
			}
			if err != nil {
				a1eiLog.Warn(err)
				return
			}
			sbMessage := &stream.SBStreamMessage{
				TargetXAppID:     a.targetXAppID,
				A1SBIMessageType: stream.EIRequestMessage,
				A1Service:        stream.EnrichmentInformation,
				A1SBIRPCType:     stream.EIJobDelete,
				Payload:          msg,
			}
			nbID := stream.ID{
				SrcEndpointID:  stream.GetEndpointIDWithTargetXAppID(a.targetXAppID, stream.EnrichmentInformation),
				DestEndpointID: "a1p-controller",
			}
			err = a.streamBroker.Send(nbID, sbMessage)
			if err != nil {
				a1eiLog.Warn(err)
			}
		}
	}
}

func (a *a1eiClient) incomingEIJobStatusQueryForwarder(ctx context.Context) {
	defer a.Close()
	for {
		select {
		case <-ctx.Done():
			a1eiLog.Warn("A1EI SBI client incoming forwarder for EI Job Status Query service is just closed")
			return
		default:
			if _, ok := a.sessions[stream.EIJobStatusQuery]; !ok {
				a1eiLog.Warn("A1EI SBI client incoming forwarder for EI Job Status Query service is just closed")
				return
			}
			msg, err := a.sessions[stream.EIJobStatusQuery].(a1.EIService_EIJobStatusQueryClient).Recv()
			if err == io.EOF || err == context.Canceled {
				a1eiLog.Warn("A1EI SBI client incoming forwarder for EI Job Status Query service is just closed")
				return
			}
			if err != nil {
				a1eiLog.Warn(err)
				return
			}
			sbMessage := &stream.SBStreamMessage{
				TargetXAppID:     a.targetXAppID,
				A1SBIMessageType: stream.EIRequestMessage,
				A1Service:        stream.EnrichmentInformation,
				A1SBIRPCType:     stream.EIJobStatusQuery,
				Payload:          msg,
			}
			nbID := stream.ID{
				SrcEndpointID:  stream.GetEndpointIDWithTargetXAppID(a.targetXAppID, stream.EnrichmentInformation),
				DestEndpointID: "a1p-controller",
			}
			err = a.streamBroker.Send(nbID, sbMessage)
			if err != nil {
				a1eiLog.Warn(err)
			}
		}
	}
}

func (a *a1eiClient) runOutgoingMsgDispatcher(ctx context.Context) error {
	msgCh := make(chan *stream.SBStreamMessage)
	sbID := stream.ID{
		SrcEndpointID:  "a1ei-controller",
		DestEndpointID: stream.GetEndpointIDWithTargetXAppID(a.targetXAppID, stream.EnrichmentInformation),
	}

	watcherID := uuid.New()
	err := a.streamBroker.Watch(sbID, msgCh, watcherID)
	if err != nil {
		return err
	}

	go func(msgCh chan *stream.SBStreamMessage) {
		for msg := range msgCh {
			a.outgoingMsgDispatcher(ctx, msg)
		}
	}(msgCh)

	<-ctx.Done()
	return errors.NewCanceled("A1EI SBI client outgoing message dispatcher is just closed - due to the context done")
}

func (a *a1eiClient) outgoingMsgDispatcher(ctx context.Context, msg *stream.SBStreamMessage) {
	var err error
	switch msg.A1SBIRPCType {
	case stream.EIQuery:
		err = a.sessions[stream.EIQuery].(a1.EIService_EIQueryClient).Send(msg.Payload.(*a1.EIResultMessage))
		if err != nil {
			a1eiLog.Warn(err)
		}
	case stream.EIJobSetup:
		err = a.sessions[stream.EIJobSetup].(a1.EIService_EIJobSetupClient).Send(msg.Payload.(*a1.EIResultMessage))
		if err != nil {
			a1eiLog.Warn(err)
		}
	case stream.EIJobUpdate:
		err = a.sessions[stream.EIJobUpdate].(a1.EIService_EIJobUpdateClient).Send(msg.Payload.(*a1.EIResultMessage))
		if err != nil {
			a1eiLog.Warn(err)
		}
	case stream.EIJobDelete:
		err = a.sessions[stream.EIJobDelete].(a1.EIService_EIJobDeleteClient).Send(msg.Payload.(*a1.EIResultMessage))
		if err != nil {
			a1eiLog.Warn(err)
		}
	case stream.EIJobStatusQuery:
		err = a.sessions[stream.EIJobStatusQuery].(a1.EIService_EIJobStatusQueryClient).Send(msg.Payload.(*a1.EIResultMessage))
		if err != nil {
			a1eiLog.Warn(err)
		}
	case stream.EIJobStatusNotify:
		ack, err := a.grpcClient.EIJobStatusNotify(ctx, msg.Payload.(*a1.EIStatusMessage))
		if err != nil {
			a1eiLog.Warn(err)
		}
		a.forwardResponseMsg(ack, stream.EIAckMessage, stream.EIJobStatusNotify)
	case stream.EIJobResultDelivery:
		ack, err := a.grpcClient.EIJobResultDelivery(ctx, msg.Payload.(*a1.EIResultMessage))
		if err != nil {
			a1eiLog.Warn(err)
		}
		a.forwardResponseMsg(ack, stream.EIAckMessage, stream.EIJobResultDelivery)
	}
}

func (a *a1eiClient) forwardResponseMsg(msg interface{}, messageType stream.A1SBIMessageType, rpcType stream.A1SBIRPCType) {
	sbMessage := &stream.SBStreamMessage{
		TargetXAppID:     a.targetXAppID,
		A1SBIMessageType: messageType,
		A1Service:        stream.EnrichmentInformation,
		A1SBIRPCType:     rpcType,
		Payload:          msg,
	}
	nbID := stream.ID{
		SrcEndpointID:  stream.GetEndpointIDWithTargetXAppID(a.targetXAppID, stream.EnrichmentInformation),
		DestEndpointID: "a1p-controller",
	}
	err := a.streamBroker.Send(nbID, sbMessage)
	if err != nil {
		a1pLog.Warn(err)
	}
}

func (a *a1eiClient) Close() {
	defer delete(a.sessions, stream.EIQuery)
	defer delete(a.sessions, stream.EIJobSetup)
	defer delete(a.sessions, stream.EIJobUpdate)
	defer delete(a.sessions, stream.EIJobDelete)
	defer delete(a.sessions, stream.EIJobStatusQuery)

	// delete stream
	deleteStream(a.targetXAppID, stream.EnrichmentInformation, a.streamBroker)
}

var _ Client = &a1eiClient{}
