# gotk3-example
* [Create a simple hello-world application](https://wiki.st.com/stm32mpu/wiki/Getting_started/STM32MP1_boards/STM32MP135x-DK/Develop_on_Arm%C2%AE_Cortex%C2%AE-A7/Create_a_simple_hello-world_application)
* [gotk3](https://github.com/gotk3/gotk3)
* [GTK CSS properties](https://docs.gtk.org/gtk3/css-properties.html)

Go `GTK3` build for stm32mp1 `armv7hf` in `Docker` container example.

## Build locally (Ubuntu/Debian)
```bash
apt install -y build-essential libgtk-3-dev libgtk-4-dev libgirepository1.0-dev
make build-gotk3-example
```

## Build for the target platform inside the Docker container
* Initial build (will take about an hour):
```bash
make build-armv7hf-docker-container
# Windows:
# make build-armv7hf-gotk3-example-initial CURDIR=%cd%
make build-armv7hf-gotk3-example-initial
```

* Consecutive builds:
```bash
make build-armv7hf-gotk3-example
```
