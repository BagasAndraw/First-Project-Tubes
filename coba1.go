package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Tipe Bentukan
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
}

// Variabel global
var slotParkir [100]SlotParkir
var kendaraanParkir []Kendaraan
var historiKendaraan []Kendaraan
var scanner = bufio.NewScanner(os.Stdin)

// Inisialisasi slot
func initSlot() {
	for i := 0; i < len(slotParkir); i++ {
		slotParkir[i] = SlotParkir{Nomor: i + 1, Kosong: true}
	}
}

// Fungsi input
func input(prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

// Fungsi: Hitung durasi parkir
func hitungDurasi(waktu Waktu) time.Duration {
	return waktu.JamKeluar.Sub(waktu.JamMasuk)
}

// Prosedur: Masukkan kendaraan dengan input nomor slot (tanpa menampilkan slot kosong)
func masukkanKendaraan() {
	plat := input("Masukkan plat nomor: ")
	jenis := input("Masukkan jenis kendaraan (Mobil/Motor): ")
	slotInput := input("Masukkan nomor slot yang diinginkan: ")
	slotNum, err := strconv.Atoi(slotInput)
	if err != nil || slotNum < 1 || slotNum > len(slotParkir) {
		fmt.Println("Nomor slot tidak valid.")
		return
	}
	if !slotParkir[slotNum-1].Kosong {
		fmt.Println("Slot sudah terisi.")
		return
	}

	now := time.Now()
	slotParkir[slotNum-1].Kosong = false
	kendaraan := Kendaraan{
		PlatNomor: plat,
		Jenis:     jenis,
		Waktu:     Waktu{JamMasuk: now},
		Slot:      slotNum,
	}
	kendaraanParkir = append(kendaraanParkir, kendaraan)
	fmt.Println("Kendaraan masuk ke slot:", slotNum)
}

// Prosedur: Keluarkan kendaraan dan simpan histori
func keluarkanKendaraan() {
	plat := input("Masukkan plat nomor kendaraan yang keluar: ")
	for i, k := range kendaraanParkir {
		if k.PlatNomor == plat {
			now := time.Now()
			k.Waktu.JamKeluar = now
			historiKendaraan = append(historiKendaraan, k)
			slotParkir[k.Slot-1].Kosong = true
			kendaraanParkir = append(kendaraanParkir[:i], kendaraanParkir[i+1:]...)
			fmt.Println("Kendaraan keluar dari slot:", k.Slot)
			return
		}
	}
	fmt.Println("Kendaraan tidak ditemukan!")
}

// Linear Search: Cari kendaraan berdasarkan plat
func CariKendaraanSequential() {
	plat := input("Masukkan plat nomor: ")
	ditemukan := false

	for i := 0; i < len(kendaraanParkir); i++ {
		if kendaraanParkir[i].PlatNomor == plat {
			fmt.Printf("Ditemukan: %s (%s), Slot %d\n", kendaraanParkir[i].PlatNomor, kendaraanParkir[i].Jenis, kendaraanParkir[i].Slot)
			ditemukan = true
		}
	}

	if !ditemukan {
		fmt.Println("Kendaraan tidak ditemukan.")
	}
}

// Opsi 4: Tampilkan daftar slot kosong tanpa input apapun
func cariSlotKosong() {
	fmt.Println("Daftar slot kosong:")
	var kosong []SlotParkir
	for _, slot := range slotParkir {
		if slot.Kosong {
			kosong = append(kosong, slot)
		}
	}
	if len(kosong) == 0 {
		fmt.Println("Tidak ada slot kosong.")
		return
	}
	for _, s := range kosong {
		fmt.Printf("- Slot %d\n", s.Nomor)
	}
}

// Selection Sort: Urutkan histori berdasarkan durasi dan tampilkan
func selectionSortDurasi() {
	for i := 0; i < len(historiKendaraan); i++ {
		min := i
		for j := i + 1; j < len(historiKendaraan); j++ {
			if hitungDurasi(historiKendaraan[j].Waktu) < hitungDurasi(historiKendaraan[min].Waktu) {
				min = j
			}
		}
		historiKendaraan[i], historiKendaraan[min] = historiKendaraan[min], historiKendaraan[i]
	}
	fmt.Println("Histori diurutkan berdasarkan durasi parkir (Selection Sort).")
	tampilkanHistori()
}

// Binary Sort: Urutkan kendaraan berdasarkan waktu masuk dan tampilkan
func urutkanDanTampilkanBerdasarkanJenis() {
	if len(kendaraanParkir) == 0 {
		fmt.Println("Tidak ada kendaraan yang sedang parkir.")
		return
	}

	// Pisahkan kendaraan berdasarkan jenis
	var motorList []Kendaraan
	var mobilList []Kendaraan

	for _, k := range kendaraanParkir {
		jenisLower := strings.ToLower(k.Jenis)
		if jenisLower == "motor" {
			motorList = append(motorList, k)
		} else if jenisLower == "mobil" {
			mobilList = append(mobilList, k)
		} else {
			// Jika jenis lain atau input salah, bisa masuk ke daftar mobil misalnya
			mobilList = append(mobilList, k)
		}
	}

	// Tampilkan hasil
	fmt.Println("Kendaraan diurutkan berdasarkan jenis:")

	if len(motorList) > 0 {
		fmt.Println("- Motor:")
		for _, m := range motorList {
			fmt.Printf("  - %s, Slot %d, Masuk: %s\n",
				m.PlatNomor, m.Slot, m.Waktu.JamMasuk.Format("15:04:05"))
		}
	} else {
		fmt.Println("- Motor: Tidak ada")
	}

	if len(mobilList) > 0 {
		fmt.Println("- Mobil:")
		for _, m := range mobilList {
			fmt.Printf("  - %s, Slot %d, Masuk: %s\n",
				m.PlatNomor, m.Slot, m.Waktu.JamMasuk.Format("15:04:05"))
		}
	} else {
		fmt.Println("- Mobil: Tidak ada")
	}
}

// Menampilkan kendaraan yang sedang parkir
func tampilkanKendaraanParkir() {
	if len(kendaraanParkir) == 0 {
		fmt.Println("Tidak ada kendaraan yang sedang parkir.")
		return
	}
	fmt.Println("Daftar kendaraan parkir:")
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
	fmt.Println("Histori kendaraan:")
	for _, h := range historiKendaraan {
		durasi := hitungDurasi(h.Waktu)
		fmt.Printf("- %s (%s), Slot %d, Durasi: %.0f menit\n",
			h.PlatNomor, h.Jenis, h.Slot, durasi.Minutes())
	}
}

func main() {
	initSlot()
	for {
		fmt.Println("\n===== MENU PARKIR =====")
		fmt.Println("1. Masukkan Kendaraan")
		fmt.Println("2. Keluarkan Kendaraan")
		fmt.Println("3. Cari Kendaraan (Sequential Search)")
		fmt.Println("4. Cari Slot Kosong (Sequential Search)")
		fmt.Println("5. Tampilkan Kendaraan yang Parkir")
		fmt.Println("6. Tampilkan Histori Kendaraan")
		fmt.Println("7. Urutkan Histori Berdasarkan Durasi (Selection Sort)")
		fmt.Println("8. Urutkan Kendaraan Berdasarkan Waktu Masuk (Binary Sort)")
		fmt.Println("0. Keluar")
		pilihan := input("Pilih menu: ")

		switch pilihan {
		case "1":
			masukkanKendaraan()
		case "2":
			keluarkanKendaraan()
		case "3":
			CariKendaraanSequential()
		case "4":
			cariSlotKosong()
		case "5":
			tampilkanKendaraanParkir()
		case "6":
			tampilkanHistori()
		case "7":
			selectionSortDurasi()
		case "8":
			urutkanDanTampilkanBerdasarkanJenis()
		case "0":
			fmt.Println("Terima kasih telah menggunakan sistem parkir!")
			return
		default:
			fmt.Println("Pilihan tidak valid.")
		}
	}
}
