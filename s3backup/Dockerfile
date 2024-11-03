ARG BUILD_FROM
FROM $BUILD_FROM

ENV LANG C.UTF-8
RUN apk add --no-cache duplicity py-boto3 mosquitto-clients inotify-tools

# Copy data for add-on
COPY *.sh /
RUN chmod a+x /run.sh
RUN chmod a+x /mqtt_discovery.sh
RUN chmod a+x /start_fake_backup.sh
RUN chmod a+x /backup_log_parser.sh
RUN chmod a+x /restore.sh
RUN chmod a+x /check_options.sh
RUN chmod a+x /watch_backup_folder.sh

CMD [ "/run.sh" ]