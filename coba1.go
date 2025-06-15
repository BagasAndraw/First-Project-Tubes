package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Waktu struct {
	JamMasuk  time.Time
	JamKeluar time.Time
}

type Kendaraan struct {
	PlatNomor string
	Jenis     string
	Waktu     Waktu
	Slot      int
}

type SlotParkir struct {
	Nomor  int
	Kosong bool
	Jenis  string
}

const jumlahSlotMotor = 15
const jumlahSlotMobil = 15

var slotParkirMotor [jumlahSlotMotor]SlotParkir
var slotParkirMobil [jumlahSlotMobil]SlotParkir
var kendaraanParkir []Kendaraan
var historiKendaraan []Kendaraan
var scanner = bufio.NewScanner(os.Stdin)

// Inisialisasi slot
func initSlot() {
	for i := 0; i < len(slotParkirMotor); i++ {
		slotParkirMotor[i] = SlotParkir{Nomor: i + 1, Kosong: true, Jenis: "Motor"}
	}
	for i := 0; i < len(slotParkirMobil); i++ {
		slotParkirMobil[i] = SlotParkir{Nomor: i + 1, Kosong: true, Jenis: "Mobil"}
	}
}

// Fungsi input 
func input() string {
	return strings.TrimSpace(scanner.Text())
}

// Fungsi: Hitung durasi parkir
func hitungDurasi(waktu Waktu) time.Duration {
	return waktu.JamKeluar.Sub(waktu.JamMasuk)
}

// Prosedur: Masukkan kendaraan dengan validasi panjang plat nomor dan input nomor slot
func masukkanKendaraan() {
	const maxPlatLength = 4

	fmt.Print("Masukkan plat nomor (4 angka): ")
	scanner.Scan()
	platStr := input()
	platNum, err := strconv.Atoi(platStr)
	if err != nil {
		fmt.Println("Input harus berupa angka!")
		return
	}
	if platNum < 1000 || platNum > 9999 {
		fmt.Println("Plat nomor harus terdiri dari 4 digit angka.")
		return
	}


	fmt.Print("Masukkan jenis kendaraan (Mobil/Motor): ")
	scanner.Scan()
	jenis := strings.ToLower(input())

	var slotNum int

	if jenis == "motor" {
		fmt.Print("Masukkan nomor slot motor yang diinginkan (1-15): ")
		scanner.Scan()
		slotInput := input()
		slotNum, err = strconv.Atoi(slotInput)
		if err != nil || slotNum < 1 || slotNum > len(slotParkirMotor) {
			fmt.Println("Nomor slot motor tidak valid.")
			return
		}
		if !slotParkirMotor[slotNum-1].Kosong {
			fmt.Println("Slot motor sudah terisi.")
			return
		}
		slotParkirMotor[slotNum-1].Kosong = false
	} else if jenis == "mobil" {
		fmt.Print("Masukkan nomor slot mobil yang diinginkan (1-15): ")
		scanner.Scan()
		slotInput := input()
		slotNum, err = strconv.Atoi(slotInput)
		if err != nil || slotNum < 1 || slotNum > len(slotParkirMobil) {
			fmt.Println("Nomor slot mobil tidak valid.")
			return
		}
		if !slotParkirMobil[slotNum-1].Kosong {
			fmt.Println("Slot mobil sudah terisi.")
			return
		}
		slotParkirMobil[slotNum-1].Kosong = false
	} else {
		fmt.Println("Jenis kendaraan tidak valid. Harus 'Mobil' atau 'Motor'.")
		return
	}

	now := time.Now()
	kendaraan := Kendaraan{
		PlatNomor: strconv.Itoa(platNum),
		Jenis:     jenis,
		Waktu:     Waktu{JamMasuk: now},
		Slot:      slotNum,
	}
	kendaraanParkir = append(kendaraanParkir, kendaraan)
	fmt.Println("Kendaraan masuk ke slot:", slotNum)
}

// Prosedur: Keluarkan kendaraan dan simpan histori
func keluarkanKendaraan() {
	fmt.Print("Masukkan plat nomor kendaraan yang keluar: ")
	scanner.Scan()
	plat := input()

	for i := 0; i < len(kendaraanParkir); i++ {
		if kendaraanParkir[i].PlatNomor == plat {
			now := time.Now()
			k := kendaraanParkir[i]
			k.Waktu.JamKeluar = now
			durasi := hitungDurasi(k.Waktu)

			historiKendaraan = append(historiKendaraan, k)

			jenis := strings.ToLower(k.Jenis)
			if jenis == "motor" {
				slotParkirMotor[k.Slot-1].Kosong = true
			} else if jenis == "mobil" {
				slotParkirMobil[k.Slot-1].Kosong = true
			}

			kendaraanParkir = append(kendaraanParkir[:i], kendaraanParkir[i+1:]...)
			fmt.Printf("Kendaraan keluar dari slot: %d\n", k.Slot)
			fmt.Printf("Jenis: %s\n", k.Jenis)
			fmt.Printf("Durasi parkir: %.0f menit\n", durasi.Minutes())
			return
		}
	}
	fmt.Println("Kendaraan tidak ditemukan!")
}

