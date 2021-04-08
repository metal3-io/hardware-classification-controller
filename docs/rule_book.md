# Rule book to identify disk type for different vendors

For `disk` selection user can use combination of `hctl` and
`rotational` parameters.

For `NIC` user have to provide vendor ID.

## `HCTL`

   Here `hctl` represents:

   1. SCSI adapter number [host]

   1. channel number [bus]

   1. id number [target]

   1. number of logical units [lun]

## `Rotational`

   Rotational value will be true if disk is HDD and false
   represents individual SSD, also software RAID can be SSD RAID or
   HDD RAID with rotational value true.

## Vendor Dell

   We are taking example of Dell hardware here.

### Disk

   1. Individual HDD/SSD : with `hctl` as 0:0:N:0 and rotational flag as
   True/False.
   1. PERC RAID of HDDs/SSDs : with `hctl` as 0:N:0:0/0:N:N:0 and rotational
   flag as True.
   1. Dell BOSS Controller Individual SSDs : with `hctl` as N:0:0:0 and
   rotational flag as False.
   1. Dell BOSS Controller Virtual Disk (RAID) : with `hctl` as N:0:0:0 and
   rotational flag as True.
   1. NVMe : No `hctl` Pattern, rotational flag as False and model name
   contains NVMe keyword.

### NIC

   1. Intel NIC Vendor ID is 0x8086.
   1. Mellanox NIC Vendor ID is 0x15b3.
   1. Broadcom NIC Vendor ID is 0x14e4.
