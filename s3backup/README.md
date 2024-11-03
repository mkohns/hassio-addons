![Supports aarch64 Architecture][aarch64-shield] ![Supports amd64 Architecture][amd64-shield] ![Supports armhf Architecture][armhf-shield] ![Supports armv7 Architecture][armv7-shield] ![Supports i386 Architecture][i386-shield]

[aarch64-shield]: https://img.shields.io/badge/aarch64-yes-green.svg
[amd64-shield]: https://img.shields.io/badge/amd64-yes-green.svg
[armhf-shield]: https://img.shields.io/badge/armhf-yes-green.svg
[armv7-shield]: https://img.shields.io/badge/armv7-yes-green.svg
[i386-shield]: https://img.shields.io/badge/i386-yes-green.svg

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/mkohns)

# Home Assistant Add-on: S3 Backup

## What is it for ?

This is my first addon! (reminder to myself: better use python next time).  
This addon is a [duplicity](https://duplicity.nongnu.org/vers7/duplicity.1.html) based tool to backup and restore the homeassistant /backup folder to an S3 bucket.  
I personaly use [IDrive](https://www.idrive.com/) as an S3 provider. But I also tested it with [Minio](https://min.io/)  

Key Features:

- always up-to-date duplicity version
- GPG Enryption support
- Backup support 
- Restore support
- MQTT based automatic virtual device creation
- Progress support
- Backup lifecycle support
- internal and external MQTT Server support
- fires events based on changes in /backup

![Device Support](https://raw.githubusercontent.com/mkohns/hassio-addons/refs/heads/main/s3backup/mqtt_device_support.png "MQTT Device Support")

## Disclaimer
This Addon only works 
- with GPG. Encryption is not optional. You will need a public and private key + passphrase. See Documentation.
- with an MQTT integration. This is not optional

## Installation

Follow these steps to get the add-on installed on your system:

0. Add this repository to your Addon Store: https://github.com/mkohns/hassio-addons
1. Navigate in your Home Assistant frontend to **Supervisor** -> **Add-on Store**.
2. Find the "S3Backup" add-on and click it.
3. Click on the "INSTALL" button.
