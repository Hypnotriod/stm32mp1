# STM32MP157F-DK2 Wi-Fi configuration on Debian with WPA Supplicant
* [Connect to Wi-Fi From Terminal on Debian with WPA Supplicant](https://www.linuxbabe.com/debian/connect-to-wi-fi-from-terminal-on-debian-wpa-supplicant)
* [WPA Supplicant - Common definitions](https://w1.fi/wpa_supplicant/devel/defs_8h.html)

* Update `/etc/network/interfaces` with:
```txt
auto wlan0
allow-hotplug wlan0
iface wlan0 inet manual
wpa-conf /etc/wpa_supplicant.conf
```

* Create `/etc/wpa_supplicant.conf`
```txt
ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev
update_config=1
network={
        ssid="ssid"
        psk="password"
        # Hidden network
        # scan_ssid=1
}
```

* Create `/etc/systemd/system/dhclient.service`
```txt
[Unit]
Description= DHCP Client
Before=network.target
After=wpa_supplicant.service

[Service]
Type=forking
ExecStart=/sbin/dhclient wlan0 -v
ExecStop=/sbin/dhclient wlan0 -r
Restart=always
 
[Install]
WantedBy=multi-user.target
```
* Enable dhcp client
```bash
sudo systemctl enable dhclient.service
```

