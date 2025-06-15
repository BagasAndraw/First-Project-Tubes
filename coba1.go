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

const jumlahSlot = 15

var (
	slotMotor        = [jumlahSlot]SlotParkir{}
	slotMobil        = [jumlahSlot]SlotParkir{}
	kendaraanParkir  []Kendaraan
	historiKendaraan []Kendaraan
	scanner          = bufio.NewScanner(os.Stdin)
)

//Fungsi Prosedur
func initSlot() {
	for i := 0; i < jumlahSlot; i++ {
		slotMotor[i] = SlotParkir{Nomor: i + 1, Kosong: true, Jenis: "Motor"}
		slotMobil[i] = SlotParkir{Nomor: i + 1, Kosong: true, Jenis: "Mobil"}
	}
}

//Fungsi
func input(prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

//Fungsi
func hitungDurasi(w Waktu) time.Duration {
	return w.JamKeluar.Sub(w.JamMasuk)
}

//Fungsi
func cariSlotKosong(jenis string) []int {
	slots := []int{}
	arr := slotMotor[:]
	if jenis == "mobil" {
		arr = slotMobil[:]
	}
	for _, s := range arr {
		if s.Kosong {
			slots = append(slots, s.Nomor)
		}
	}
	return slots
}

//Fungsi
func cariSlotDanIsi(jenis string, nomor int) bool {
	if jenis == "motor" && nomor >= 1 && nomor <= jumlahSlot && slotMotor[nomor-1].Kosong {
		slotMotor[nomor-1].Kosong = false
		return true
	}
	if jenis == "mobil" && nomor >= 1 && nomor <= jumlahSlot && slotMobil[nomor-1].Kosong {
		slotMobil[nomor-1].Kosong = false
		return true
	}
	return false
}

//Fungsi Prosedur
func kosongkanSlot(jenis string, nomor int) {
	if jenis == "motor" {
		slotMotor[nomor-1].Kosong = true
	} else {
		slotMobil[nomor-1].Kosong = true
	}
}

//Fungsi Prosedur
func masukkanKendaraanIO() {
	plat := input("Plat nomor (maks 4 huruf/angka): ")
	if len(plat) > 4 {
		fmt.Println("Plat terlalu panjang.")
		return
	}
	jenis := strings.ToLower(input("Jenis kendaraan (motor/mobil): "))
	if jenis != "motor" && jenis != "mobil" {
		fmt.Println("Jenis tidak valid.")
		return
	}
	slotInput := input(fmt.Sprintf("Pilih slot %s (1-%d): ", jenis, jumlahSlot))
	slotNum, err := strconv.Atoi(slotInput)
	if err != nil || !cariSlotDanIsi(jenis, slotNum) {
		fmt.Println("Slot tidak valid atau sudah terisi.")
		return
	}

	now := time.Now()
	kendaraanParkir = append(kendaraanParkir, Kendaraan{
		PlatNomor: plat,
		Jenis:     jenis,
		Slot:      slotNum,
		Waktu:     Waktu{JamMasuk: now},
	})
	fmt.Printf("Kendaraan %s masuk slot %d.\n", plat, slotNum)
}

//Fungsi Prosedur
func keluarkanKendaraanIO() {
	plat := input("Plat nomor kendaraan keluar: ")
	for i, k := range kendaraanParkir {
		if k.PlatNomor == plat {
			k.Waktu.JamKeluar = time.Now()
			durasi := hitungDurasi(k.Waktu)
			historiKendaraan = append(historiKendaraan, k)
			kendaraanParkir = append(kendaraanParkir[:i], kendaraanParkir[i+1:]...)
			kosongkanSlot(k.Jenis, k.Slot)

			fmt.Printf("Kendaraan %s keluar dari slot %d\n", k.PlatNomor, k.Slot)
			fmt.Printf("Durasi parkir: %.0f menit\n", durasi.Minutes())
			return
		}
	}
	fmt.Println("Kendaraan tidak ditemukan.")
}

//Fungsi Prosedur
func tampilkanKendaraanParkirIO() {
	if len(kendaraanParkir) == 0 {
		fmt.Println("Tidak ada kendaraan.")
		return
	}
	for _, k := range kendaraanParkir {
		fmt.Printf("- %s (%s), Slot %d, Masuk: %s\n", k.PlatNomor, k.Jenis, k.Slot, k.Waktu.JamMasuk.Format("15:04"))
	}
}

//Fungsi Prosedur
func tampilkanHistoriIO() {
	if len(historiKendaraan) == 0 {
		fmt.Println("Belum ada histori.")
		return
	}
	for _, h := range historiKendaraan {
		durasi := hitungDurasi(h.Waktu)
		fmt.Printf("- %s (%s), Slot %d, Keluar: %s, Durasi: %.0f menit\n",
			h.PlatNomor, h.Jenis, h.Slot, h.Waktu.JamKeluar.Format("15:04"), durasi.Minutes())
	}
}

//Fungsi Prosedur
func cariKendaraanSequentialIO() {
	plat := input("Plat yang dicari: ")
	found := false
	for _, k := range kendaraanParkir {
		if k.PlatNomor == plat {
			fmt.Printf("Ditemukan: %s (%s), Slot %d\n", k.PlatNomor, k.Jenis, k.Slot)
			found = true
		}
	}
	if !found {
		fmt.Println("Tidak ditemukan.")
	}
}

//Fungsi Prosedur
func urutkanKendaraanDurasiIO() {
	if len(kendaraanParkir) == 0 {
		fmt.Println("Tidak ada kendaraan.")
		return
	}
	for i := 0; i < len(kendaraanParkir); i++ {
		min := i
		for j := i + 1; j < len(kendaraanParkir); j++ {
			if time.Since(kendaraanParkir[j].Waktu.JamMasuk) < time.Since(kendaraanParkir[min].Waktu.JamMasuk) {
				min = j
			}
		}
		kendaraanParkir[i], kendaraanParkir[min] = kendaraanParkir[min], kendaraanParkir[i]
	}
	tampilkanKendaraanParkirIO()
}

//Fungsi Prosedur
func urutkanHistoriIO() {
	if len(historiKendaraan) == 0 {
		fmt.Println("Belum ada histori.")
		return
	}

	motor, mobil := []Kendaraan{}, []Kendaraan{}
	for _, k := range historiKendaraan {
		if k.Jenis == "motor" {
			motor = append(motor, k)
		} else {
			mobil = append(mobil, k)
		}
	}
	sortByKeluar := func(list []Kendaraan) {
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
	sortByKeluar(motor)
	sortByKeluar(mobil)

	fmt.Println("Motor:")
	for _, m := range motor {
		fmt.Printf("  - %s, Keluar: %s\n", m.PlatNomor, m.Waktu.JamKeluar.Format("15:04"))
	}
	fmt.Println("Mobil:")
	for _, m := range mobil {
		fmt.Printf("  - %s, Keluar: %s\n", m.PlatNomor, m.Waktu.JamKeluar.Format("15:04"))
	}
}

//Fungsi Prosedur
func cariSlotKosongIO() {
	fmt.Println("Slot Motor Kosong:")
	for _, s := range slotMotor {
		status := "Terisi"
		if s.Kosong {
			status = "Kosong"
		}
		fmt.Printf("- Slot %d: %s\n", s.Nomor, status)
	}
	fmt.Println("Slot Mobil Kosong:")
	for _, s := range slotMobil {
		status := "Terisi"
		if s.Kosong {
			status = "Kosong"
		}
		fmt.Printf("- Slot %d: %s\n", s.Nomor, status)
	}
}

func main() {
	initSlot()
	for {
		fmt.Println("\n===== MENU PARKIR =====")
		fmt.Println("1. Masukkan Kendaraan")
		fmt.Println("2. Keluarkan Kendaraan")
		fmt.Println("3. Cari Kendaraan (Sequential Search)")
		fmt.Println("4. Tampilkan Slot Kosong")
		fmt.Println("5. Tampilkan Kendaraan Parkir")
		fmt.Println("6. Tampilkan Histori")
		fmt.Println("7. Urutkan Berdasarkan Durasi Parkir")
		fmt.Println("8. Urutkan Histori per Jenis & Waktu Keluar")
		fmt.Println("0. Keluar")
		pil := input("Pilih menu: ")
		switch pil {
		case "1":
			masukkanKendaraanIO()
		case "2":
			keluarkanKendaraanIO()
		case "3":
			cariKendaraanSequentialIO()
		case "4":
			cariSlotKosongIO()
		case "5":
			tampilkanKendaraanParkirIO()
		case "6":
			tampilkanHistoriIO()
		case "7":
			urutkanKendaraanDurasiIO()
		case "8":
			urutkanHistoriIO()
		case "0":
			fmt.Println("Terima kasih.")
			return
		default:
			fmt.Println("Menu tidak valid.")
		}
	}
}