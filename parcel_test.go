package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
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
	db,err:=sql.Open("sqlite","tracker.db")
	require.NoError(t,err)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id,err:=store.Add(parcel)
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	require.NoError(t,err)
	require.Greater(t,id,0)


	// get
	p,err:=store.Get(id)
	require.NoError(t,err)
	parcel.Number = id
	assert.Equal(t, p, parcel, "Parcel does not match.")

	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel

	// delete
	err=store.Delete(id)
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	require.NoError(t,err)
	// проверьте, что посылку больше нельзя получить из БД
	_,err=store.Get(id)
	require.Error(t, err, "Parcel should not exist in the DB after delete command")
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db,err:=sql.Open("sqlite","tracker.db")
	require.NoError(t,err)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id,err:=store.Add(parcel)
	require.NoError(t,err)
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора

	// set address
	// обновите адрес, убедитесь в отсутствии ошибки
	newAddress := "new test address"
	err=store.SetAddress(id,newAddress)
	// check
	require.NoError(t,err)
	// получите добавленную посылку и убедитесь, что адрес обновился
	p,err:=store.Get(id)
	require.NoError(t,err)
	assert.Equal(t, p.Address, newAddress, "Adress does not match.")
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db,err:=sql.Open("sqlite","tracker.db")
	require.NoError(t,err)
	store := NewParcelStore(db)
	parcel := getTestParcel()


	// add
	id,err:=store.Add(parcel)
	require.NoError(t,err)
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора

	// set status
	err=store.SetStatus(id,ParcelStatusSent)
	require.NoError(t,err)
	// обновите статус, убедитесь в отсутствии ошибки
	p,err:=store.Get(id)
	require.NoError(t,err)
	assert.Equal(t, ParcelStatusSent, p.Status, "Status does not match")
	// check
	// получите добавленную посылку и убедитесь, что статус обновился
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db,err:=sql.Open("sqlite","tracker.db")
	require.NoError(t,err)
	store := NewParcelStore(db)


	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t,err)
		// обновляем идентификатор добавленной у посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	// убедитесь в отсутствии ошибки
	require.NoError(t,err)
	// убедитесь, что количество полученных посылок совпадает с количеством добавленных
	assert.Equal(t,len(storedParcels),len(parcels),"Amount of parcels retrieved from DB does not match amount added")
	for _, parcel := range storedParcels {
		expectedParcel, ok := parcelMap[parcel.Number]
	require.True(t, ok, "Unexpected parcel found with ID: %d", parcel.Number)
    assert.Equal(t, expectedParcel, parcel, "Parcel does not match for parcel ID: %d", parcel.Number)
}
}
