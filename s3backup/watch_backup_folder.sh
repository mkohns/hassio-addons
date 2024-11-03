#!/usr/bin/with-contenv bashio

# Directory to monitor (set this to your desired folder)
MONITOR_DIR="/backup"

# Configurable wait time for file stability check (in seconds)
STABILITY_WAIT_TIME=5

# Home Assistant API token for authentication
HA_TOKEN="${SUPERVISOR_TOKEN}"  

# Home Assistant events API
HA_API_URL="http://supervisor/core/api/events" 

send_event() {
    local event_type="$1"
    local file="$2"
    bashio::log.green [`date`] send event $event_type / $file
    curl -s -X POST -H "Authorization: Bearer $HA_TOKEN" \
         -H "Content-Type: application/json" \
         -d "{\"filename\": \"$file\"}" \
         "$HA_API_URL/$event_type"
}

send_event "backup_watcher_started" "now"

# Function to check file stability
check_file_stability() {
    local file="$1"
    local last_size=-1

    while true; do
        # Check file size
        current_size=$(stat -c %s "$file" 2>/dev/null)
        
        # If file no longer exists, exit the loop
        if [ ! -f "$file" ]; then
            return
        fi

        if [ "$current_size" -eq "$last_size" ]; then
            bashio::log.green [`date`] file completed: $file
            send_event "backup_ended" "$file"
            return
        fi

        last_size=$current_size
        sleep "$STABILITY_WAIT_TIME"
    done
}

# Monitor folder for newly created or deleted files
inotifywait -m -e create -e delete "$MONITOR_DIR" --format '%e %w%f' | while read event file
do
    if [[ "$event" == "CREATE" ]]; then
        bashio::log.green [`date`] new file created: $file
        send_event "backup_started" "$file"        
        # Check file stability in the background
        check_file_stability "$file" &
    elif [[ "$event" == "DELETE" ]]; then
        bashio::log.green [`date`] file removed: $file
        send_event "backup_file_deleted" "$file"
    fi
done