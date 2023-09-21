# CashCorns Backend

Build
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
```

Copy files to VM
```bash
scp  -i "~/.aws/cw-bb-test-zarf-agent-metrics.pem" cashcorns-backend ubuntu@ec2-18-117-82-51.us-east-2.compute.amazonaws.com:/home/ubuntu

scp  -i "~/.aws/cw-bb-test-zarf-agent-metrics.pem" tls.key ubuntu@ec2-18-117-82-51.us-east-2.compute.amazonaws.com:/home/ubuntu

scp  -i "~/.aws/cw-bb-test-zarf-agent-metrics.pem" tls.crt ubuntu@ec2-18-117-82-51.us-east-2.compute.amazonaws.com:/home/ubuntu
```

```bash
# Generate a private key
openssl genpkey -algorithm RSA -out tls.key

# Generate a self-signed certificate
openssl req -new -key tls.key -out tls.csr -subj "/CN=ec2-18-117-82-51.us-east-2.compute.amazonaws"

# Sign the certificate using the private key
openssl x509 -req -in tls.csr -signkey tls.key -out tls.crt
```

Run 
```bash
go run main.go serve -f ApprovalDates.json
```