// Sequential Search: Cari kendaraan berdasarkan plat
func CariKendaraanSequential() {
	fmt.Print("Masukkan plat nomor: ")
	scanner.Scan()
	plat := input()

	ditemukan := false
	for _, k := range kendaraanParkir {
		if k.PlatNomor == plat {
			fmt.Printf("Ditemukan: %s (%s), Slot %d\n", k.PlatNomor, k.Jenis, k.Slot)
			ditemukan = true
		}
	}
	if !ditemukan {
		fmt.Println("Kendaraan tidak ditemukan.")
	}
}

// Binary Search: Cari kendaraan berdasarkan jam masuk (HH:MM), dalam rentang
func cariKendaraanBerdasarkanJam() {
	fmt.Print("Masukkan jam mulai (HH:MM): ")
	scanner.Scan()
	startStr := input()
	fmt.Print("Masukkan jam akhir (HH:MM): ")
	scanner.Scan()
	endStr := input()

	parseJamKeMenit := func(s string) int {
		parts := strings.Split(s, ":")
		if len(parts) != 2 {
			return -1
		}
		jam, _ := strconv.Atoi(parts[0])
		menit, _ := strconv.Atoi(parts[1])
		return jam*60 + menit
	}

	startMenit := parseJamKeMenit(startStr)
	endMenit := parseJamKeMenit(endStr)
	if startMenit == -1 || endMenit == -1 || startMenit > endMenit {
		fmt.Println("Format waktu tidak valid atau rentang salah.")
		return
	}

	for i := 1; i < len(kendaraanParkir); i++ {
		key := kendaraanParkir[i]
		j := i - 1
		for j >= 0 && kendaraanParkir[j].Waktu.JamMasuk.After(key.Waktu.JamMasuk) {
			kendaraanParkir[j+1] = kendaraanParkir[j]
			j--
		}
		kendaraanParkir[j+1] = key
	}

	getMenit := func(t time.Time) int {
		return t.Hour()*60 + t.Minute()
	}

	low, high := 0, len(kendaraanParkir)-1
	startIdx := -1
	for low <= high {
		mid := (low + high) / 2
		if getMenit(kendaraanParkir[mid].Waktu.JamMasuk) >= startMenit {
			startIdx = mid
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	low, high = 0, len(kendaraanParkir)-1
	endIdx := -1
	for low <= high {
		mid := (low + high) / 2
		if getMenit(kendaraanParkir[mid].Waktu.JamMasuk) <= endMenit {
			endIdx = mid
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	if startIdx == -1 || endIdx == -1 || startIdx > endIdx {
		fmt.Println("Tidak ada kendaraan dalam rentang waktu tersebut.")
		return
	}

	fmt.Printf("Kendaraan yang masuk antara %s dan %s:\n", startStr, endStr)
	for i := startIdx; i <= endIdx; i++ {
		k := kendaraanParkir[i]
		fmt.Printf("- %s (%s), Slot %d, Masuk: %s\n",
			k.PlatNomor, k.Jenis, k.Slot, k.Waktu.JamMasuk.Format("15:04"))
	}
}

// Sequential search: Tampilkan daftar slot kosong tanpa input apapun
func cariSlotKosong() {
	fmt.Println("Slot parkir motor:")
	for _, slot := range slotParkirMotor {
		if slot.Kosong {
			fmt.Printf("Slot %d \n", slot.Nomor)
		} else {
			fmt.Printf("Slot %d Sudah Terisi\n", slot.Nomor)
		}
	}

	fmt.Println("Slot parkir mobil:")
	for _, slot := range slotParkirMobil {
		if slot.Kosong {
			fmt.Printf("Slot %d \n", slot.Nomor)
		} else {
			fmt.Printf("Slot %d Sudah Terisi\n", slot.Nomor)
		}
	}
}

// Selection Sort: Urutkan histori berdasarkan durasi dan tampilkan
func urutkanKendaraanParkirBerdasarkanDurasi() {
	if len(kendaraanParkir) == 0 {
		fmt.Println("Tidak ada kendaraan yang sedang parkir.")
		return
	}
	for i := 0; i < len(kendaraanParkir); i++ {
		min := i
		for j := i + 1; j < len(kendaraanParkir); j++ {
			durasiJ := time.Since(kendaraanParkir[j].Waktu.JamMasuk)
			durasiMin := time.Since(kendaraanParkir[min].Waktu.JamMasuk)
			if durasiJ < durasiMin {
				min = j
			}
		}
		kendaraanParkir[i], kendaraanParkir[min] = kendaraanParkir[min], kendaraanParkir[i]
	}
	fmt.Println("Kendaraan parkir diurutkan berdasarkan durasi:")
	for _, k := range kendaraanParkir {
		fmt.Printf("- %s (%s), Slot %d, Durasi: %.0f menit\n",
			k.PlatNomor, k.Jenis, k.Slot, time.Since(k.Waktu.JamMasuk).Minutes())
	}
}

// Insertion Sort: Urutkan kendaraan berdasarkan waktu masuk dan tampilkan
func urutkanHistoriBerdasarkanJenisDanJamKeluar() {
	if len(historiKendaraan) == 0 {
		fmt.Println("Belum ada histori kendaraan.")
		return
	}

	var motorList, mobilList []Kendaraan
	for _, k := range historiKendaraan {
		if strings.ToLower(k.Jenis) == "motor" {
			motorList = append(motorList, k)
		} else {
			mobilList = append(mobilList, k)
		}
	}

	insertionSort := func(list []Kendaraan) {
		for i := 1; i < len(list); i++ {
			key := list[i]
			j := i - 1
			for j >= 0 && list[j].Waktu.JamKeluar.After(key.Waktu.JamKeluar) {
				list[j+1] = list[j]
				j--
			}
			list[j+1] = key
		}
	}

	insertionSort(motorList)
	insertionSort(mobilList)

	fmt.Println("Motor: ")
	for _, m := range motorList {
		fmt.Printf("  - %s, Slot %d, Keluar: %s\n", m.PlatNomor, m.Slot, m.Waktu.JamKeluar.Format("15:04:05"))
	}

	fmt.Println("Mobil: ")
	for _, m := range mobilList {
		fmt.Printf("  - %s, Slot %d, Keluar: %s\n", m.PlatNomor, m.Slot, m.Waktu.JamKeluar.Format("15:04:05"))
	}
}

// Menampilkan kendaraan yang sedang parkir
func tampilkanKendaraanParkir() {
	if len(kendaraanParkir) == 0 {
		fmt.Println("Tidak ada kendaraan yang sedang parkir.")
		return
	}
	for _, k := range kendaraanParkir {
		fmt.Printf("- %s (%s), Slot %d, Masuk: %s\n", k.PlatNomor, k.Jenis, k.Slot, k.Waktu.JamMasuk.Format("15:04:05"))
	}
}

// Menampilkan histori kendaraan
func tampilkanHistori() {
	if len(historiKendaraan) == 0 {
		fmt.Println("Belum ada histori kendaraan.")
		return
	}
	for _, h := range historiKendaraan {
		durasi := hitungDurasi(h.Waktu)
		fmt.Printf("- %s (%s), Slot %d, Durasi: %.0f menit\n", h.PlatNomor, h.Jenis, h.Slot, durasi.Minutes())
	}
}

func main() {
	initSlot()
	for {
		fmt.Println("\n===== MENU PARKIR =====")
		fmt.Println("1. Masukkan Kendaraan")
		fmt.Println("2. Keluarkan Kendaraan")
		fmt.Println("3. Cari Kendaraan (Sequential Search)")
		fmt.Println("4. Cari Kendaraan Berdasarkan Waktu (Binary Search)")
		fmt.Println("5. Cari Slot Kosong (Sequential Search)")
		fmt.Println("6. Tampilkan Kendaraan yang Parkir")
		fmt.Println("7. Tampilkan Histori Kendaraan")
		fmt.Println("8. Urutkan Riwayat Kendaraan Berdasarkan Durasi Parkir (Selection Sort)")
		fmt.Println("9. Urutkan Histori Berdasarkan Jenis dan Waktu Keluar (Insertion Sort)")
		fmt.Println("0. Keluar")

		fmt.Print("Pilih menu: ")
		scanner.Scan()
		pilihan := input()

		switch pilihan {
		case "1":
			masukkanKendaraan()
		case "2":
			keluarkanKendaraan()
		case "3":
			CariKendaraanSequential()
		case "4":
			cariKendaraanBerdasarkanJam()
		case "5":
			cariSlotKosong()
		case "6":
			tampilkanKendaraanParkir()
		case "7":
			tampilkanHistori()
		case "8":
			urutkanKendaraanParkirBerdasarkanDurasi()
		case "9":
			urutkanHistoriBerdasarkanJenisDanJamKeluar()
		case "0":
			fmt.Println("Terima kasih telah menggunakan sistem parkir!")
			return
		default:
			fmt.Println("Pilihan tidak valid.")
		}
	}
}