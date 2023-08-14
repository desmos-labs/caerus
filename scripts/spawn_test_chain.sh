#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
DESMOS_HOME="$SCRIPT_DIR/.desmos"
KEYRING_PASS=pass1234

# APIs account information
APIS_ACCOUNT=bondscape-apis
APIS_ACCOUNT_ADDRESS=desmos1nm6kh6jwqmsezwtnmgdd4w4tzyk9f8gvqu5en0
APIS_ACCOUNT_MNEMONIC="hour harbor fame unaware bunker junk garment decrease federal vicious island smile warrior fame right suit portion skate analyst multiply magnet medal fresh sweet"

# User information
USER_ACCOUNT=bondscape-user
USER_ACCOUNT_ADDRESS=desmos1c7ms9zhtgwmv5jy6ztj2vq0jj67zenw3gdl2gr
USER_ACCOUNT_MNEMONIC="shine velvet envelope baby bicycle slot betray water mask cake rule useless absorb lens enable fork lake make depart slogan kiwi iron decorate relief"

# Background flags, will be set to true if the user pass the -b argument.
BACKGROUND=false

while getopts "b" arg; do
  case $arg in
    b)
      BACKGROUND=true
      ;;
  esac
done

desmos() {
	"$SCRIPT_DIR/desmos" --home="$DESMOS_HOME" "$@"
}

rm -r -f "$DESMOS_HOME"
desmos tendermint unsafe-reset-all
desmos init testchain --chain-id=testchain
# Add a default reason to the reports module params
jq '.app_state.reports.params.standard_reasons[0] |= . + {"id":"1","title":"Spam","description":"Spam user or content"}' "$DESMOS_HOME/config/genesis.json" > "$DESMOS_HOME/config/genesis-patched.json"
mv "$DESMOS_HOME/config/genesis-patched.json" "$DESMOS_HOME/config/genesis.json"

(echo "$APIS_ACCOUNT_MNEMONIC"; echo $KEYRING_PASS; echo $KEYRING_PASS) | desmos keys add "$APIS_ACCOUNT" --recover --keyring-backend=file
(echo "$USER_ACCOUNT_MNEMONIC"; echo $KEYRING_PASS; echo $KEYRING_PASS) | desmos keys add "$USER_ACCOUNT" --recover --keyring-backend=file
echo $KEYRING_PASS | desmos add-genesis-account $APIS_ACCOUNT 200000000000000stake --keyring-backend=file
echo $KEYRING_PASS | desmos add-genesis-account $USER_ACCOUNT 200000000000000stake --keyring-backend=file
echo $KEYRING_PASS | desmos gentx $APIS_ACCOUNT 100000000000stake --amount 100000000000stake --chain-id=testchain --keyring-backend=file
desmos collect-gentxs


if [ $BACKGROUND = true ] ; then
  desmos start &> "$SCRIPT_DIR/desmos.log" &
else
  desmos start
fi

