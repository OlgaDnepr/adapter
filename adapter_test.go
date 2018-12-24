package main

import (
	"context"

	"github.com/OlgaDnepr/adapter/mocks"
	"github.com/OlgaDnepr/adapter/pb"
	"github.com/golang/mock/gomock"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = ginkgo.Describe("Adapter tests", func() {
	mockController := gomock.NewController(ginkgo.GinkgoT())
	mockServer := mocks.NewMockServerClient(mockController)
	adapter := newAdapterServer(mockServer)
	monkeyReply := &pb.Reply{Message: pb.MonkeyFollow_Monkey}
	followReply := &pb.Reply{Message: pb.MonkeyFollow_Follow}
	invalidReply := &pb.Reply{Message: pb.MonkeyFollow(-1)}

	ginkgo.It("Marco success", func() {
		mockServer.EXPECT().Get(gomock.Any(), monkeyReply).Return(followReply, nil)
		request := &pb.Request{Message: pb.MarcoPolo_Marco}
		out, err := adapter.Get(context.Background(), request)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(out.Message).To(gomega.Equal(pb.MarcoPolo_Polo))
	})

	ginkgo.It("Polo success", func() {
		mockServer.EXPECT().Get(gomock.Any(), followReply).Return(monkeyReply, nil)
		request := &pb.Request{Message: pb.MarcoPolo_Polo}
		out, err := adapter.Get(context.Background(), request)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(out.Message).To(gomega.Equal(pb.MarcoPolo_Marco))
	})

	ginkgo.It("Error: invalid message", func() {
		request := &pb.Request{Message: pb.MarcoPolo(-1)}
		out, err := adapter.Get(context.Background(), request)
		gomega.Expect(err).ToNot(gomega.BeNil())
		gomega.Expect(out).To(gomega.BeNil())
	})

	ginkgo.It("Error: no request", func() {
		out, err := adapter.Get(context.Background(), nil)
		gomega.Expect(err).ToNot(gomega.BeNil())
		gomega.Expect(out).To(gomega.BeNil())
	})

	ginkgo.It("Error in server communication", func() {
		mockServer.EXPECT().Get(gomock.Any(), monkeyReply).Return(nil, errors.New("invalid reply"))
		request := &pb.Request{Message: pb.MarcoPolo_Marco}
		out, err := adapter.Get(context.Background(), request)
		gomega.Expect(err).ToNot(gomega.BeNil())
		gomega.Expect(out).To(gomega.BeNil())
	})

	ginkgo.It("Error invalid reply", func() {
		mockServer.EXPECT().Get(gomock.Any(), monkeyReply).Return(invalidReply, nil)
		request := &pb.Request{Message: pb.MarcoPolo_Marco}
		out, err := adapter.Get(context.Background(), request)
		gomega.Expect(err).ToNot(gomega.BeNil())
		gomega.Expect(out).To(gomega.BeNil())
	})

	ginkgo.It("Error no reply", func() {
		mockServer.EXPECT().Get(gomock.Any(), followReply).Return(nil, nil)
		request := &pb.Request{Message: pb.MarcoPolo_Polo}
		out, err := adapter.Get(context.Background(), request)
		gomega.Expect(err).ToNot(gomega.BeNil())
		gomega.Expect(out).To(gomega.BeNil())
	})
})
