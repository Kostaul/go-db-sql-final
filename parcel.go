package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	query := "INSERT INTO parcel (client, status, address, created_at) VALUES (?,?,?,?)"

	result,err:= s.db.Exec(query,p.Client, p.Status, p.Address, p.CreatedAt)
	if err!=nil{
		return 0,err
	}
	id, err := result.LastInsertId()
	if err!=nil{
		return 0,err
	}
	// верните идентификатор последней добавленной записи
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	query:= "SELECT client,status,address,created_at FROM parcel WHERE number=?"
	p := Parcel{}

	err := s.db.QueryRow(query, number).Scan(&p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{},err
	}


	// заполните объект Parcel данными из таблицы

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	query:= "SELECT * FROM parcel WHERE client=?"
	// здесь из таблицы может вернуться несколько строк
	var res []Parcel
	rows,err := s.db.Query(query, client)
	if err != nil {
		return nil,err
	}
	defer rows.Close()
	// заполните срез Parcel данными из таблицы
	for rows.Next() {
		var p Parcel
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	query:="UPDATE parcel SET status=? WHERE number=?"
	_, err := s.db.Exec(query, status, number)
	if err!=nil{
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	query := "UPDATE parcel SET address = ? WHERE number = ? AND status = ?"
	result, err := s.db.Exec(query, address, number, ParcelStatusRegistered)
	if err!=nil{
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("parcel not found or not registered")
	}
	// менять адрес можно только если значение статуса registered
	

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	
	query:="DELETE FROM parcel WHERE number=? AND status=?"
	result,err:=s.db.Exec(query,number,ParcelStatusRegistered)
	if err!=nil{
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("parcel not found or not registered")
	}

	return nil
}
