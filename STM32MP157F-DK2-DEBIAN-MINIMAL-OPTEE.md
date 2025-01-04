# STM32MP157F-DK2 Debian build on Ubuntu PC
Based on: [debian-getting-started-with-the-stm32mp157](https://forum.digikey.com/t/debian-getting-started-with-the-stm32mp157/12459)

# Prepare environment inside your working directory
```bash
export WORKSPACE_DIR=`pwd`
export UBOOT_DIR=${WORKSPACE_DIR}/u-boot
export OPTEE_DIR=${WORKSPACE_DIR}/optee_os
export MACHINE=stm32mp157f-dk2
```

# STM32mpu SDK
* [STM32MPU Developer Package](https://wiki.st.com/stm32mpu/wiki/STM32MPU_Developer_Package#Installing_the_SDK)
* [Download the SDK](https://www.st.com/en/embedded-software/stm32mp1dev.html#get-software)
```bash
cd ${WORKSPACE_DIR}
tar xvf en.SDK-x86_64-stm32mp1-openstlinux-6.6-yocto-scarthgap-mpu-v24.11.06.tar.gz
chmod +x stm32mp1-openstlinux-6.6-yocto-scarthgap-mpu-v24.11.06/sdk/st-image-weston-openstlinux-weston-stm32mp1.rootfs-x86_64-toolchain-5.0.3-openstlinux-6.6-yocto-scarthgap-mpu-v24.11.06.sh
./stm32mp1-openstlinux-6.6-yocto-scarthgap-mpu-v24.11.06/sdk/st-image-weston-openstlinux-weston-stm32mp1.rootfs-x86_64-toolchain-5.0.3-openstlinux-6.6-yocto-scarthgap-mpu-v24.11.06.sh -d ${WORKSPACE_DIR}/sdk
sudo apt install build-essential swig libgnutls28-dev
sudo apt install libncurses5-dev libncursesw5-dev libyaml-dev u-boot-tools
```

# U-boot
* [U-Boot for stm32mp1](https://docs.u-boot.org/en/latest/board/st/stm32mp1.html)
* [STM32MP15 U-Boot](https://wiki.stmicroelectronics.cn/stm32mpu/wiki/STM32MP25_U-Boot)
* [STM32MP15 U-Boot build](https://wiki.st.com/stm32mpu/wiki/U-Boot_overview#U-Boot_build)
```bash
cd ${WORKSPACE_DIR}
git clone -b v2023.10-stm32mp https://github.com/STMicroelectronics/u-boot --depth=1
cd u-boot/
source ${WORKSPACE_DIR}/sdk/environment-setup
export KBUILD_OUTPUT=${UBOOT_DIR}/out
make distclean
make stm32mp15_defconfig
make DEVICE_TREE=${MACHINE} all
```
Artifacts:
* u-boot/out/u-boot.bin
* u-boot/out/u-boot-nodtb.bin

# OP-TEE
* [How to build OP-TEE components](https://wiki.st.com/stm32mpu/wiki/How_to_build_OP-TEE_components)
```bash
cd ${WORKSPACE_DIR}
# git clone -b 4.0.0-stm32mp https://github.com/STMicroelectronics/optee_os.git --depth=1
git clone -b 3.19.0-stm32mp-r2.1 https://github.com/STMicroelectronics/optee_os.git --depth=1
cd optee_os/
source ${WORKSPACE_DIR}/sdk/environment-setup
unset LDFLAGS;
unset CFLAGS;
make PLATFORM=stm32mp1 CFG_EMBED_DTB_SOURCE_FILE=${MACHINE}.dts CFG_TEE_CORE_LOG_LEVEL=2 O=build all
```
Artifacts:
* optee_os/build/core/tee-header_v2.bin
* optee_os/build/core/tee-pageable_v2.bin
* optee_os/build/core/tee-pager_v2.bin

# Arm trusted firmware
* [Trusted Firmware build instructions](https://trustedfirmware-a.readthedocs.io/en/lts-v2.10/plat/st/stm32mp1.html#build-instructions)
* [Trusted Firmware doc](https://trustedfirmware-a.readthedocs.io/en/stable/plat/st/stm32mp1.html)
* `-pedantic` flag in `HOSTCCFLAGS` of `arm-trusted-firmware/tools/fiptool/Makefile` may cause the `ISO C does not support the '_FloatXX' type` errors during the build. Can be removed.
```bash
cd ${WORKSPACE_DIR}
git clone -b v2.10-stm32mp https://github.com/STMicroelectronics/arm-trusted-firmware.git --depth=1
cd arm-trusted-firmware/
source ${WORKSPACE_DIR}/sdk/environment-setup
unset LDFLAGS;
unset CFLAGS;
make realclean
make PLAT=stm32mp1 \
    STM32MP13=0 \
    STM32MP15=1 \
    STM32MP_SDMMC=1 \
    STM32MP_EMMC=0 \
    ARCH=aarch32 \
    ARM_ARCH_MAJOR=7 \
    AARCH32_SP=optee \
    DTB_FILE_NAME=${MACHINE}.dtb \
    BL33=${UBOOT_DIR}/out/u-boot-nodtb.bin \
    BL33_CFG=${UBOOT_DIR}/out/u-boot.dtb \
    BL32=${OPTEE_DIR}/build/core/tee-header_v2.bin \
    BL32_EXTRA1=${OPTEE_DIR}/build/core/tee-pager_v2.bin \
    BL32_EXTRA2=${OPTEE_DIR}/build/core/tee-pageable_v2.bin \
    all fip
```
Create fip with fiptool (alternative way). Skip this part.
```bash
fiptool create \
    --tos-fw ${OPTEE_DIR}/build/core/tee-header_v2.bin \
    --tos-fw-extra1 ${OPTEE_DIR}/build/core/tee-pager_v2.bin \
    --tos-fw-extra2 ${OPTEE_DIR}/build/core/tee-pageable_v2.bin \
    --nt-fw ${UBOOT_DIR}/out/u-boot-nodtb.bin \
    --hw-config ${UBOOT_DIR}/out/u-boot.dtb \
    build/stm32mp1/release/fip.bin
```
Artifacts:
* arm-trusted-firmware/build/stm32mp1/release/fip.bin
* arm-trusted-firmware/build/stm32mp1/release/tf-a-stm32mp157f-dk2.stm32

# Linux kernel
* [Modify, rebuild and reload the Linux kernel](https://wiki.st.com/stm32mpu/wiki/Getting_started/STM32MP2_boards/STM32MP257x-DK/Develop_on_Arm_Cortex-A35/Modify,_rebuild_and_reload_the_Linux_kernel)
```bash
cd ${WORKSPACE_DIR}
git clone -b v6.6-stm32mp https://github.com/STMicroelectronics/linux.git --depth=1
cd linux
export OUTPUT_BUILD_DIR=$PWD/build
mkdir -p ${OUTPUT_BUILD_DIR}
source ${WORKSPACE_DIR}/sdk/environment-setup
make O="${OUTPUT_BUILD_DIR}" defconfig fragment*.config
cd ${OUTPUT_BUILD_DIR}
```
Set `CONFIG_LOCALVERSION_AUTO` to `n` in `linux/build/.config` to remove the version suffix.  
Also launching `make` with the `LOCALVERSION=` helps to get rid of the `+` sign.
```bash
export IMAGE_KERNEL="uImage"
make ${IMAGE_KERNEL} vmlinux dtbs LOADADDR=0xC2000040 O="${OUTPUT_BUILD_DIR}" LOCALVERSION=
make modules O="${OUTPUT_BUILD_DIR}" LOCALVERSION=
make INSTALL_MOD_PATH="${OUTPUT_BUILD_DIR}/install_artifact" modules_install O="${OUTPUT_BUILD_DIR}" LOCALVERSION=
mkdir -p ${OUTPUT_BUILD_DIR}/install_artifact/boot/dtbs/
rm ${OUTPUT_BUILD_DIR}/install_artifact/lib/modules/6.6.48/build
cp ${OUTPUT_BUILD_DIR}/arch/${ARCH}/boot/${IMAGE_KERNEL} ${OUTPUT_BUILD_DIR}/install_artifact/boot/
find ${OUTPUT_BUILD_DIR}/arch/${ARCH}/boot/dts/ -name 'st*.dtb' -exec cp '{}' ${OUTPUT_BUILD_DIR}/install_artifact/boot/dtbs/ \;
```
Artifacts:
* linux/build/install_artifact/boot/uImage
* linux/build/install_artifact/boot/dtbs/stm32mp157f-dk2.dtb
* linux/build/install_artifact/lib/modules/6.6.48/

# Debian rootfs
```bash
cd ${WORKSPACE_DIR}
wget -c https://rcn-ee.com/rootfs/eewiki/minfs/debian-12.1-minimal-armhf-2023-08-22.tar.xz
tar xf debian-12.1-minimal-armhf-2023-08-22.tar.xz 
```

# Populate the SD card
* [Flash Layout SD card](https://wiki.st.com/stm32mpu/wiki/STM32CubeProgrammer_flashlayout#SD_card)  
* [STM32 MPU Flash mapping](https://wiki.st.com/stm32mpu/wiki/STM32_MPU_Flash_mapping)  

Call `lsblk` to determine the device entry for the SD card.  
In case of `/dev/sdX` do:
```bash
export DISK=/dev/sdX
export DISK_P=${DISK}
```
In case of `/dev/mmcblkX` do:
```bash
export DISK=/dev/mmcblkX
export DISK_P=${DISK}p
```
```bash
cd ${WORKSPACE_DIR}
umount ${DISK_P}X
sudo dd if=/dev/zero of=${DISK} bs=1M count=10
sudo sgdisk -o ${DISK}
sudo sgdisk --resize-table=128 -a 1 \
    -n 1:34:545    -c 1:fsbl1   \
    -n 2:546:1057  -c 2:fsbl2   \
    -n 3:1058:5153 -c 3:fip     \
    -n 4:5154:6177 -c 4:u-boot-env \
    -n 5:6178:     -c 5:rootfs  \
    -A 5:set:2                  \
    -p ${DISK}
sudo dd if=./arm-trusted-firmware/build/stm32mp1/release/tf-a-${MACHINE}.stm32 of=${DISK_P}1
sudo dd if=./arm-trusted-firmware/build/stm32mp1/release/tf-a-${MACHINE}.stm32 of=${DISK_P}2
sudo dd if=./arm-trusted-firmware/build/stm32mp1/release/fip.bin of=${DISK_P}3
sudo dd if=/dev/zero of=${DISK_P}4 bs=512K count=1
sudo mkfs.ext4 -L rootfs ${DISK_P}5
sudo tar xfvp ./debian-*-*-armhf-*/armhf-rootfs-*.tar -C /media/${USER}/rootfs/
sync
sudo mkdir -p /media/${USER}/rootfs/boot/extlinux/
sudo sh -c "echo 'LABEL Linux' > /media/${USER}/rootfs/boot/extlinux/extlinux.conf"
sudo sh -c "echo '    KERNEL /boot/uImage' >> /media/${USER}/rootfs/boot/extlinux/extlinux.conf"
sudo sh -c "echo '    APPEND console=ttySTM0,115200 root=/dev/mmcblk0p5 rw rootfstype=ext4 rootwait' >> /media/${USER}/rootfs/boot/extlinux/extlinux.conf"
sudo sh -c "echo '    FDTDIR /boot/dtbs' >> /media/${USER}/rootfs/boot/extlinux/extlinux.conf"
sudo cp -rv ./linux/build/install_artifact/boot/* /media/${USER}/rootfs/boot/
sudo cp -rv ./linux/build/install_artifact/lib/* /media/${USER}/rootfs/lib/
sudo sh -c "echo '/dev/mmcblk0p5  /  auto  errors=remount-rw  0  1' >> /media/${USER}/rootfs/etc/fstab"
sync
```

