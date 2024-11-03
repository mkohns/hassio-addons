#!/usr/bin/with-contenv bashio

bashio::log.green "Sending MQTT Discovery"
date

#####################
## USER PARAMETERS ##
#####################

sendDiscovery () {
bashio::log.green "Sending MQTT Discovery to: $MQTT_DISCOVERY_TOPIC" 
mosquitto_pub -r -d -h $MQTT_HOST -u $MQTT_USER -P $MQTT_PASSWORD -t "$MQTT_DISCOVERY_TOPIC" -m "$MQTT_DISCOVERY_MSG"
}

# ACTIVE ENTITY
MQTT_DISCOVERY_MSG=$(cat <<EOF
{
   "device_class":"running",
   "name":"Status",
   "availability_topic":"$WILL_TOPIC",
   "state_topic":"$STATE_TOPIC",
   "value_template":"{{ value_json.active }}",
   "unique_id":"s3backup_active",
   "device":{
      "identifiers":[
          "s3backup"
      ],
      "name":"S3Backup",
      "manufacturer": "mkohns",
      "model": "Virtual Sensor",
      "model_id": "S3"
   }
}
EOF
)
MQTT_DISCOVERY_TOPIC='homeassistant/binary_sensor/s3backup/active/config'
sendDiscovery

# PROGRESS ENTITY
MQTT_DISCOVERY_MSG=$(cat <<EOF
{
   "device_class":"power_factor",
   "name": "Progress",
   "unit_of_measurement": "%",
   "state_topic":"$STATE_TOPIC",
   "value_template":"{{ value_json.progress }}",
   "unique_id":"s3backup_progress",
   "availability_topic":"$WILL_TOPIC",
   "device":{
      "identifiers":[
          "s3backup"
      ]
   }
}
EOF
)
MQTT_DISCOVERY_TOPIC='homeassistant/sensor/s3backup/progress/config'
sendDiscovery

# UPLOADED SIZE ENTITY
MQTT_DISCOVERY_MSG=$(cat <<EOF
{
   "device_class":"data_size",
   "name":"Uploaded Bytes",
   "unit_of_measurement": "B",
   "state_topic":"$STATE_TOPIC",
   "value_template":"{{ value_json.uploaded_size }}",
   "unique_id":"s3backup_uploadedsize",
   "availability_topic":"$WILL_TOPIC",
   "device":{
      "identifiers":[
          "s3backup"
      ]
   }
}
EOF
)
MQTT_DISCOVERY_TOPIC='homeassistant/sensor/s3backup/uploadedsize/config'
sendDiscovery

# START TIME ENTITY
MQTT_DISCOVERY_MSG=$(cat <<EOF
{
   "device_class":"timestamp",
   "name":"Start Time",
   "state_topic":"$STATE_TOPIC",
   "value_template":"{{ value_json.start_time }}",
   "unique_id":"s3backup_start_time",
   "availability_topic":"$WILL_TOPIC",
   "device":{
      "identifiers":[
          "s3backup"
      ]
   }
}
EOF
)
MQTT_DISCOVERY_TOPIC='homeassistant/sensor/s3backup/starttime/config'
sendDiscovery

# END TIME ENTITY
MQTT_DISCOVERY_MSG=$(cat <<EOF
{
   "device_class":"timestamp",
   "name":"End Time",
   "state_topic":"$STATE_TOPIC",
   "value_template":"{{ value_json.end_time }}",
   "unique_id":"s3backup_end_time",
   "availability_topic":"$WILL_TOPIC",
   "device":{
      "identifiers":[
          "s3backup"
      ]
   }
}
EOF
)
MQTT_DISCOVERY_TOPIC='homeassistant/sensor/s3backup/endtime/config'
sendDiscovery

# ELAPSED TIME ENTITY
MQTT_DISCOVERY_MSG=$(cat <<EOF
{
   "device_class":"duration",
   "name":"Elapsed Time",
   "unit_of_measurement": "s",
   "state_topic":"$STATE_TOPIC",
   "value_template":"{{ value_json.elapsed_time }}",
   "unique_id":"s3backup_elapsed_time",
   "availability_topic":"$WILL_TOPIC",
   "device":{
      "identifiers":[
          "s3backup"
      ]
   }
}
EOF
)
MQTT_DISCOVERY_TOPIC='homeassistant/sensor/s3backup/elapsedtime/config'
sendDiscovery

# SOURCE FILES ENTITY
MQTT_DISCOVERY_MSG=$(cat <<EOF
{
   "name":"Source Files",
   "state_topic":"$STATE_TOPIC",
   "value_template":"{{ value_json.source_files }}",
   "unique_id":"s3backup_source_files",
   "availability_topic":"$WILL_TOPIC",
   "device":{
      "identifiers":[
          "s3backup"
      ]
   }
}
EOF
)
MQTT_DISCOVERY_TOPIC='homeassistant/sensor/s3backup/sourcefiles/config'
sendDiscovery

# SOURCE FILES SIZE ENTITY
MQTT_DISCOVERY_MSG=$(cat <<EOF
{
   "name":"Source Files Size",
   "state_topic":"$STATE_TOPIC",
   "device_class":"data_size",
   "unit_of_measurement": "B",   
   "value_template":"{{ value_json.source_files_size }}",
   "unique_id":"s3backup_source_files_size",
   "availability_topic":"$WILL_TOPIC",
   "device":{
      "identifiers":[
          "s3backup"
      ]
   }
}
EOF
)
MQTT_DISCOVERY_TOPIC='homeassistant/sensor/s3backup/sourcefilessize/config'
sendDiscovery

