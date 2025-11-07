package dbrepo

import (
	"context"
	"errors"
	"time"

	"githup.com/Mo277210/booking/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgreDBRepo) AllUsers() bool {
	return true
}
//InsertReservation inserts a reservation into the database
func (m *postgreDBRepo) InsertReservation(res models.Reservation)(int,error){

	ctx,cancel:=context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

	var newID int

stmt := `insert into reservations (first_name, last_name,email,phone,start_date,
end_date,room_id,created_at,updated_at) 
values ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`

err:=m.DB.QueryRowContext(ctx,
	stmt,
	res.FirstName,
	res.LastName,
	res.Email,
	res.Phone,
	res.StartDate,
	res.EndDate,
	res.RoomID,
	time.Now(),
	time.Now(),

).Scan(&newID)
if err!=nil{
	return 0, err
}


	return newID,nil
}




//InsertRoomRestriction inserts a room restriction into the database

func(m *postgreDBRepo)InsertRoomRestriction(r models.RoomRestrictions)error{
	
	
	ctx,cancel:=context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()


stmt:= `
insert into room_restrictions (start_date,end_date,room_id,
reservations_id,restrictions_id,created_at,updated_at)
values ($1,$2,$3,$4,$5,$6,$7)
`




	_,err:=m.DB.ExecContext(
		ctx,stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now(),
		// r.RestrictionID,
	)

	if err!=nil{
		return err
	}
	return nil
}


// SearchAvailabilityByDatesByRoomID returns true if availability exists, and false if no availability
func(m *postgreDBRepo)SearchAvailabilityByDatesByRoomID(start,end time.Time,roomID int)(bool,error){

	ctx,cancel:=context.WithTimeout(context.Background(),3*time.Second)
	
   defer cancel()
	
   var numRows int
  
   query:=`
select
	 
	count(id) 
 from
	  room_restrictions 
 where
	  
	  room_id=$1 
	  and $2 < end_date and $3 > start_date`

	  // row:=m.DB.QueryRowContext(ctx,query,start,end)
	
	 row:=m.DB.QueryRowContext(ctx,query,roomID,start,end)
	 
	 err:=row.Scan(&numRows)
	 if err!=nil{
		 return false,err
	 }

	 if numRows==0{
		 return true,nil
	 }
	return false,nil
}

//SearchAvailabilityForAllRooms returns a slice of available rooms
func(m *postgreDBRepo)SearchAvailabilityForAllRooms(start,end time.Time)([]models.Room, error){
	
	
	ctx,cancel:=context.WithTimeout(context.Background(),3*time.Second)
	
   defer cancel()

var rooms []models.Room

	query := `
select 
   r.id, r.rooms_name
from 
   rooms r
where 
   r.id not in (
      select room_id 
      from room_restrictions rr 
      where $1 < rr.end_date 
      and $2 > rr.start_date
   )`

	rows,err:=m.DB.QueryContext(ctx,query,start,end)

	if err!=nil{
		return rooms,err
	}

	for rows.Next(){
		var room models.Room
		err:=rows.Scan(&room.ID,&room.RoomName)
		if err!=nil{
			return rooms,err
		}
		rooms=append(rooms,room)
	}

	
	return rooms, nil
}

//GetRoomByID gets a room by id
func(m *postgreDBRepo)GetRoomByID(id int) (models.Room,error){
	
	ctx,cancel:=context.WithTimeout(context.Background(),3*time.Second)
	
   defer cancel()

var room models.Room

	query := `
select 
   id,rooms_name,created_at,updated_at
from 
   rooms 
where 
   id = $1`

   row:=m.DB.QueryRowContext(ctx,query,id)

   err:=row.Scan(&room.ID,&room.RoomName,&room.CreatedAt,&room.UpdatedAt)
   if err!=nil{
		return room,err
   }

	return room, nil
}

//GetUserByID gets a user by id
func (m *postgreDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT id, first_name, last_name, email, password, access_level, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreateAt,
		&u.UpdateAt,
	)
	if err != nil {
		return models.User{}, err
	}

	return u, nil
}

