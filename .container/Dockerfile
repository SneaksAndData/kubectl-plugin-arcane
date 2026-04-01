FROM --platform=$BUILDPLATFORM golang:1.25-bookworm AS build-stage

ARG TARGETOS
ARG TARGETARCH
ARG KUBECTL_VERSION="v1.34.0"

WORKDIR /app

COPY . ./

RUN go mod verify && go mod vendor

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /kubectl-arcane -v  .


FROM --platform=$BUILDPLATFORM public.ecr.aws/amazonlinux/amazonlinux:2 AS aws-cli-stage

ARG TARGETOS
ARG TARGETARCH
ARG KUBECTL_VERSION="v1.34.0"

RUN yum update -y && yum install -y unzip curl tar gzip \
    ###\
    ### AWS CLI v2 \
    ###\
    && cli_arch=$(test "$TARGETARCH" = "amd64" && echo "x86_64" || echo "aarch64") \
    && curl "https://awscli.amazonaws.com/awscli-exe-linux-$cli_arch.zip" -o "awscliv2.zip" \
    && unzip awscliv2.zip \
    # The --bin-dir is specified so that we can copy the \
    # entire bin directory from the installer stage into \
    # into /usr/local/bin of the final stage without \
    # accidentally copying over any other executables that \
    # may be present in /usr/local/bin of the installer stage. \
    && ./aws/install --bin-dir /aws-cli-bin/ \
    ###\
    ### kubectl \
    ###\
    && curl -L "https://dl.k8s.io/release/$KUBECTL_VERSION/bin/linux/$TARGETARCH/kubectl" -o /usr/local/bin/kubectl \
    && chmod +x /usr/local/bin/kubectl


FROM gcr.io/distroless/static-debian12
COPY --from=aws-cli-stage /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --from=build-stage /kubectl-arcane /usr/local/bin
COPY --from=aws-cli-stage /usr/local/aws-cli/ /usr/local/aws-cli/
COPY --from=aws-cli-stage /aws-cli-bin/ /usr/local/bin/
ENTRYPOINT ["kubectl"]