# NEW FILES ENTITY
MQTT_DISCOVERY_MSG=$(cat <<EOF
{
   "name":"New Files",
   "state_topic":"$STATE_TOPIC",
   "value_template":"{{ value_json.new_files }}",
   "unique_id":"s3backup_new_files",
   "availability_topic":"$WILL_TOPIC",
   "device":{
      "identifiers":[
          "s3backup"
      ]
   }
}
EOF
)
MQTT_DISCOVERY_TOPIC='homeassistant/sensor/s3backup/newfiles/config'
sendDiscovery

# NEW FILES SIZE ENTITY
MQTT_DISCOVERY_MSG=$(cat <<EOF
{
   "name":"New Files Size",
   "state_topic":"$STATE_TOPIC",
   "device_class":"data_size",
   "unit_of_measurement": "B",   
   "value_template":"{{ value_json.new_files_size }}",
   "unique_id":"s3backup_new_files_size",
   "availability_topic":"$WILL_TOPIC",
   "device":{
      "identifiers":[
          "s3backup"
      ]
   }
}
EOF
)
MQTT_DISCOVERY_TOPIC='homeassistant/sensor/s3backup/newfilessize/config'
sendDiscovery

# DELETED FILES ENTITY
MQTT_DISCOVERY_MSG=$(cat <<EOF
{
   "name":"Deleted Files",
   "state_topic":"$STATE_TOPIC",
   "value_template":"{{ value_json.deleted_files }}",
   "unique_id":"s3backup_deleted_files",
   "availability_topic":"$WILL_TOPIC",
   "device":{
      "identifiers":[
          "s3backup"
      ]
   }
}
EOF
)
MQTT_DISCOVERY_TOPIC='homeassistant/sensor/s3backup/deletedfiles/config'
sendDiscovery

# DELETED FILES ENTITY
MQTT_DISCOVERY_MSG=$(cat <<EOF
{
   "name":"Changed Files",
   "state_topic":"$STATE_TOPIC",
   "value_template":"{{ value_json.changed_files }}",
   "unique_id":"s3backup_changed_files",
   "availability_topic":"$WILL_TOPIC",
   "device":{
      "identifiers":[
          "s3backup"
      ]
   }
}
EOF
)
MQTT_DISCOVERY_TOPIC='homeassistant/sensor/s3backup/changedfiles/config'
sendDiscovery

# DELTA ENTRIES FILES ENTITY
MQTT_DISCOVERY_MSG=$(cat <<EOF
{
   "name":"Delta Entries",
   "state_topic":"$STATE_TOPIC",
   "value_template":"{{ value_json.delta_entries }}",
   "unique_id":"s3backup_delta_entries",
   "availability_topic":"$WILL_TOPIC",
   "device":{
      "identifiers":[
          "s3backup"
      ]
   }
}
EOF
)
MQTT_DISCOVERY_TOPIC='homeassistant/sensor/s3backup/deltaentries/config'
sendDiscovery

# ERROR ENTITY
MQTT_DISCOVERY_MSG=$(cat <<EOF
{
   "name":"Errors",
   "state_topic":"$STATE_TOPIC",
   "value_template":"{{ value_json.errors }}",
   "unique_id":"s3backup_errors",
   "availability_topic":"$WILL_TOPIC",
   "device":{
      "identifiers":[
          "s3backup"
      ]
   }
}
EOF
)
MQTT_DISCOVERY_TOPIC='homeassistant/sensor/s3backup/errors/config'
sendDiscovery

# BACKUP BUTTON ENTITY
MQTT_DISCOVERY_MSG=$(cat <<EOF
{
   "name":"Backup",
   "state_topic":"$STATE_TOPIC",
   "unique_id":"s3backup_start_backup",
   "command_topic": "$COMMAND_TOPIC",
   "availability_topic":"$WILL_TOPIC",
   "payload_press":"BACKUP",   
   "device":{
      "identifiers":[
          "s3backup"
      ]
   }
}
EOF
)
MQTT_DISCOVERY_TOPIC='homeassistant/button/s3backup/startbackup/config'
sendDiscovery

# RESTORE BUTTON ENTITY
MQTT_DISCOVERY_MSG=$(cat <<EOF
{
   "name":"Restore",
   "state_topic":"$STATE_TOPIC",
   "unique_id":"s3backup_start_restore",
   "command_topic": "$COMMAND_TOPIC",
   "availability_topic":"$WILL_TOPIC",
   "payload_press":"RESTORE",   
   "device":{
      "identifiers":[
          "s3backup"
      ]
   }
}
EOF
)
MQTT_DISCOVERY_TOPIC='homeassistant/button/s3backup/startrestore/config'
sendDiscovery

# wait a few seconds before sending first state 
#sleep 5

#STATE=$(cat <<EOF
#{
#    "active": "OFF",
#    "progress": 0,
#    "uploaded_size": 0,
#    "start_time": null,
#    "end_time": null,
#    "elapsed_time": 0,
#    "source_files": 0,
#    "source_files_size": 0,
#    "new_files": 0,
#    "new_files_size": 0,
#    "deleted_files": 0,
#    "changed_files": 0,
#    "delta_entries": 0,
#    "errors": 0
#}
#EOF
#)
#mosquitto_pub -d -h $MQTT_HOST -u $MQTT_USER -P $MQTT_PASSWORD -t "$STATE_TOPIC" -m "$STATE"