//UpdateUser updates a user in the database
func (m *postgreDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
	UPDATE users
	SET first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5
	
	`
_,err:=m.DB.ExecContext(ctx,query,
u.FirstName,
u.LastName,
u.Email,
u.AccessLevel,
time.Now(),
)
if err!=nil{
	return err
}
return nil

}

//Authenticate authenticates a user
func (m *postgreDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var id int
	var hashedPassword string
	
	query := `
	SELECT id, password
	FROM users
	WHERE email = $1
	`
	row := m.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}
	return id, hashedPassword, nil
}

//AllReservations returns a slice of all reservations
func (m *postgreDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
	SELECT r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
	r.end_date, r.room_id, r.created_at, r.updated_at, rm.id,r.processed,
	rm.rooms_name
	FROM reservations r
	LEFT JOIN rooms rm ON (r.room_id = rm.id)
	ORDER BY r.start_date ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}

//AllNewReservations returns a slice of all new reservations
func (m *postgreDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
	SELECT r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
	r.end_date, r.room_id, r.created_at, r.updated_at, rm.id,
	rm.rooms_name
	FROM reservations r
	LEFT JOIN rooms rm ON (r.room_id = rm.id)
	WHERE processed = 0
	ORDER BY r.start_date ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}
	return reservations, nil
}

//GetReservationByID gets a reservation by id
func (m *postgreDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var res models.Reservation

	query := `
	SELECT r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date,
	r.end_date, r.room_id, r.created_at, r.updated_at, r.processed, rm.id,
	rm.rooms_name
	FROM reservations r
	LEFT JOIN rooms rm ON (r.room_id = rm.id)
	WHERE r.id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Phone,
		&res.StartDate,
		&res.EndDate,
		&res.RoomID,
		&res.CreatedAt,
		&res.UpdatedAt,
		&res.Processed,
		&res.Room.ID,
		&res.Room.RoomName,
	)
	if err != nil {
		return res, err
	}
	return res, nil
}

//UpdateReservation updates a reservation in the database
func (m *postgreDBRepo) UpdateResertvation(u models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
	UPDATE reservations
	SET first_name = $1, last_name = $2, email = $3, phone = $4, updated_at = $5
	WHERE id = $6
	`
_,err:=m.DB.ExecContext(ctx,query,
u.FirstName,
u.LastName,
u.Email,
u.Phone,
time.Now(),
u.ID,
)
if err!=nil{
	return err
}
return nil

}

//DeleteReservation deletes a reservation from the database
func (m *postgreDBRepo) DeleteReservation(id int) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		query:= `
		delete from reservations where id=$1
		`
		_,err:=m.DB.ExecContext(ctx,query,id)
		if err!=nil{
			return err
		}
		return nil

}

//UpdateProcessedReservation updates a reservation's processed status
func (m *postgreDBRepo) UpdateProcessedReservation(id int, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	UPDATE reservations
	SET processed = $1
	WHERE id = $2
	`
	_, err := m.DB.ExecContext(ctx, query, processed, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgreDBRepo) AllRooms () ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room
	query := `
	select id, rooms_name, created_at, updated_at
	from rooms
	order by rooms_name
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var rm models.Room
		err := rows.Scan(
			&rm.ID,
			&rm.RoomName,
			&rm.CreatedAt,
			&rm.UpdatedAt,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, rm)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

//GetRestrictionsForRoomByDate returns a slice of room restrictions for a given room and date range
func (m *postgreDBRepo) GetRestrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestrictions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restrictions []models.RoomRestrictions
	query := `
	select id, coalesce(reservations_id, 0), restrictions_id, room_id, start_date, end_date
	from room_restrictions
	where start_date <= $2 and end_date >= $1
	and room_id = $3
	`

	rows, err := m.DB.QueryContext(ctx, query, start, end, roomID)
	if err != nil {
		return restrictions, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.RoomRestrictions
		err := rows.Scan(
			&r.ID,
			&r.ReservationID,
			&r.RestrictionID,
			&r.RoomID,
			&r.StartDate,
			&r.EndDate,
		)
		if err != nil {
			return nil, err
		}
		restrictions = append(restrictions, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return restrictions, nil
}

//InsertBlockForRoom inserts a room restriction
func (m *postgreDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	insert into room_restrictions (start_date, end_date, room_id, restrictions_id, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $6)
	`
	_, err := m.DB.ExecContext(ctx, query, startDate, startDate, id, 2, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}


//DeleteBlockByID deletes a room restriction
func (m *postgreDBRepo) DeleteBlockByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	delete from room_restrictions where id = $1
	`
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}