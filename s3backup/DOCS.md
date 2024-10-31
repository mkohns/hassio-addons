[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/mkohns)

# Home Assistant Add-on: S3 Backup

## Prerequisites

- Make sure you have the [internal mqtt broker](https://www.home-assistant.io/integrations/mqtt/) activated or have an external mqtt server
- Mare sure you have the MQTT Integration activated (same link as above)
- Make sure you have public and private GPG Key + Passphrase in place


If not - here is how to generate it

1. Make sure you have a linux box around or install the [visual studio code addon](https://github.com/hassio-addons/addon-vscode/blob/main/README.md) an start a terminal in there

```
root@605acc316f79:/# gpg --gen-key
gpg (GnuPG) 2.2.40; Copyright (C) 2022 g10 Code GmbH
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

Real name: Your Name
Email address: you@yourdomain.com
You selected this USER-ID:
    "Your Name <you@yourdomain.com>"

Change (N)ame, (E)mail, or (O)kay/(Q)uit? O
We need to generate a lot of random bytes. It is a good idea to perform
some other action (type on the keyboard, move the mouse, utilize the
disks) during the prime generation; this gives the random number
generator a better chance to gain enough entropy.

gpg: /root/.gnupg/trustdb.gpg: trustdb created
gpg: directory '/root/.gnupg/openpgp-revocs.d' created
gpg: revocation certificate stored as '/root/.gnupg/openpgp-revocs.d/5F7C581AE8FF77856C625B855D2BACF61ED7F9DF.rev'
public and secret key created and signed.

pub   rsa3072 2024-10-30 [SC] [expires: 2026-10-30]
      5F7C581AE8FF77856C625B855D2BACF61ED7F9DF
uid                      Your Name <you@yourdomain.com>
sub   rsa3072 2024-10-30 [E] [expires: 2026-10-30]
```

> Notice the above key id 5F7C581AE8FF77856C625B855D2BACF61ED7F9DF in the output.  
> In the below example commands, you need to replace the keyid with YOUR key id!

1.1 If gpg is not installed in your instance do the following and repeat step 1

```
 apt update
 apt install gpg
```

2. After successful generation, export the keys in pem format with _YOUR_ given key-id in the output

```
gpg -a --output public.pem --export 5F7C581AE8FF77856C625B855D2BACF61ED7F9DF
```

3. And same for the private key. GPG will ask you for a passphrase which protects your private key. You need the passphrase later on.

```
gpg -a --output private.pem --export-secret-key 5F7C581AE8FF77856C625B855D2BACF61ED7F9DF
```
4. Last, copy both files on your homeassistant instance into the /ssl folder

That´s all - all prerequisites are set! Congrats!

## Options

### S3 Bucketname (s3settings_bucketName)
Specify the bucket the backup go in.  
You will need to create this bucket upfront within your S3.  
The bucketname is without any slashes

### S3 Endpoint URL (s3settings_endpointUrl)
Specify the full URL of your S3 hoster.  
This is the URL without the bucketname.  
The base URL of your S3 server without trailing slash.

### S3 AccessKey (s3settings_accessKey)
Specify your access_key.  
Make sure you have granted access to the above bucket

### S3 SecretKey (s3settings_secretKey)
Specify your secret_key.  
Make sure you have granted access to the above bucket

### GPG Fingerprint (gpgsettings_fingerprint)
The fingerprint id of your public key.  
This is the above mentions key-id which you can find in the output of the key-gen.  
If you have forgotten it, you can get it again with:   
```
gpg --list-keys
```

### GPG Passphrase (gpgsettings_passphrase)
The passphrase your private key is protected with.  
See above documentation!

### GPG Public Key Filename (gpgsettings_publicKey)
The filename under /ssl of your public key.

### GPG Private Key Filename (gpgsettings_privatKey)
The filename under /ssl of your private key.

### Duplicity SourceDir (duplicitysettings_sourceDir)
This should typically contain /backup.

### Duplicity RestoreDir (duplicitysettings_restoreMediaDir)
I decided that it is better to restore the backup in an own folder.  
The base folder is /media.  
You can here specify the restore folder without any slashes!

### Duplicity IncrementFor (duplicitysettings_incrementalFor)
Time to increment before making a new full backup.
This is transfered to duplicity option: full-if-older-than.   
The allowed time formats are: D, W, M, or Y

### Duplicity Remove Older Than (duplicitysettings_removeOlderThan)
Removing older backups from S3.
This is transfered to duplicity options: remove-older-than 
The allowed time formats are: D, W, M, or Y

### Use internal MQTT (mqttsettings_use_internal)
Switch this on to use the internal mqtt server.
With this option the username and password is auto injected.

### MQTT Server Name/IP (mqttsettings_server)
The hostname or IP of your external (non-homeassistnat) mqtt server. Connection will be done on port 1883.

### MQTT Username (mqttsettings_user)
The mqtt user you created

### MQTT Password (mqttsettings_pass)
The password of the mqtt user.

## License
MIT License

Copyright (c) 2021-2024 Franck Nijhof

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
