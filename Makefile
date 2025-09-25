DEBUG_FLAGS := DEBUG=true

EXECUTABLE := ./build/zarinit-server

.PHONY: build install enable-service disable-service update-service restart test

build: 
	go build -o $(EXECUTABLE) ./cmd/server/main.go
	chmod 775 $(EXECUTABLE)

test:
	go test ./...
	
#=== SERVICE
SERVICE_NAME := zarinit-server.service


install: build
	
	mkdir -p /opt/zarinit/
	cp $(EXECUTABLE) /opt/zarinit/zarinit-server
	
	# install services
	cp ./services/$(SERVICE_NAME) /lib/systemd/system/
	cp ./services/z-hostapd-2.service /lib/systemd/system/
	cp ./services/z-hostapd-5.service /lib/systemd/system/
	
	# install configs
	mkdir -p /etc/zarinit/
	cp ./cloud-config.yml /etc/zarinit
	cp ./router-config.yml /etc/zarinit

enable-service: install
	systemctl enable --now $(SERVICE_NAME) || journalctl -xeu $(SERVICE_NAME)

disable-service:
	systemctl disable --now $(SERVICE_NAME)

update-service: disable-service install enable-service
