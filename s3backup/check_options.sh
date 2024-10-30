#!/usr/bin/with-contenv bashio

if [ -z "${BUCKET_NAME}" ]; then
    bashio::log.green s3settings_bucketName must be set
    exit
fi

if [ -z "${ENDPOINT_URL}" ]; then
    bashio::log.green s3settings_endpointUrl must be set
    exit
fi

if [ -z "${AWS_ACCESS_KEY_ID}" ]; then
    bashio::log.green s3settings_accessKey must be set
    exit
fi

if [ -z "${AWS_SECRET_ACCESS_KEY}" ]; then
    bashio::log.green s3settings_secretKey must be set
    exit
fi

if [ -z "${GPG_FINGERPRINT}" ]; then
    bashio::log.green gpgsettings_fingerprint must be set
    exit
fi

if [ -z "${GPG_PASSPHRASE}" ]; then
    bashio::log.green gpgsettings_passphrase must be set
    exit
fi

if [ -z "${GPG_PUBLICKEYFILENAME}" ]; then
    bashio::log.green gpgsettings_publicKey must be set
    exit
fi

if [ -z "${GPG_PRIVATEKEYFILENAME}" ]; then
    bashio::log.green gpgsettings_privatKey must be set
    exit
fi

if [ -z "${SOURCE_DIR}" ]; then
    bashio::log.green duplicitysettings_sourceDir must be set
    exit
fi

if [ -z "${RESTORE_DIR}" ]; then
    bashio::log.green duplicitysettings_restoreMediaDir must be set
    exit
fi

if [ -z "${INCREMENTAL_FOR}" ]; then
    bashio::log.green duplicitysettings_incrementalFor must be set
    exit
fi

if [ -z "${REMOVE_OLDER_THAN}" ]; then
    bashio::log.green duplicitysettings_removeOlderThan must be set
    exit
fi

if [ -z "${MQTT_HOST}" ]; then
    bashio::log.green could not find mqtt server name
    exit
fi

if [ -z "${MQTT_USER}" ]; then
    bashio::log.green could not find mqtt user name
    exit
fi

if [ -z "${MQTT_PASSWORD}" ]; then
    bashio::log.green could not find mqtt user password
    exit
fi
