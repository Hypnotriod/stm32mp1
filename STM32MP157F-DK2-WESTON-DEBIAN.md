# Weston on Debian
* [Wayland Weston overview](https://wiki.st.com/stm32mpu/wiki/Wayland_Weston_overview)
* [Running Weston](https://wayland.pages.freedesktop.org/weston/toc/running-weston.html)
* [weston.ini](https://manpages.ubuntu.com/manpages/focal/man5/weston.ini.5.html)
* [VivanteGPUIP](https://www.verisilicon.com/en/IPPortfolio/VivanteGPUIP)

# Copy the GPU drivers to the SD card
```bash
export ROOTFS=/media/${USER}/rootfs
sudo cp -r sdk/sysroots/cortexa7t2hf-neon-vfpv4-ostl-linux-gnueabi/vendor/lib/* ${ROOTFS}/lib
```

# Copy additional resources (optional)
```bash
sudo mkdir -p ${ROOTFS}/home/debian/Pictures/
sudo cp resources/logo/debian-logo-480-800.jpg ${ROOTFS}/home/debian/Pictures/
```

# Install Weston on target
```bash
sudo apt update
sudo apt install weston
```

* Create `/etc/xdg/weston/weston.ini`
```
[core]
modules=screen-share.so
shell=desktop-shell.so
backend=drm-backend.so
idle-time=0
repaint-window=100
require-input=false
remoting=remoting-plugin.so

[shell]
#background-image=/home/debian/Pictures/debian-logo-480-800.jpg
background-color=0xff000000
background-type=scale-crop
clock-format=minutes
panel-color=0xffaaaaaa
panel-position=bottom
locking=false
animation=none
startup-animation=none
close-animation=none
focus-animation=none
binding-modifier=ctrl

[autolaunch]
#path=/path/to/the/application

[keyboard]
keymap_layout=us

# HDMI connector
[output]
name=HDMI-A-1
mode=off
#mode=preffered

# DSI connector
[output]
name=DSI-1
#mode=off
mode=preferred
transform=rotate-90
app-ids=1000

# LTDC connector
[output]
name=DPI-1
mode=off
#mode=preffered

# LVDS connector
[output]
name=LVDS-1
mode=off
#mode=preffered

[libinput]
touchscreen_calibrator=true
calibration_helper=/bin/echo

[screen-share]
command=/usr/bin/weston --backend=rdp --shell=fullscreen --no-clients-resize
```

* Create `/usr/lib/systemd/user/weston.service`
```
[Unit]
Description=Weston, a Wayland compositor, as a user service
Documentation=man:weston(1) man:weston.ini(5)
Documentation=https://wayland.freedesktop.org/

# Activate using a systemd socket
Requires=weston.socket
After=weston.socket

# Since we are part of the graphical session, make sure we are started before
Before=graphical-session.target

[Service]
Type=notify
TimeoutStartSec=60
WatchdogSec=20
# Defaults to journal
#StandardOutput=journal
StandardError=journal

EnvironmentFile=-/etc/default/debian
Environment="XDG_RUNTIME_DIR=/home/debian"
Environment="WESTON_USER=debian"
Environment="WL_EGL_GBM_FENCE=0"
# add a ~/.config/weston.ini and weston will pick-it up

ExecStart=/usr/bin/weston --modules=systemd-notify.so
ExecStop=/usr/bin/killall weston

[Install]
WantedBy=graphical-session.target
```

* Create `/usr/lib/systemd/user/weston.socket`
```
[Unit]
Description=Weston, a Wayland compositor
Documentation=man:weston(1) man:weston.ini(5)
Documentation=https://wayland.freedesktop.org/

[Socket]
ListenStream=%t/wayland-0
```

* Create `/usr/lib/systemd/system/weston-graphical-session.service`
```
[Unit]
Description=Weston graphical session

# Make sure we are started after logins are permitted.
Requires=systemd-user-sessions.service
After=systemd-user-sessions.service

# If you want you can make it part of the graphical session
#Before=graphical.target

# Not necessary but just in case
#ConditionPathExists=/dev/tty7

# D-Bus is necessary for contacting logind. Logind is required.
Wants=dbus.socket
After=dbus.socket

[Service]
Type=simple
Environment=XDG_SESSION_TYPE=wayland
ExecStart=/bin/systemctl --wait --user start weston.service
RemainAfterExit=yes

# The user to run the session as. Pick one!
User=debian
Group=debian

# Make sure the working directory is the users home directory
WorkingDirectory=/home/debian

# Set up a full user session for the user, required by Weston.
PAMName=debian-autologin

# A virtual terminal is needed.
TTYPath=/dev/tty7
TTYReset=yes
TTYVHangup=yes
TTYVTDisallocate=yes

# Fail to start if not controlling the tty.
StandardInput=tty-fail

# Defaults to journal, in case it doesn't adjust it accordingly
#StandardOutput=journal
StandardError=journal

# Log this user with utmp, letting it show up with commands 'w' and 'who'.
UtmpIdentifier=tty7
UtmpMode=user

[Install]
Alias=display-manager.service
WantedBy=multi-user.target
```

# Start Weston session process (root user)
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
