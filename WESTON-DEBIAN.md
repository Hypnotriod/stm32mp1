# Weston on Debian
* [Wayland Weston overview](https://wiki.st.com/stm32mpu/wiki/Wayland_Weston_overview)
* [Running Weston](https://wayland.pages.freedesktop.org/weston/toc/running-weston.html)
* [weston.ini](https://manpages.ubuntu.com/manpages/focal/man5/weston.ini.5.html)
* [VivanteGPUIP](https://www.verisilicon.com/en/IPPortfolio/VivanteGPUIP)

# Prepare the SD card
```bash
sudo cp -r sdk/sysroots/cortexa7t2hf-neon-vfpv4-ostl-linux-gnueabi/vendor/lib/* /media/${USER}/rootfs/lib
sudo mkdir -p /media/${USER}/rootfs/usr/lib/systemd/user/
sudo mkdir -p /media/${USER}/rootfs/etc/xdg/weston/
sudo cp weston/weston.service /media/${USER}/rootfs/usr/lib/systemd/user/
sudo cp weston/weston.socket /media/${USER}/rootfs/usr/lib/systemd/user/
sudo cp weston/weston-graphical-session.service /media/${USER}/rootfs/usr/lib/systemd/system/
sudo cp weston/weston.ini /media/${USER}/rootfs/etc/xdg/weston/weston.ini
```

# Install Weston on target
```bash
sudo apt update
sudo apt install weston
```

# Start Weston session (root user)
```bash
sudo XDG_RUNTIME_DIR=/run weston --tty=1
```

# Weston session as a systemd service (debian user)
```bash
# Start/Stop
sudo systemctl start weston-graphical-session.service
sudo systemctl status weston-graphical-session.service
sudo systemctl stop weston-graphical-session.service
# Enable on startup
sudo systemctl enable weston-graphical-session.service
```
