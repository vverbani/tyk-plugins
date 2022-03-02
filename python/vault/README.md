# vault python middleware example.

hvac has some issues with the six module. You need to import six from your middleware or it won't load when calling it via the python embedded interpreter

## Setup vault
```
$ docker container create --cap-add=IPC_LOCK -e 'VAULT_DEV_ROOT_TOKEN_ID=myroot' -e 'VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:1234' -p 8200:1234 --name vault vault
$ docker start vault
$ docker container exec -it vault /bin/sh
# export VAULT_ADDR="http://localhost:1234"
# vault login myroot
# export VAULT_TOKEN=myroot
# vault kv put secret/data fred=jim
# vault kv get secret/data
```

## create the bundle
This needs to be done inside a tyk docker image unless you install tyk locally
```
rm -rf vendor
pip3 install -r requirements.txt --prefix vendor
/opt/tyk-gateway/tyk bundle build -m manifest.json -y -o vault.zip
echo Press enter to add the vendor directory to the bundle
read ready
zip -ur bundle.zip vendor/
```

## load the bundle and use it
