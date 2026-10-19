FROM golang:{{.Inputs.GoVersion.Value}} AS build
WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -trimpath \
    -ldflags "-X github.com/{{.Inputs.Owner.Value}}/{{.Inputs.ProjectName.Value}}/internal/version.Version=${VERSION} -s -w" \
    -o /out/{{.Inputs.ProjectName.Value}} ./cmd/{{.Inputs.ProjectName.Value}}

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/{{.Inputs.ProjectName.Value}} /{{.Inputs.ProjectName.Value}}

EXPOSE 8080
ENTRYPOINT ["/{{.Inputs.ProjectName.Value}}"]
