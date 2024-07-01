# SPDX-FileCopyrightText: 2022 2020-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0

ARG ONOS_BUILD_VERSION=undefined

FROM onosproject/golang-build:v1.0 as build

ENV GO111MODULE=on
ARG ONOS_MAKE_TARGET=build

COPY Makefile go.mod go.sum /go/src/github.com/onosproject/onos-a1txapp/
COPY cmd/ /go/src/github.com/onosproject/onos-a1txapp/cmd/
COPY pkg/ /go/src/github.com/onosproject/onos-a1txapp/pkg/
COPY vendor/ /go/src/github.com/onosproject/onos-a1txapp/vendor/
COPY build/ /go/src/github.com/onosproject/onos-a1txapp/build

RUN cd /go/src/github.com/onosproject/onos-a1txapp && GOFLAGS=-mod=vendor make ${ONOS_MAKE_TARGET}

FROM alpine:3.20
RUN apk add libc6-compat

USER nobody

COPY --from=build /go/src/github.com/onosproject/onos-a1txapp/build/_output/onos-a1txapp /usr/local/bin/onos-a1txapp

ENTRYPOINT ["onos-a1txapp"]
