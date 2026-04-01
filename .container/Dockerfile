# Now copy it into our base image.

ARG TARGETOS
ARG TARGETARCH

FROM golang:1.18 as build

RUN yum update -y && yum install -y unzip curl \
    ### AWS CLI v2
    cli_arch=$(test "$TARGETARCH" = "amd64" && echo "x86_64" || echo "aarch64") \
    && curl "https://awscli.amazonaws.com/awscli-exe-linux-$cli_arch.zip" -o "aws-cli-v2.zip" \
    && unzip aws-cli-v2.zip \
    # The --bin-dir is specified so that we can copy the
    # entire bin directory from the installer stage into
    # into /usr/local/bin of the final stage without
    # accidentally copying over any other executables that
    # may be present in /usr/local/bin of the installer stage.
    && ./aws/install --bin-dir /aws-cli-bin/ \
    ### kubectl
    && curl -L "https://dl.k8s.io/release/$KUBECTL_VERSION/bin/linux/$TARGETARCH/kubectl" -o /usr/local/bin/kubectl \
    && chmod +x /usr/local/bin/kubectl \
    ### krew \
    set -x; cd "$(mktemp -d)" && \
    KREW="krew-${TARGETOS}_${TARGETARCH}" && \
    curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/latest/download/${KREW}.tar.gz" && \
    tar zxvf "${KREW}.tar.gz" && \
    ./"${KREW}" install krew
#
#
#
#CMD ["/app"]
#
#
#FROM --platform=$BUILDPLATFORM golang:1.24-bookworm AS build-stage
#
#ARG APPVERSION
#ARG BUILDNUMBER
#ARG KUBECTL_VERSION="v1.34.0"
#
#WORKDIR /app
#
#COPY . ./
#
#RUN go mod verify && go mod vendor
#
#RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
#    -ldflags "-X 'kubectl-plugin-arcane/cmd.Version=$APPVERSION'" \
#    -o /kubectl-arcane \
#    -v  \
#    .
#
#FROM --platform=$BUILDPLATFORM public.ecr.aws/amazonlinux/amazonlinux:2 AS aws-cli-stage
#
#ARG TARGETARCH
#ARG KUBECTL_VERSION="v1.34.0"
#
#RUN yum update -y && yum install -y unzip curl \
#    && cli_arch=$(test "$TARGETARCH" = "amd64" && echo "x86_64" || echo "aarch64") \
#    && curl "https://awscli.amazonaws.com/awscli-exe-linux-$cli_arch.zip" -o "awscliv2.zip" \
#    && unzip awscliv2.zip \
#  # The --bin-dir is specified so that we can copy the
#  # entire bin directory from the installer stage into
#  # into /usr/local/bin of the final stage without
#  # accidentally copying over any other executables that
#  # may be present in /usr/local/bin of the installer stage.
#    && ./aws/install --bin-dir /aws-cli-bin/ \
#    && curl -L "https://dl.k8s.io/release/$KUBECTL_VERSION/bin/linux/$TARGETARCH/kubectl" -o /usr/local/bin/kubectl \
#    && chmod +x /usr/local/bin/kubectl
#
#
## Deploy the application binary into a lean image
## We cannot use distroless here because we need to run /bin/sh -c in the entrypoint script
#FROM golang:1.24-bookworm AS build-release-stage
#
#WORKDIR /app
#
#RUN apt-get update \
# && apt-get --assume-yes upgrade \
# && apt-get --no-install-recommends --assume-yes install ca-certificates
#
#COPY --from=build-stage /snd /app
#
#COPY --from=aws-cli-stage /usr/local/aws-cli/ /usr/local/aws-cli/
#COPY --from=aws-cli-stage /aws-cli-bin/ /usr/local/bin/
#COPY --from=aws-cli-stage /usr/local/bin/kubectl /usr/local/bin/
#
#ENV PATH="$PATH:/app"
#
#ENTRYPOINT ["/app/snd"]