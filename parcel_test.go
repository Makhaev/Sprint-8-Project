package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // настройте подключение к БД
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)
	parcel.Number = id
	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel

	got, err := store.Get(id)
	require.NoError(t, err)

	require.Equal(t, parcel.Number, got.Number)
	require.Equal(t, parcel.Client, got.Client)
	require.Equal(t, parcel.Status, got.Status)
	require.Equal(t, parcel.Address, got.Address)
	require.Equal(t, parcel.CreatedAt, got.CreatedAt)

	// === DELETE ===
	err = store.Delete(id)
	require.NoError(t, err)

	// === CHECK DELETE ===
	_, err = store.Get(id)
	require.Error(t, err) // Должна быть ошибка, так как записи больше нет
	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	// проверьте, что посылку больше нельзя получить из БД
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") // подключение к БД
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// === ADD ===
	id, err := store.Add(parcel) // добавляем посылку
	require.NoError(t, err)
	require.NotZero(t, id)

	// обновляем номер посылки
	parcel.Number = id

	// === SET ADDRESS ===
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress) // обновляем адрес
	require.NoError(t, err)

	// === CHECK ===
	got, err := store.Get(id) // получаем посылку из БД
	require.NoError(t, err)
	require.Equal(t, newAddress, got.Address) // проверяем, что адрес обновлён
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// === PREPARE ===
	db, err := sql.Open("sqlite", "tracker.db") // подключение к БД
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// === ADD ===
	id, err := store.Add(parcel) // добавляем посылку
	require.NoError(t, err)
	require.NotZero(t, id)

	// обновляем номер посылки
	parcel.Number = id

	// === SET STATUS ===
	newStatus := ParcelStatusSent // например, обновим статус на "sent"
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err)

	// === CHECK ===
	got, err := store.Get(id) // получаем посылку из БД
	require.NoError(t, err)
	require.Equal(t, newStatus, got.Status) // проверяем, что статус обновился
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// === PREPARE ===
	db, err := sql.Open("sqlite", "tracker.db") // подключение к БД
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	for i := range parcels {
		parcels[i].Client = client
	}

	// создаём map для проверки результата
	parcelMap := map[int]Parcel{}

	// === ADD ===
	for i := range parcels {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotZero(t, id)

		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	// === GET BY CLIENT ===
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Len(t, storedParcels, len(parcels)) // сравниваем количество

	// === CHECK ===
	for _, p := range storedParcels {
		original, ok := parcelMap[p.Number]
		require.True(t, ok, "Посылка с номером %d не найдена в parcelMap", p.Number)

		require.Equal(t, original.Client, p.Client)
		require.Equal(t, original.Status, p.Status)
		require.Equal(t, original.Address, p.Address)
		require.Equal(t, original.CreatedAt, p.CreatedAt)
	}
}
