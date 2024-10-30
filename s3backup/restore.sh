#!/usr/bin/with-contenv bashio

export state_active="ON"
export state_progress=0
export state_uploaded_size=0
export state_start_time=$(date -u  +"%Y-%m-%dT%H:%M:%SZ")
export state_end_time="unknown"
export state_elapsed_time=0
export state_source_files=0
export state_source_files_size=0
export state_new_files=0
export state_new_files_size=0
export state_deleted_files=0
export state_changed_files=0
export state_delta_entries=0
export state_errors=0


sendState() {
    STATE=$(cat <<EOF
{
    "active": "$state_active",
    "progress": $state_progress,
    "uploaded_size": $state_uploaded_size,
    "start_time": "$state_start_time",
    "end_time": "$state_end_time",
    "elapsed_time": $state_elapsed_time,
    "source_files": $state_source_files,
    "source_files_size": $state_source_files_size,
    "new_files": $state_new_files,
    "new_files_size": $state_new_files_size,
    "deleted_files": $state_deleted_files,
    "changed_files": $state_changed_files,
    "delta_entries": $state_delta_entries,
    "errors": $state_errors
}
EOF
)
    mosquitto_pub -d -h $MQTT_HOST -u $MQTT_USER -P $MQTT_PASSWORD -t "$STATE_TOPIC" -m "$STATE"
}

sendState

bashio::log.green "starting restore"

duplicity -v5 \
            --name homeassistant \
            --encrypt-key=$GPG_FINGERPRINT \
            --s3-endpoint-url=$ENDPOINT_URL \
            --archive-dir=/data \
            restore $BUCKET_NAME /media/$RESTORE_DIR &

if [ $? -eq 0 ]; then
    bashio::log.green "restore finished"
    state_active="OFF"
    state_progress=100
    state_end_time=$(date -u  +"%Y-%m-%dT%H:%M:%SZ")
    sendState
else
    bashio::log.green "restore failed"
    state_end_time=$(date -u  +"%Y-%m-%dT%H:%M:%SZ")
    state_active="OFF"
    state_progress=100
    state_errors=1
    sendState
fi

