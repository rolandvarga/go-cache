#!/bin/zsh

payload='{
	"id": "samplestring-85c86c465-12345/9/1586875121821105370/65598300_146",
}'

post_endpoint="http://localhost:8080/new"
entries_endpoint="http://localhost:8080/entries"


echo "sending payload to endpoint ${post_endpoint}"
curl -L -vv --silent -i --location --request POST $post_endpoint \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--data-raw $payload

curl -L -vv -XGET $entries_endpoint \
--header 'Accept: application/json'
