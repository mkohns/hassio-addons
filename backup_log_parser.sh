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

bashio::log.green "wait for backup.log"
# Warte auf das Vorhandensein der Datei
while [ ! -f backup.log ]; do
    bashio::log.green "wait for backup.log"
    sleep 1
done

bashio::log.green "backup.log available"

bashio::log.green "starting backup.log parser"

# Initialize line counter
line_count=0
export stop=0

# Continuously monitor the file
while true; do
    # Get the current total line count of the file
    new_line_count=$(wc -l < "backup.log")

    # Check if there are new lines to process
    if (( new_line_count > line_count )); then
        # Process only the new lines
        while IFS= read -r line; do
            array=($line)
            # Stop reading if a line contains "."
            if [[ "$line" == ". " ]]; then
                bashio::log.green "stop processing"
                sendState
                stop=1
            elif [[ "$line" == "NOTICE 16"* ]]; then
                bashio::log.green "detected progress: $line"
                state_progress=${array[4]}
                state_uploaded_size=${array[2]}
                state_elapsed_time=${array[3]}
                sendState
            elif [[ "$line" == *"StartTime"* ]]; then
                bashio::log.green "detected StartTime: $line"
                time=${array[2]}
                time_int=${time%.*}
                state_start_time=$(date -u -d @$time_int +"%Y-%m-%dT%H:%M:%SZ")
            elif [[ "$line" == *"EndTime"* ]]; then
                bashio::log.green "detected EndTime: $line"
                time=${array[2]}
                time_int=${time%.*}
                state_end_time=$(date -u -d @$time_int +"%Y-%m-%dT%H:%M:%SZ")
            elif [[ "$line" == *"ElapsedTime"* ]]; then
                bashio::log.green "detected ElapsedTime: $line"
                time=${array[2]}
                time_int=${time%.*}
                state_elapsed_time=$(($state_elapsed_time + $time_int))
            elif [[ "$line" == *"SourceFiles"* ]]; then
                bashio::log.green "detected SourceFiles: $line"
                state_source_files=${array[2]}
            elif [[ "$line" == *"SourceFileSize"* ]]; then
                bashio::log.green "detected SourceFileSize: $line"
                state_source_files_size=${array[2]}
            elif [[ "$line" == *"NewFiles"* ]]; then
                bashio::log.green "detected NewFiles: $line"
                state_new_files=${array[2]}
            elif [[ "$line" == *"NewFileSize"* ]]; then
                bashio::log.green "detected NewFileSize: $line"
                state_new_files_size=${array[2]}
            elif [[ "$line" == *"DeletedFiles"* ]]; then
                bashio::log.green "detected DeletedFiles: $line"
                state_deleted_files=${array[2]}
            elif [[ "$line" == *"ChangedFiles"* ]]; then
                bashio::log.green "detected ChangedFiles: $line"
                state_changed_files=${array[2]} 
            elif [[ "$line" == *"DeltaEntries"* ]]; then
                bashio::log.green "detected DeltaEntries: $line"
                state_delta_entries=${array[2]}
            elif [[ "$line" == *"Errors"* ]]; then
                bashio::log.green "detected Errors: $line"
                state_errors=${array[2]}                
            fi              

            # Print each line (optional)
            bashio::log.yellow "$line"
        done < <(tail -n $((new_line_count - line_count)) "backup.log")  # Process tail output

        if (( $stop > 0 )); then
            bashio::log.green "sending break"
            break
        fi

        # Update the line counter to the new line count
        line_count=$new_line_count
    fi

    # Pause briefly before checking again to avoid high CPU usage
    sleep 1
done

bashio::log.green "stopping backup.log parser"

state_active="OFF"
sendState