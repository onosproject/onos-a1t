# SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0

ARG ONOS_BUILD_VERSION=undefined

FROM onosproject/golang-build:v1.3.1 as build

ENV GO111MODULE=on
ARG ONOS_MAKE_TARGET=build

COPY Makefile go.mod go.sum /go/src/github.com/onosproject/onos-a1t/
COPY cmd/ /go/src/github.com/onosproject/onos-a1t/cmd/
COPY pkg/ /go/src/github.com/onosproject/onos-a1t/pkg/
COPY vendor/ /go/src/github.com/onosproject/onos-a1t/vendor/
COPY api/ /go/src/github.com/onosproject/onos-a1t/api
COPY build/ /go/src/github.com/onosproject/onos-a1t/build


RUN cd /go/src/github.com/onosproject/onos-a1t && GOFLAGS=-mod=vendor make ${ONOS_MAKE_TARGET}

FROM frolvlad/alpine-glibc:alpine-3.19

USER nobody

COPY --from=build /go/src/github.com/onosproject/onos-a1t/build/_output/onos-a1t /usr/local/bin/onos-a1t

ENTRYPOINT ["onos-a1t"]
