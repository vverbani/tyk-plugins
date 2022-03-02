#!/bin/bash

rm -rf vendor
pip3 install -r requirements.txt --prefix vendor
/opt/tyk-gateway/tyk bundle build -m manifest.json -y -o vault.zip
echo Press enter to add the vendor directory to the bundle
read ready
zip -ur ticket.zip vendor/
