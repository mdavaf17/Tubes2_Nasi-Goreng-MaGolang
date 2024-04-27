# Tubes2_Nasi-Goreng-MaGolang
Tugas Besar 2 IF2211 Strategi Algoritma Semester II tahun 2023/2024 Pemanfaatan Algoritma IDS dan BFS dalam Permainan WikiRace

## Iterative Deepening Search (IDS)
1. Menerima URL awal dan URL tujuan
2. Menganggap URL awal, URL tujuan, dan kumpulan URL yang akan dikunjungi sebagai node
3. Dilakukan iterasi untuk setiap kedalaman maksimal
4. Pada iterasi pertama, ketika kedalaman maksimal sama dengan nol, dilakukan pengecekan apakah URL awal sudah sama dengan URL tujuan, jika sama maka penelusuran berhenti
5. Pada iterasi kedua, kedalaman maksimal bertambah satu, dilakukan scrapping untuk mendapatkan semua URL pada laman tersebut.
6. Proses dilanjutkan dengan Iterasi setiap URL pada laman yang ditemukan dan melakukan rekursif dengan elemen iterasi sekarang sebagai URL. Jika URL awal dan URL tujuan sudah sama maka penelusuran dihentikan dan tidak perlu dilanjutkan ke iterasi berikutnya dengan kedalaman yang berbeda. Jika belum sama, dilakukan pengecekan apakah kedalaman sekarang sudah lebih dari atau sama dengan dengan kedalaman maksimal. Apabila sudah, rekursif tidak akan dilakukan, kemudian berlanjut ke iterasi berikutnya dengan kedalaman maksimal yang ditambah dengan satu. Apabila kedalaman sekarang belum lebih dari atau sama dengan kedalaman maksimal maka proses rekursif akan dilanjutkan
7. Selama penelusuran jalur yang sudah dikunjungi disimpan pada suatu array, apabila jalur tersebut tidak mengarah pada URL tujuan, maka akan di pop atau dihapus dari array. Jalur akhir yang tetap berada pada array adalah jalur yang mengarah pada URL tujuan.

## Breadth First Search (BFS)
1. Algoritma menerima masukan berupa URL awal dan URL tujuan
2. Buat sebuah sebuah array “list” yang dapat di-iterasi menggunakan sebuah iterator “i”, beserta counter untuk mencatat jumlah artikel yang diperiksa. “List” merupakan gabungan dari queue berisi URL yang akan dikunjungi dan array visited berisi URL yang telah dikunjungi.
3. Periksa URL yang berada di depan queue (list[i]) dan increment counter. Jika URL adalah URL tujuan, pencarian selesai
4. Jika bukan, cari URL lain dalam URL tersebut yang menuju ke sebuah wiki page.
5. Jika ditemukan URL yang merupakan URL tujuan, pencarian selesain.
6. Jika tidak, masukkan semua URL tersebut ke dalam list. Tambahkan juga simpul-simpul pada tree dengan isi simpul tersebut berupa URL lain yang ditemukan dan parent dari simpul tersebut berupa URL yang sedang diperiksa sekarang.
7. Increment iterator “i” dan ulangi langkah 3-6 sampai pencarian selesai atau list sudah ditelusuri sepenuhnya.
8. Cari rute menuju URL tujuan dengan menelusuri tree. Jika solusi tidak ditemukan, rute kosong.
9. Kembalikan rute yang diperoleh beserta counter.

## Requirements
* [Docker](https://docs.docker.com/get-docker/)

## How To Setup
1. Move to src/ folder
```
cd src
```
2. Build the docker image
```
docker build -t tubes2_nasi_goreng_magolang .
```
3. Run the docker image as a conainer
```
docker run -p 8030:8030 tubes2_nasi_goreng_magolang
```
4. Open your web browser at [http://localhost:8030](http://localhost:8030) or [http://127.0.0.1/:8030](http://127.0.0.1/:8030)
## Author 

### Group 30 | Nasi Goreng MaGolang

| No. | Name                           | NIM |
|-----|--------------------------------|------------|
| 1.  | Ariel Herfrison                | 13522002   |
| 2.  | Panjri Sri Kuncara             | 13522028   |
| 3.  | Muhammad Dava Fathurrahman    | 13522114   |

