#!/usr/bin/with-contenv bashio

bashio::log.green "S3BACKUP"
bashio::log.green `date`

#####################
## USER PARAMETERS ##
#####################

# REQUIRED
# s3settings
export BUCKET_NAME="s3://$(bashio::config 's3settings_bucketName')"
export ENDPOINT_URL="$(bashio::config 's3settings_endpointUrl')"
export AWS_ACCESS_KEY_ID="$(bashio::config 's3settings_accessKey')"
export AWS_SECRET_ACCESS_KEY="$(bashio::config 's3settings_secretKey')"
# gpgsettings
export GPG_FINGERPRINT="$(bashio::config 'gpgsettings_fingerprint')"
export GPG_PASSPHRASE="$(bashio::config 'gpgsettings_passphrase')"
export PASSPHRASE="$(bashio::config 'gpgsettings_passphrase')"
export GPG_PUBLICKEYFILENAME="$(bashio::config 'gpgsettings_publicKey')"
export GPG_PRIVATEKEYFILENAME="$(bashio::config 'gpgsettings_privatKey')"
# duplicitysettings
export SOURCE_DIR="$(bashio::config 'duplicitysettings_sourceDir')"
export RESTORE_DIR="$(bashio::config 'duplicitysettings_restoreMediaDir')"
export INCREMENTAL_FOR="$(bashio::config 'duplicitysettings_incrementalFor')"
export REMOVE_OLDER_THAN="$(bashio::config 'duplicitysettings_removeOlderThan')"
# mqttsetting
export MQTT_USE_INTERNAL="$(bashio::config 'mqttsettings_use_internal')"
if [[ ${MQTT_USE_INTERNAL} == "true" ]]; then
    export MQTT_HOST="$(bashio::services mqtt "host")"
    export MQTT_USER="$(bashio::services mqtt "username")"
    export MQTT_PASSWORD="$(bashio::services mqtt "password")"
else
    export MQTT_HOST="$(bashio::config 'mqttsettings_server')"
    export MQTT_USER="$(bashio::config 'mqttsettings_user')"
    export MQTT_PASSWORD="$(bashio::config 'mqttsettings_pass')"
fi

# TOPIC CONFIG
export STATE_TOPIC="s3backup/state"
export COMMAND_TOPIC="s3backup/cmd"
export WILL_TOPIC="s3backup/status"

bashio::log.green "MQTT USER: $MQTT_USER"

# TRAF FUNC
send_availability() {
    bashio::log.green "Sending MQTT LastWill to: $WILL_TOPIC" 
    mosquitto_pub -d -h $MQTT_HOST -u $MQTT_USER -P $MQTT_PASSWORD -t "$WILL_TOPIC" -m "offline" -r
}
# Trap termination signals
trap send_availability SIGTERM SIGINT

# Start the device creation
. ./mqtt_discovery.sh

bashio::log.green "Showing duplicity version"
duplicity --version
bashio::log.green "Showing content of /ssl"
ls -al /ssl
bashio::log.green "Showing content persistent storage"
ls -al /data
bashio::log.green "Creating restore folder"
mkdir -p /media/$RESTORE_DIR
bashio::log.green "Showing content media"
ls -al /media
bashio::log.green "Importing Public Key"
gpg --import /ssl/$GPG_PUBLICKEYFILENAME
gpg --list-keys
bashio::log.green "Importing Private Key"
gpg --import --batch --passphrase "$GPG_PASSPHRASE" --pinentry-mode loopback  /ssl/$GPG_PRIVATEKEYFILENAME
gpg --list-keys
bashio::log.green "Trusting the key"
echo "${GPG_FINGERPRINT}:6:" | gpg --import-ownertrust
gpg --list-key

bashio::log.green "Starting the command topic listener"
bashio::log.green "Sending MQTT LastWill to: $WILL_TOPIC" 
mosquitto_pub -d -h $MQTT_HOST -u $MQTT_USER -P $MQTT_PASSWORD -t "$WILL_TOPIC" -m "online" -r

while read MESSAGE 
do 

    bashio::log.green "new mqtt message arrived: $MESSAGE"
    if [ "$MESSAGE" = "BACKUP" ]; then
        bashio::log.green "starting backup"
        rm -f backup.log
        duplicity --progress \
                  --name homeassistant \
                  --full-if-older-than $INCREMENTAL_FOR \
                  --skip-if-no-change \
                  --encrypt-key=$GPG_FINGERPRINT \
                  --s3-endpoint-url=$ENDPOINT_URL \
                  --log-file backup.log \
                  --archive-dir=/data \
                  backup $SOURCE_DIR $BUCKET_NAME &

        # Starting fake backup
        #bashio::log.green "starting fake backup"
        #./start_fake_backup.sh &

        # Starting the backup log analyser
        bashio::log.green "starting backup log parser"
        . ./backup_log_parser.sh

        # clean up the backup target
        duplicity --force \
                  --name homeassistant \
                  --encrypt-key=$GPG_FINGERPRINT \
                  --s3-endpoint-url=$ENDPOINT_URL \
                  --archive-dir=/data \
                 remove-older-than ${REMOVE_OLDER_THAN} "${BUCKET_NAME}"

        # cleanup
        bashio::log.green "cleanup ..."
        rm backup.log
    elif [ "$MESSAGE" = "RESTORE" ]; then
        bashio::log.red "Restoring backup to /media/$RESTORE_DIR"
        . ./restore.sh
    else
        bashio::log.red "message not handled"
    fi

done < <(mosquitto_sub -h $MQTT_HOST -u $MQTT_USER -P $MQTT_PASSWORD -t $COMMAND_TOPIC --will-retain --will-topic $WILL_TOPIC --will-payload offline)