# STM32MP157F-DK2 Debian build on Ubuntu PC
Based on: [debian-getting-started-with-the-stm32mp157](https://forum.digikey.com/t/debian-getting-started-with-the-stm32mp157/12459)

## Prepare the environment inside your working directory
```bash
export WORKSPACE_DIR=$PWD
export SDK_DIR=${WORKSPACE_DIR}/sdk
export UBOOT_DIR=${WORKSPACE_DIR}/u-boot
export OPTEE_DIR=${WORKSPACE_DIR}/optee_os
```
Choose the machine name:
```bash
export MACHINE=stm32mp157f-dk2
```
```bash
export MACHINE=stm32mp157c-dk2
```
```bash
export MACHINE=stm32mp157c-odyssey
```

## STM32mpu SDK
* [STM32MPU Developer Package](https://wiki.st.com/stm32mpu/wiki/STM32MPU_Developer_Package#Installing_the_SDK)  

Download the latest [STM32MP1 OpenSTLinux Developer Package](https://www.st.com/en/embedded-software/stm32mp1dev.html#get-software) .tar archive.  
Extract and install the SDK:   
```bash
tar xvf en.SDK-x86_64-stm32mp1-openstlinux-6.6-yocto-scarthgap-mpu-v24.11.06.tar.gz
chmod +x stm32mp1-openstlinux-*/sdk/st-image-weston-*.sh
./stm32mp1-openstlinux-*/sdk/st-image-weston-*.sh -d ${SDK_DIR}
sudo apt install build-essential swig libgnutls28-dev
sudo apt install libncurses5-dev libncursesw5-dev libyaml-dev u-boot-tools
```

## U-boot (Universal Bootloader)
* [U-Boot for stm32mp1](https://docs.u-boot.org/en/latest/board/st/stm32mp1.html)
* [STM32MP15 U-Boot](https://wiki.stmicroelectronics.cn/stm32mpu/wiki/STM32MP25_U-Boot)
* [STM32MP15 U-Boot build](https://wiki.st.com/stm32mpu/wiki/U-Boot_overview#U-Boot_build)

<details>
<summary>stm32mp157c-odyssey patch</summary>

Based on the [meta-st-odyssey](https://github.com/Seeed-Studio/meta-st-odyssey/tree/openstlinux-6.1-yocto-mickledore). Execute after the `git clone` command:
```bash
cp resources/stm32mp157c-odyssey/u-boot-patch/stm32mp1.c u-boot/board/st/stm32mp1/
cp resources/stm32mp157c-odyssey/u-boot-patch/stm32mp15_defconfig u-boot/configs/
cp resources/stm32mp157c-odyssey/u-boot-patch/stm32mp157c-odyssey-som-u-boot.dtsi u-boot/arch/arm/dts/
cp resources/stm32mp157c-odyssey/u-boot-patch/stm32mp157c-odyssey-som.dtsi u-boot/arch/arm/dts/
cp resources/stm32mp157c-odyssey/u-boot-patch/stm32mp157c-odyssey.dts u-boot/arch/arm/dts/
```
</details>

```bash
cd ${WORKSPACE_DIR}
git clone -b v2023.10-stm32mp-r1.2 https://github.com/STMicroelectronics/u-boot --depth=1
cd u-boot/
source ${SDK_DIR}/environment-setup
export KBUILD_OUTPUT=${UBOOT_DIR}/out
make distclean
make stm32mp15_defconfig
make DEVICE_TREE=${MACHINE} all
```
Artifacts:
* u-boot/out/u-boot.dtb
* u-boot/out/u-boot-nodtb.bin

## OP-TEE (Open Portable Trusted Execution Environment)
* [How to build OP-TEE components](https://wiki.st.com/stm32mpu/wiki/How_to_build_OP-TEE_components)

<details>
<summary>stm32mp157c-odyssey patch</summary>

Based on the [meta-st-odyssey](https://github.com/Seeed-Studio/meta-st-odyssey/tree/openstlinux-6.1-yocto-mickledore). Execute after the `git clone` command:
```bash
cp resources/stm32mp157c-odyssey/optee_os-patch/clk-stm32mp15.c optee_os/core/drivers/clk/
cp resources/stm32mp157c-odyssey/optee_os-patch/conf.mk optee_os/core/arch/arm/plat-stm32mp1/
cp resources/stm32mp157c-odyssey/optee_os-patch/platform_config.h optee_os/core/arch/arm/plat-stm32mp1/
cp resources/stm32mp157c-odyssey/optee_os-patch/shared_resources.c optee_os/core/arch/arm/plat-stm32mp1/
cp resources/stm32mp157c-odyssey/optee_os-patch/stm32mp157c-odyssey.dts optee_os/core/arch/arm/dts/
cp resources/stm32mp157c-odyssey/optee_os-patch/stm32mp157c-odyssey.dtsi optee_os/core/arch/arm/dts/
```
</details>

```bash
cd ${WORKSPACE_DIR}
git clone -b 4.0.0-stm32mp-r1.2 https://github.com/STMicroelectronics/optee_os.git --depth=1
cd optee_os/
source ${SDK_DIR}/environment-setup
unset LDFLAGS;
unset CFLAGS;
make PLATFORM=stm32mp1 CFG_EMBED_DTB_SOURCE_FILE=${MACHINE}.dts \
    CFLAGS32=--sysroot=${SDKTARGETSYSROOT} \
    CFG_TEE_CORE_LOG_LEVEL=2 O=build all
```
Artifacts:
* optee_os/build/core/tee-header_v2.bin
* optee_os/build/core/tee-pageable_v2.bin
* optee_os/build/core/tee-pager_v2.bin

## TF-A (Arm Trusted Firmware-A) / FIP (Firmware Image Package)
* [Trusted Firmware build instructions](https://trustedfirmware-a.readthedocs.io/en/lts-v2.10/plat/st/stm32mp1.html#build-instructions)
* [Trusted Firmware doc](https://trustedfirmware-a.readthedocs.io/en/stable/plat/st/stm32mp1.html)
* The `-pedantic` flag in `HOSTCCFLAGS` from `arm-trusted-firmware/tools/fiptool/Makefile` may cause the `ISO C does not support the '_FloatXX' type` errors during the fiptool compilation. Can be removed.

<details>
<summary>stm32mp157c-odyssey patch</summary>

Based on the [meta-st-odyssey](https://github.com/Seeed-Studio/meta-st-odyssey/tree/openstlinux-6.1-yocto-mickledore). Execute after the `git clone` command:
```bash
cp resources/stm32mp157c-odyssey/arm-trusted-firmware-patch/stm32mp1_def.h arm-trusted-firmware/plat/st/stm32mp1/
cp resources/stm32mp157c-odyssey/arm-trusted-firmware-patch/stm32mp1_shared_resources.c arm-trusted-firmware/plat/st/stm32mp1/
cp resources/stm32mp157c-odyssey/arm-trusted-firmware-patch/stm32mp157c-odyssey-som.dtsi arm-trusted-firmware/fdts/
cp resources/stm32mp157c-odyssey/arm-trusted-firmware-patch/stm32mp157c-odyssey.dts arm-trusted-firmware/fdts/
```
</details>

```bash
cd ${WORKSPACE_DIR}
git clone -b v2.10-stm32mp-r1.2 https://github.com/STMicroelectronics/arm-trusted-firmware.git --depth=1
cd arm-trusted-firmware/
# Remove the -pedantic flag to be able to compile the fiptool
sed -i 's/ -pedantic//g' tools/fiptool/Makefile
source ${SDK_DIR}/environment-setup
unset LDFLAGS;
unset CFLAGS;
make realclean
```
In case of the `SD card` do:
```bash
make PLAT=stm32mp1 \
    STM32MP13=0 \
    STM32MP15=1 \
    STM32MP_SDMMC=1 \
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
In case of the `eMMC` do:
```bash
make PLAT=stm32mp1 \
    STM32MP13=0 \
    STM32MP15=1 \
    STM32MP_EMMC=1 \
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

<details>
<summary>Create fip with the fiptool (alternative way).</summary>

```bash
fiptool create \
    --tos-fw ${OPTEE_DIR}/build/core/tee-header_v2.bin \
    --tos-fw-extra1 ${OPTEE_DIR}/build/core/tee-pager_v2.bin \
    --tos-fw-extra2 ${OPTEE_DIR}/build/core/tee-pageable_v2.bin \
    --nt-fw ${UBOOT_DIR}/out/u-boot-nodtb.bin \
    --hw-config ${UBOOT_DIR}/out/u-boot.dtb \
    build/stm32mp1/release/fip.bin
```
</details>

Artifacts:
* arm-trusted-firmware/build/stm32mp1/release/fip.bin
* arm-trusted-firmware/build/stm32mp1/release/tf-a-stm32mp157f-dk2.stm32

## Linux kernel
* [Modify, rebuild and reload the Linux kernel](https://wiki.st.com/stm32mpu/wiki/Getting_started/STM32MP2_boards/STM32MP257x-DK/Develop_on_Arm_Cortex-A35/Modify,_rebuild_and_reload_the_Linux_kernel)

<details>
<summary>stm32mp157c-odyssey patch</summary>

Based on the [meta-st-odyssey](https://github.com/Seeed-Studio/meta-st-odyssey/tree/openstlinux-6.1-yocto-mickledore). Execute after the `git clone` command:
```bash
cp resources/stm32mp157c-odyssey/linux-patch/stm32mp157c-odyssey-scmi.dtsi linux/arch/arm/boot/dts/st/
cp resources/stm32mp157c-odyssey/linux-patch/stm32mp157c-odyssey-som.dtsi linux/arch/arm/boot/dts/st/
cp resources/stm32mp157c-odyssey/linux-patch/stm32mp157c-odyssey.dts linux/arch/arm/boot/dts/st/
```
</details>

```bash
cd ${WORKSPACE_DIR}
git clone -b v6.6-stm32mp-r1.2 https://github.com/STMicroelectronics/linux.git --depth=1
cd linux
export OUTPUT_BUILD_DIR=$PWD/build
mkdir -p ${OUTPUT_BUILD_DIR}
source ${SDK_DIR}/environment-setup
make O="${OUTPUT_BUILD_DIR}" defconfig fragment*.config
cd ${OUTPUT_BUILD_DIR}
sed -i 's/CONFIG_LOCALVERSION_AUTO=y/CONFIG_LOCALVERSION_AUTO=n/g' .config
export IMAGE_KERNEL="uImage"
make ${IMAGE_KERNEL} vmlinux dtbs LOADADDR=0xC2000040 O="${OUTPUT_BUILD_DIR}" LOCALVERSION=
make modules O="${OUTPUT_BUILD_DIR}" LOCALVERSION=
make INSTALL_MOD_PATH="${OUTPUT_BUILD_DIR}/install_artifact" modules_install O="${OUTPUT_BUILD_DIR}" LOCALVERSION=
mkdir -p ${OUTPUT_BUILD_DIR}/install_artifact/boot/dtbs/
rm ${OUTPUT_BUILD_DIR}/install_artifact/lib/modules/6.6.48/build
cp ${OUTPUT_BUILD_DIR}/arch/${ARCH}/boot/${IMAGE_KERNEL} ${OUTPUT_BUILD_DIR}/install_artifact/boot/
cp ${OUTPUT_BUILD_DIR}/arch/${ARCH}/boot/dts/st/${MACHINE}.dtb ${OUTPUT_BUILD_DIR}/install_artifact/boot/dtbs/
# To copy all the available ST *.dtb files do:
# find ${OUTPUT_BUILD_DIR}/arch/${ARCH}/boot/dts/ -name 'st*.dtb' -exec cp '{}' ${OUTPUT_BUILD_DIR}/install_artifact/boot/dtbs/ \;
```
Artifacts:
* linux/build/install_artifact/boot/uImage
* linux/build/install_artifact/boot/dtbs/stm32mp157f-dk2.dtb
* linux/build/install_artifact/lib/modules/6.6.48/

## Debian rootfs
```bash
cd ${WORKSPACE_DIR}
wget -c https://rcn-ee.com/rootfs/eewiki/minfs/debian-12.9-minimal-armhf-2025-02-05.tar.xz
tar xf debian-12.9-minimal-armhf-2025-02-05.tar.xz
export ROOTFS_TAR=${WORKSPACE_DIR}/debian-12.9-minimal-armhf-2025-02-05/armhf-rootfs-debian-bookworm.tar
```

## Populate the SD card or create the eMMC / SD card image file
* [Flash Layout](https://wiki.st.com/stm32mpu/wiki/STM32CubeProgrammer_flashlayout)
* [STM32 MPU Flash mapping](https://wiki.st.com/stm32mpu/wiki/STM32_MPU_Flash_mapping)
* [extlinux.conf Menu Customization](https://www.willhaley.com/blog/extlinux-menu/)
* [Debian logo wallpaper](https://github.com/shriramters/wallpapers/blob/main/bin/debian-swirl-4k-dark.png)

Call `lsblk` to determine the device entry for the SD card.  
In case of `sdX` (replace the `sdX` with the correct one) do:
```bash
export DISK=/dev/sdX
export DISK_P=${DISK}
```
In case of `mmcblkX` (replace the `mmcblkX` with the correct one) do:
```bash
export DISK=/dev/mmcblkX
export DISK_P=${DISK}p
```
Also unmount all the previous SD card partitions and erase the partition table:
```bash
umount ${DISK_P}X
sudo dd if=/dev/zero of=${DISK} bs=1M count=10
```
Alternatively to make an `SD card` / `eMMC` image file (2GB example) using the `loop device` do:
```bash
cd ${WORKSPACE_DIR}
export IMAGE_FILE=${MACHINE}.img
dd if=/dev/zero of=${IMAGE_FILE} bs=1G count=2
export DISK=$(sudo losetup --partscan --show --find ${IMAGE_FILE})
export DISK_P=${DISK}p
```
Format the disk and populate with the artifacts  
In case of the `SD card` or `SD card image` do:
```bash
cd ${WORKSPACE_DIR}
export ROOTFS_PARTUUID=e91c4e10-16e6-4c0e-bd0e-77becf4a3582
sudo sgdisk -o ${DISK}
sudo sgdisk --resize-table=128 -a 1 \
    -n 1:34:545    -c 1:fsbl1   \
    -n 2:546:1057  -c 2:fsbl2   \
    -n 3:1058:5153 -c 3:fip     \
    -n 4:5154:6177 -c 4:u-boot-env \
    -n 5:6178:     -c 5:rootfs  \
    -A 5:set:2                  \
    -u 5:${ROOTFS_PARTUUID}     \
    -p ${DISK}
sudo dd if=./arm-trusted-firmware/build/stm32mp1/release/tf-a-${MACHINE}.stm32 of=${DISK_P}1
sudo dd if=./arm-trusted-firmware/build/stm32mp1/release/tf-a-${MACHINE}.stm32 of=${DISK_P}2
sudo dd if=./arm-trusted-firmware/build/stm32mp1/release/fip.bin of=${DISK_P}3
sudo dd if=/dev/zero of=${DISK_P}4 bs=512K count=1
sudo mkfs.ext4 -L rootfs ${DISK_P}5
```

In case of the `eMMC image` do:
```bash
cd ${WORKSPACE_DIR}
export ROOTFS_PARTUUID=491f6117-415d-4f53-88c9-6e0de54deac6
sudo sgdisk -o ${DISK}
sudo sgdisk --resize-table=128 -a 1 \
    -n 1:1024:5119 -c 1:fip     \
    -n 2:5120:6143 -c 2:u-boot-env \
    -n 3:6144:     -c 3:rootfs  \
    -A 3:set:2                  \
    -u 3:${ROOTFS_PARTUUID}     \
    -p ${DISK}
sudo dd if=./arm-trusted-firmware/build/stm32mp1/release/fip.bin of=${DISK_P}1
sudo dd if=/dev/zero of=${DISK_P}2 bs=512K count=1
sudo mkfs.ext4 -L rootfs ${DISK_P}3
```

Mount the `rootfs` partition. The default `/media/${USER}/rootfs` path is used.
```bash
export ROOTFS=/media/${USER}/rootfs
sudo tar xpvf ${ROOTFS_TAR} -C ${ROOTFS}/
sync
sudo mkdir -p ${ROOTFS}/boot/extlinux/
# Skip the next 2 lines if you do not need the U-Boot splash screen:
sudo cp resources/logo/debian-logo-480-800-16bit.bmp ${ROOTFS}/boot/
sudo sh -c "echo 'MENU BACKGROUND /boot/debian-logo-480-800-16bit.bmp' >> ${ROOTFS}/boot/extlinux/extlinux.conf"
sudo sh -c "echo 'TIMEOUT 10' >> ${ROOTFS}/boot/extlinux/extlinux.conf"
sudo sh -c "echo 'DEFAULT Linux' >> ${ROOTFS}/boot/extlinux/extlinux.conf"
sudo sh -c "echo 'LABEL Linux' >> ${ROOTFS}/boot/extlinux/extlinux.conf"
sudo sh -c "echo '    KERNEL /boot/uImage' >> ${ROOTFS}/boot/extlinux/extlinux.conf"
sudo sh -c "echo '    APPEND console=ttySTM0,115200 root=PARTUUID=${ROOTFS_PARTUUID} rw rootfstype=ext4 rootwait' >> ${ROOTFS}/boot/extlinux/extlinux.conf"
sudo sh -c "echo '    FDTDIR /boot/dtbs' >> ${ROOTFS}/boot/extlinux/extlinux.conf"
sudo cp -rv ./linux/build/install_artifact/boot/* ${ROOTFS}/boot/
sudo cp -rv ./linux/build/install_artifact/lib/* ${ROOTFS}/lib/
sudo sh -c "echo 'PARTUUID=${ROOTFS_PARTUUID}  /  auto  errors=remount-ro  0  1' > ${ROOTFS}/etc/fstab"
sudo mkdir -p ${ROOTFS}/boot/firmware/
sudo cp resources/sysconf.txt ${ROOTFS}/boot/firmware/
sudo sh -c "echo 'Debian GNU/Linux 12 \134n \l \n' > ${ROOTFS}/etc/issue"
# Copy/replace the Broadcom/Cypress 802.11 wireless card firmware
sudo cp -rv ${SDK_DIR}/sysroots/cortexa7t2hf-neon-vfpv4-ostl-linux-gnueabi/usr/lib/firmware/* ${ROOTFS}/usr/lib/firmware/
# Copy the GPU drivers
sudo cp -rv ${SDK_DIR}/sysroots/cortexa7t2hf-neon-vfpv4-ostl-linux-gnueabi/vendor/lib/* ${ROOTFS}/lib/
sync
```

<details>
<summary>stm32mp157c-odyssey patch</summary>

Configure the Wi-Fi module firmware
```bash
sudo cp resources/stm32mp157c-odyssey/brcmfmac43430-sdio.seeed,stm32mp157c-odyssey.txt ${ROOTFS}/usr/lib/firmware/brcm/
sudo ln -sr ${ROOTFS}/usr/lib/firmware/brcm/brcmfmac43430-sdio.bin ${ROOTFS}/usr/lib/firmware/brcm/brcmfmac43430-sdio.seeed,stm32mp157c-odyssey.bin
```
</details>

Edit the `${ROOTFS}/boot/firmware/sysconf.txt` file to setup essential system configurations, such as `user name`, `user password`, `hostname`, etc, at the system boot using the `bbbio-set-sysconf.service`. You can disable it with `sudo systemctl disable bbbio-set-sysconf`  
Unmount the rootfs partition and detach a loop device (if used)  
```bash
sudo umount ${ROOTFS}
# In case of the loop device
sudo losetup -d ${DISK}
```

## Flash the eMMC image on the target device from Linux
* [Bootable eMMC](https://linux-sunxi.org/Bootable_eMMC)
* [mmc-utils / mmc](https://manpages.debian.org/testing/mmc-utils/mmc.1.en.html)  

Transfer the `arm-trusted-firmware/build/stm32mp1/release/tf-a-*.stm32` TF-A file (built with the `STM32MP_EMMC=1` option) to the target device.  
Transfer the eMMC image `*.img` file to the target device.  
On target, call `lsblk` to determine the `eMMC` block device name.  
Set the correct values for the `BLOCK_DEVICE` and `MACHINE` environment variables:  
```bash
export BLOCK_DEVICE=mmcblkX
export MACHINE=stm32mp157c-odyssey
sudo apt install -y mmc-utils
sudo sh -c "echo 0 > /sys/block/${BLOCK_DEVICE}boot0/force_ro"
sudo sh -c "echo 0 > /sys/block/${BLOCK_DEVICE}boot1/force_ro"
sudo dd if=tf-a-${MACHINE}.stm32 of=/dev/${BLOCK_DEVICE}boot0 conv=notrunc
sudo dd if=tf-a-${MACHINE}.stm32 of=/dev/${BLOCK_DEVICE}boot1 conv=notrunc
sudo mmc bootbus set single_backward x1 x1 /dev/${BLOCK_DEVICE}
sudo mmc bootpart enable 1 1 /dev/${BLOCK_DEVICE}
sudo dd if=${MACHINE}.img of=/dev/${BLOCK_DEVICE}
# Grow the rootfs partition 3 to the full eMMC capacity
sudo parted /dev/${BLOCK_DEVICE} resizepart 3 -- -1
sudo e2fsck -f /dev/${BLOCK_DEVICE}p3
sudo resize2fs /dev/${BLOCK_DEVICE}p3
```

## In case of an image file, after it was flashed to a microSD card
Use the following code snippet to grow the `rootfs` partition to the full microSD card capacity:  
Call `lsblk` to determine the device entry for the SD card.  
In case of `sdX` (replace the `sdX` with the correct one) do:
```bash
export DISK=/dev/sdX
export DISK_P=${DISK}
```
In case of `mmcblkX` (replace the `mmcblkX` with the correct one) do:
```bash
export DISK=/dev/mmcblkX
export DISK_P=${DISK}p
```
```bash
# Unmount the rootfs partition 5 of the SD card
sudo umount ${DISK_P}5
# Resize the rootfs partition 5 to the full capacity
sudo parted ${DISK} resizepart 5 -- -1
# Check the rootfs ext4 partition for bad blocks
sudo e2fsck -f ${DISK_P}5
# Resize the rootfs partition 5 file system
sudo resize2fs ${DISK_P}5
```

## Issues
* [systemd-binfmt fails on boot if binfmt is missing from fstab](https://github.com/systemd/systemd/issues/28501)  
Replace `ConditionPathIsReadWrite=/proc/sys/` with `ConditionPathIsMountPoint=/proc/sys/fs/binfmt_misc` in `/usr/lib/systemd/system/systemd-binfmt.service`
```txt
#ConditionPathIsReadWrite=/proc/sys/
ConditionPathIsMountPoint=/proc/sys/fs/binfmt_misc
```
* `STM32MP157C ODYSSEY` Bluetooth is not working.  
  No solution yet.
