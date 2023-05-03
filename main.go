package main

import (
    "fmt"
    "hash/crc32"
    "net"
    "os/exec"
    "runtime"
    "strconv"
    "strings"
)

func main() {
    // 获取 CPU ID、硬盘序列号和网卡 MAC 地址
    cpuID, err := getCPUID()
    if err != nil {
        fmt.Println(err)
        return
    }
    diskID, err := getDiskID()
    if err != nil {
        fmt.Println(err)
        return
    }
    macAddr, err := getMACAddr()
    if err != nil {
        fmt.Println(err)
        return
    }

    // 将 CPU ID、硬盘序列号和网卡 MAC 地址合并，并计算其 CRC32 校验值
    machineCode := strconv.FormatUint(uint64(crc32.ChecksumIEEE([]byte(cpuID+diskID+macAddr))), 10)

    // 将机器码按每四个字符一组进行分组，并添加分隔符 "-"
    machineCodeWithDash := strings.Join(splitByLength(machineCode, 4), "-")

    // 输出机器码
    fmt.Println("Machine code: " + machineCodeWithDash)
}

// 获取 CPU ID
func getCPUID() (string, error) {
    var cmd string
    if runtime.GOOS == "windows" {
        // Windows 上使用 WMIC 命令获取 CPU ID
        cmd = "wmic cpu get ProcessorId"
    } else if runtime.GOOS == "darwin" {
        // macOS 上使用 system_profiler 命令获取 CPU ID
        cmd = "system_profiler SPHardwareDataType | grep 'Processor Name' | awk '{print $4}'"
    } else {
        // Linux 上使用 cat /proc/cpuinfo 命令获取 CPU ID
        cmd = "cat /proc/cpuinfo"
    }
    output, err := exec.Command("sh", "-c", cmd).Output()
    if err != nil {
        return "", err
    }
    cpuID := ""
    if runtime.GOOS == "windows" {
        // Windows 上 CPU ID 在第二行，包含在输出的第二个单词中
        lines := strings.Split(string(output), "\n")
        if len(lines) > 1 {
            words := strings.Fields(lines[1])
            if len(words) > 1 {
                cpuID = words[1]
            }
        }
    } else if runtime.GOOS == "darwin" {
        // macOS 上 CPU ID 包含在 Processor Name 的值中
        lines := strings.Split(string(output), "\n")
        for _, line := range lines {
            if strings.HasPrefix(line, "Processor Name") {
                cpuID = strings.TrimSpace(strings.TrimPrefix(line, "Processor Name:"))
                break
            }
        }
    } else {
        // Linux 上 CPU ID 包含在 CPU 特性的值中
        lines := strings.Split(string(output), "\n")
        for _, line := range lines {
            if strings.HasPrefix(line, "cpu family") {
                words := strings.Fields(line)
                if len(words) > 2 {
                    cpuID = words[2]
                }
                break
            }
        }
    }
    return cpuID, nil
}

// 获取硬盘序列号
func getDiskID() (string, error) {
    var cmd string
    if runtime.GOOS == "windows" {
        // Windows 上使用 WMIC 命令获取硬盘序列号
        cmd = "wmic diskdrive get SerialNumber"
    } else if runtime.GOOS == "darwin" {
        // macOS 上使用 system_profiler 命令获取硬盘序列号
        cmd = "system_profiler SPSerialATADataType | grep 'Serial Number' | awk '{print $3}'"
    } else {
        // Linux 上使用 lsblk 命令获取硬盘序列号
        cmd = "lsblk -dn -o serial"
    }
    output, err := exec.Command("sh", "-c", cmd).Output()
    if err != nil {
        return "", err
    }
    diskID := ""
    if runtime.GOOS == "windows" {
        // Windows 上硬盘序列号在第二行，包含在输出的第二个单词中
        lines := strings.Split(string(output), "\n")
        if len(lines) > 1 {
            words := strings.Fields(lines[1])
            if len(words) > 1 {
                diskID = words[1]
            }
        }
    } else {
        // macOS 和 Linux 上硬盘序列号在输出的第一行
        lines := strings.Split(string(output), "\n")
        if len(lines) > 0 {
            diskID = lines[0]
        }
    }
    return diskID, nil
}

// 获取网卡 MAC 地址
func getMACAddr() (string, error) {
    ifas, err := net.Interfaces()
    if err != nil {
        return "", err
    }
    macAddr := ""
    // 遍历所有网卡，找到第一个非零 MAC 地址
    for _, ifa := range ifas {
        if ifa.HardwareAddr != nil {
            mac := ifa.HardwareAddr.String()
            if mac != "" && mac != "00:00:00:00:00:00" {
                macAddr = mac
                break
            }
        }
    }
    return macAddr, nil
}

// 将字符串按指定长度进行分组
func splitByLength(str string, length int) []string {
    var result []string
    for i := 0; i < len(str); i += length {
        end := i + length
        if end > len(str) {
            end = len(str)
        }
        result = append(result, str[i:end])
    }
    return result
}
