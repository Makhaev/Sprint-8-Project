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
	res, err := s.db.Exec(
		`INSERT INTO parsel (client,status,addres,created_at)
		VALUES(?,?,?,?)
		`,
		p.Client, p.Status, p.Address, p.CreatedAt,
	)

	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	row := s.db.QueryRow(`
		SELECT number, client, status, address, created_at
		FROM parcel
		WHERE number = ?`, number)

	var p Parcel
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	return p, nil

	// здесь из таблицы должна вернуться только одна строка

	// заполните объект Parcel данными из таблицы

}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query(`
		SELECT number, client, status, address, created_at
		FROM parcel
		WHERE client = ?`, client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Parcel
	for rows.Next() {
		var p Parcel
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	res, err := s.db.Exec("UPDATE parcel SET status = ? WHERE number = ?", status, number)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("посылка с номером %d не найдена", number)
	}

	return nil

}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	var status string
	err := s.db.QueryRow("SELECT status FROM parcel WHERE number = ?", number).Scan(&status)
	if err != nil {
		return fmt.Errorf("ошибка при получении статуса: %v", err)
	}
	// Проверяем статус
	if status != ParcelStatusRegistered {
		return fmt.Errorf("адрес можно изменить только при статусе 'registered'")
	}
	// Обновляем адрес
	_, err = s.db.Exec("UPDATE parcel SET address = ? WHERE number = ?", address, number)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении адреса: %v", err)
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	var status string
	err := s.db.QueryRow("SELECT status FROM parcel WHERE number = ?", number).Scan(&status)
	if err != nil {
		return fmt.Errorf("ошибка при получении статуса: %v", err)
	}

	if status != ParcelStatusRegistered {
		return fmt.Errorf("адрес можно изменить только при статусе 'registered'")
	}

	// Удаляем строку
	_, err = s.db.Exec("DELETE FROM parcel WHERE number = ?", number)
	if err != nil {
		return fmt.Errorf("ошибка при удалении: %v", err)
	}

	return nil
}
