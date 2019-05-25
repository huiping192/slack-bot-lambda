env GOOS=linux GOARCH=amd64 go build -o /tmp/main
zip -j /tmp/main.zip /tmp/main
aws lambda update-function-code --function-name slack-events-endpoint \
--zip-file fileb:///tmp/main.zip
