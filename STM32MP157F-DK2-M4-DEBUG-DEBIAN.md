# STM32MP157F-DK2 M4 core debug in Production Mode with STM32CubeIDE on Debian 
* [Linux remoteproc framework overview](https://wiki.st.com/stm32mpu/wiki/Linux_remoteproc_framework_overview)
* [Remote debugging using gdbserver](https://wiki.st.com/stm32mpu/wiki/GDB#Remote_debugging_using_gdbserver)
* [Modify, rebuild and reload a firmware](https://wiki.st.com/stm32mpu/wiki/Getting_started/STM32MP1_boards/STM32MP157x-DK2/Develop_on_Arm%C2%AE_Cortex%C2%AE-M4/Modify,_rebuild_and_reload_a_firmware)
* [How to setup target password in STM32CubeIDE](https://wiki.stmicroelectronics.cn/stm32mpu/wiki/How_to_setup_target_password_in_STM32CubeIDE)
* [Error executing event gdb-attach on target](https://forum.digikey.com/t/debian-on-stm32mp157-debug-cm4-core-in-stm32cubeide/15533)

## Configure the target board:
* Install the `gdbserver`
```bash
sudo apt install gdbserver
```
* Set a password for the `root` user and update the `/etc/ssh/sshd_config` with:
```txt
PermitRootLogin yes
```

## Generate a "blink" project in the `STM32CubeMX`:
* Choose the `STM32MP157F-DK2` board at the `Board Selector`
* `Initialize all the peripherals with default mode` -> `Yes`
* Search for the `PH7` or `LED_Y` in `Pinout View`
* Left click on the highlighted pin -> set `GPIO_Output`
* Right click on the highlighted pin -> `Pin Reservation` -> `Cortex M4`
* Fill the project settings in the `Project Manager` tab. Select the `STM32CubeIDE` option in the `Toolchain / IDE` and generate the project.
* Update the `main.c` with the "blinky" code:
```c
  /* USER CODE BEGIN WHILE */
  while (1)
  {
    HAL_GPIO_WritePin(LED_Y_GPIO_Port, LED_Y_Pin, GPIO_PIN_SET);
    HAL_Delay(500);
    HAL_GPIO_WritePin(LED_Y_GPIO_Port, LED_Y_Pin, GPIO_PIN_RESET);
    HAL_Delay(500);
    /* USER CODE END WHILE */

    /* USER CODE BEGIN 3 */
  }
```
* In `STM32CubeIDE` specify the `Inet address` of the target board at `Run -> Debug Configurations -> Debugger`
* In `STM32CubeIDE` update the password for the `root` user at `Window -> Preferences -> Remote Development -> Remote Connections -> Password based authentication`
* To fix the `Error executing event gdb-attach on target` issue do:  
In autogenerated `./CM4/RemoteProc/fw_cortex_m4.sh` file replace/fix next two lines:
```bash
#if [ $1 == "start" ]
if [ "$1" = "start" ]
...
#if [ $1 == "stop" ]
if [ "$1" = "stop" ]
```
