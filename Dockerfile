FROM golang:{{.Inputs.go_version.Value}} AS build
WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -trimpath \
    -ldflags "-X github.com/{{.Inputs.owner.Value}}/{{.Inputs.project_name.Value}}/internal/version.Version=${VERSION} -s -w" \
    -o /out/{{.Inputs.project_name.Value}} ./cmd/{{.Inputs.project_name.Value}}

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/{{.Inputs.project_name.Value}} /{{.Inputs.project_name.Value}}

EXPOSE 8080
ENTRYPOINT ["/{{.Inputs.project_name.Value}}"]
