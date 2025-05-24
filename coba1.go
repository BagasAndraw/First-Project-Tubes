package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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

// Prosedur: Masukkan kendaraan
func masukkanKendaraan() {
	plat := input("Masukkan plat nomor: ")
	jenis := input("Masukkan jenis kendaraan (Mobil/Motor): ")
	now := time.Now()

	for i := 0; i < len(slotParkir); i++ {
		if slotParkir[i].Kosong {
			slotParkir[i].Kosong = false
			kendaraan := Kendaraan{
				PlatNomor: plat,
				Jenis:     jenis,
				Waktu:     Waktu{JamMasuk: now},
				Slot:      slotParkir[i].Nomor,
			}
			kendaraanParkir = append(kendaraanParkir, kendaraan)
			fmt.Println("✅ Kendaraan masuk ke slot:", slotParkir[i].Nomor)
			return
		}
	}
	fmt.Println("❌ Slot penuh!")
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
			fmt.Println("✅ Kendaraan keluar dari slot:", k.Slot)
			return
		}
	}
	fmt.Println("❌ Kendaraan tidak ditemukan!")
}

// Linear Search: Cari kendaraan berdasarkan plat
func cariKendaraanLinear() {
	plat := input("Masukkan plat nomor: ")
	for _, k := range kendaraanParkir {
		if k.PlatNomor == plat {
			fmt.Printf("✅ Ditemukan: %s (%s), Slot %d\n", k.PlatNomor, k.Jenis, k.Slot)
			return
		}
	}
	fmt.Println("❌ Kendaraan tidak ditemukan.")
}

// Binary Search: Cari slot kosong berdasarkan nomor slot
func cariSlotKosongBinary() {
	sort.Slice(slotParkir[:], func(i, j int) bool {
		return slotParkir[i].Nomor < slotParkir[j].Nomor
	})
	slotInput := input("Masukkan nomor slot yang dicari: ")
	slotNum, err := strconv.Atoi(slotInput)
	if err != nil {
		fmt.Println("❌ Nomor tidak valid.")
		return
	}
	low, high := 0, len(slotParkir)-1
	for low <= high {
		mid := (low + high) / 2
		if slotParkir[mid].Nomor == slotNum {
			if slotParkir[mid].Kosong {
				fmt.Println("✅ Slot", slotNum, "kosong.")
			} else {
				fmt.Println("❌ Slot", slotNum, "sudah terisi.")
			}
			return
		} else if slotParkir[mid].Nomor < slotNum {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	fmt.Println("❌ Slot tidak ditemukan.")
}

// Selection Sort: Urutkan histori berdasarkan durasi
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
	fmt.Println("✅ Histori diurutkan berdasarkan durasi parkir (Selection Sort).")
}

// Binary Sort: Urutkan kendaraan berdasarkan waktu masuk
func binarySortByWaktuMasuk() {
	sort.Slice(kendaraanParkir, func(i, j int) bool {
		return kendaraanParkir[i].Waktu.JamMasuk.Before(kendaraanParkir[j].Waktu.JamMasuk)
	})
	fmt.Println("✅ Kendaraan diurutkan berdasarkan waktu masuk (Binary Sort).")
}

// Menampilkan kendaraan yang sedang parkir
func tampilkanKendaraanParkir() {
	if len(kendaraanParkir) == 0 {
		fmt.Println("🚫 Tidak ada kendaraan yang sedang parkir.")
		return
	}
	fmt.Println("🚗 Daftar kendaraan parkir:")
	for _, k := range kendaraanParkir {
		fmt.Printf("- %s (%s), Slot %d, Masuk: %s\n", k.PlatNomor, k.Jenis, k.Slot, k.Waktu.JamMasuk.Format("15:04:05"))
	}
}

// Menampilkan histori kendaraan
func tampilkanHistori() {
	if len(historiKendaraan) == 0 {
		fmt.Println("📭 Belum ada histori kendaraan.")
		return
	}
	fmt.Println("📚 Histori kendaraan:")
	for _, h := range historiKendaraan {
		durasi := hitungDurasi(h.Waktu)
		fmt.Printf("- %s (%s), Slot %d, Durasi: %.0f menit\n",
			h.PlatNomor, h.Jenis, h.Slot, durasi.Minutes())
	}
}

// Main
func main() {
	initSlot()
	for {
		fmt.Println("\n===== MENU PARKIR =====")
		fmt.Println("1. Masukkan Kendaraan")
		fmt.Println("2. Keluarkan Kendaraan")
		fmt.Println("3. Cari Kendaraan (Squential Search)")
		fmt.Println("4. Cari Slot Kosong (Binary Search)")
		fmt.Println("5. Tampilkan Kendaraan yang Parkir")
		fmt.Println("6. Tampilkan Histori Kendaraan")
		fmt.Println("7. Urutkan Histori Berdasarkan Durasi (Selection Sort)")
		fmt.Println("8. Urutkan Kendaraan Berdasarkan Waktu Masuk (Quick Sort)")
		fmt.Println("0. Keluar")
		fmt.Println()

		pilihan := input("Pilih menu: ")

		switch pilihan {
		case "1":
			masukkanKendaraan()
		case "2":
			keluarkanKendaraan()
		case "3":
			cariKendaraanLinear()
		case "4":
			cariSlotKosongBinary()
		case "5":
			tampilkanKendaraanParkir()
		case "6":
			tampilkanHistori()
		case "7":
			selectionSortDurasi()
		case "8":
			binarySortByWaktuMasuk()
		case "0":
			fmt.Println("👋 Terima kasih telah menggunakan sistem parkir!")
			return
		default:
			fmt.Println("❌ Pilihan tidak valid.")
		}
	}
}
