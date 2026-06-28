# Contoh Perhitungan Jarak Tempuh Drone

## Skenario

Diasumsikan:

* Ukuran estate: **20 x 20**
* Total plot: **400 plot**

Terdapat dua pohon:

| Pohon | X | Y | Tinggi |
| ----- | - | - | ------ |
| A     | 2 | 2 | 5 m    |
| B     | 2 | 8 | 6 m    |

Seluruh plot lainnya kosong.

---

## Aturan Ketinggian

### Plot Kosong

Ketinggian drone saat melewati plot kosong:

```text
1 meter
```

### Plot Berisi Pohon

Ketinggian drone saat melewati pohon:

```text
tinggi_pohon + 1 meter
```

Sehingga:

| Lokasi      | Ketinggian |
| ----------- | ---------- |
| Plot kosong | 1 m        |
| (2,2)       | 6 m        |
| (2,8)       | 7 m        |

---

# Perhitungan Jarak Horizontal

Estate berukuran:

```text
20 × 20 = 400 plot
```

Drone mengunjungi setiap plot tepat satu kali.

Jumlah perpindahan horizontal:

```text
400 - 1 = 399 perpindahan
```

Setiap perpindahan horizontal bernilai:

```text
10 meter
```

Maka:

```text
399 × 10
= 3990 meter
```

**Total jarak horizontal = 3990 meter**

---

# Perhitungan Jarak Vertikal

## 1. Lepas Landas Awal

Drone berangkat dari tanah:

```text
0 m → 1 m
```

Jarak vertikal:

```text
1 meter
```

---

## 2. Pohon di (2,2)

Urutan kunjungan:

```text
Plot ke-39
```

### Sebelum memasuki pohon

Posisi sebelumnya:

```text
(3,2)
```

Ketinggian:

```text
1 m
```

Naik ke:

```text
6 m
```

Jarak:

```text
6 - 1 = 5 m
```

### Setelah melewati pohon

Posisi berikutnya:

```text
(1,2)
```

Turun dari:

```text
6 m → 1 m
```

Jarak:

```text
5 m
```

### Kontribusi Pohon (2,2)

```text
5 + 5 = 10 meter
```

---

## 3. Pohon di (2,8)

Urutan kunjungan:

```text
Plot ke-159
```

### Sebelum memasuki pohon

Posisi sebelumnya:

```text
(3,8)
```

Ketinggian:

```text
1 m
```

Naik ke:

```text
7 m
```

Jarak:

```text
7 - 1 = 6 m
```

### Setelah melewati pohon

Posisi berikutnya:

```text
(1,8)
```

Turun dari:

```text
7 m → 1 m
```

Jarak:

```text
6 m
```

### Kontribusi Pohon (2,8)

```text
6 + 6 = 12 meter
```

---

## 4. Pendaratan Akhir

Plot terakhir:

```text
(20,20)
```

Ketinggian plot:

```text
1 m
```

Drone turun ke tanah:

```text
1 m → 0 m
```

Jarak:

```text
1 meter
```

---

# Total Jarak Vertikal

```text
Lepas landas awal = 1 m
Pohon (2,2)       = 10 m
Pohon (2,8)       = 12 m
Pendaratan akhir  = 1 m
--------------------------------
Total vertikal    = 24 m
```

---

# Total Jarak Tempuh Drone

```text
Horizontal = 3990 m
Vertikal   =   24 m
-------------------
Total      = 4014 m
```

---

# Hasil API

```json
{
  "distance": 4014
}
```

## Kesimpulan

Jarak total yang ditempuh drone adalah:

```text
4014 meter
```
