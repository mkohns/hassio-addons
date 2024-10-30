#!/usr/bin/with-contenv bashio

write2file() {
  echo "$@" >> backup.log
}

write2file "NOTICE 1"
write2file ". Local and Remote metadata are synchronized, no sync needed."
sleep 1

write2file "NOTICE 1"
write2file ". Last full backup date: none"
sleep 1

write2file "WARNING 1"
write2file ". No signatures found, switching to full backup."
sleep 1

write2file "NOTICE 16 0 3 0 3 0 0"
write2file ". 0.0KB 00:00:03 [0.0KB/s] [>                                        ] 0% ETA 3sec"
sleep 5

write2file "NOTICE 16 0 6 0 6 0 0"
write2file ". 0.0KB 00:00:06 [0.0KB/s] [>                                        ] 0% ETA 6sec"
sleep 5

write2file "NOTICE 16 0 9 0 6 0 1"
write2file ". 0.0KB 00:00:09 [0.0B/s] [>                                        ] 0% ETA Stalled!"
sleep 5

write2file "NOTICE 16 15224572 12 29 22 1522167 0"
write2file ". 14.5MB 00:00:12 [1.5MB/s] [===========>                             ] 29% ETA < 30sec"
sleep 5

write2file "NOTICE 16 22040316 15 42 16 1746928 0"
write2file ". 21.0MB 00:00:15 [1.7MB/s] [================>                        ] 42% ETA < 30sec"
sleep 5

write2file "NOTICE 16 23672837 18 100 0 1746928 0"
write2file ". 22.6MB 00:00:18 [1.7MB/s] [========================================>] 100% ETA 0sec"
sleep 5

write2file "NOTICE 1"
write2file ". --------------[ Backup Statistics ]--------------"
write2file ". StartTime 1730115176.38 (Mon Oct 28 11:32:56 2024)"
write2file ". EndTime 1730115186.12 (Mon Oct 28 11:33:06 2024)"
write2file ". ElapsedTime 9.75 (9.75 seconds)"
write2file ". SourceFiles 5031"
write2file ". SourceFileSize 52433212 (50.0 MB)"
write2file ". NewFiles 5031"
write2file ". NewFileSize 52433212 (50.0 MB)"
write2file ". DeletedFiles 0"
write2file ". ChangedFiles 0"
write2file ". ChangedFileSize 0 (0 bytes)"
write2file ". ChangedDeltaSize 0 (0 bytes)"
write2file ". DeltaEntries 5031"
write2file ". RawDeltaSize 47599899 (45.4 MB)"
write2file ". TotalDestinationSizeChange 22040316 (21.0 MB)"
write2file ". Errors 0"
write2file ". -------------------------------------------------"
write2file "."