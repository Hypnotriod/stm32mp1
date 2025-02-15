# STM32MP157F-DK2 Wi-Fi configuration on Debian
* [How to setup a WLAN connection](https://wiki.st.com/stm32mpu/wiki/How_to_setup_a_WLAN_connection)
* [Connect to Wi-Fi From Terminal on Debian with WPA Supplicant](https://www.linuxbabe.com/debian/connect-to-wi-fi-from-terminal-on-debian-wpa-supplicant)
* [WPA Supplicant - Common definitions](https://w1.fi/wpa_supplicant/devel/defs_8h.html)

## WPA Supplicant
* Update `/etc/network/interfaces` with:
```txt
auto wlan0
iface wlan0 inet dhcp
wpa-conf /etc/wpa_supplicant/wpa_supplicant.conf
```

* Create `/etc/wpa_supplicant/wpa_supplicant.conf`  
Use `wpa_passphrase your_ssid your_password | grep -vE "#psk"` to generate WPA PSK from an ASCII passphrase
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

* Restart network interface
```bash
sudo systemctl restart systemd-networkd.service
# or with networking services
sudo systemctl restart networking.service
```
