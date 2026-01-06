package data

import (
	"context"
	"database/sql"
	"time"
    "golang.org/x/crypto/bcrypt"
    "errors"
)

type Associate struct {
	ID             int       `json:"id"`
	FirstName      string    `json:"FirstName"`
	LastName       string    `json:"LastName"`
	Title          string    `json:"Title"`
	Department     string    `json:"Department"`
	Office         string    `json:"Office"`
	Status         string    `json:"Status"`
	StartDate      time.Time `json:"StartDate"`
	EmplStatus     string    `json:"EmplStatus"`
	Salary         int       `json:"Salary"`
	DOB            time.Time `json:"DOB"`
	ProfilePicture string    `json:"profile_picture"`
    Password       string    `json:"-"` // Don't expose password in JSON
    Email          string    `json:"Email"`
    PhoneNumber    string    `json:"PhoneNumber"`
    Gender         string    `json:"Gender"`
    PrivateEmail   string    `json:"PrivateEmail"`
    ManagerID      *int      `json:"manager_id"`
}

type Models struct {
	Associates         AssociateModel
	Offices            OfficeModel
	Departments        DepartmentModel
	Tasks              TaskModel
	Thanks             ThankModel
	TimeOffRequests    TimeOffRequestModel
	DocumentCategories DocumentCategoryModel
	MenuPermissions    MenuPermissionModel
	AppSettings        AppSettingsModel
	TimeEntries        TimeEntryModel
	ThanksCategories   ThanksCategoryModel
	Holidays           HolidayModel
}

type AssociateModel struct {
	DB *sql.DB
}

type OfficeModel struct {
	DB *sql.DB
}

type DepartmentModel struct {
	DB *sql.DB
}

func New(db *sql.DB) Models {
	return Models{
		Associates:         AssociateModel{DB: db},
		Offices:            OfficeModel{DB: db},
		Departments:        DepartmentModel{DB: db},
		Tasks:              TaskModel{DB: db},
		Thanks:             ThankModel{DB: db},
		TimeOffRequests:    TimeOffRequestModel{DB: db},
		DocumentCategories: DocumentCategoryModel{DB: db},
		MenuPermissions:    MenuPermissionModel{DB: db},
		AppSettings:        AppSettingsModel{DB: db},
		TimeEntries:        TimeEntryModel{DB: db},
		ThanksCategories:   ThanksCategoryModel{DB: db},
		Holidays:           HolidayModel{DB: db},
	}
}

func (m AssociateModel) Insert(associate Associate) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(associate.Password), 12)
    if err != nil {
        return 0, err
    }

	stmt := `INSERT INTO Associates (first_name, last_name, title, department, office, status, start_date, empl_status, salary, dob, profile_picture, password, email, phone_number, gender, private_email, manager_id)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := m.DB.ExecContext(ctx, stmt,
		associate.FirstName,
		associate.LastName,
		associate.Title,
		associate.Department,
		associate.Office,
		associate.Status,
		associate.StartDate,
		associate.EmplStatus,
		associate.Salary,
		associate.DOB,
		associate.ProfilePicture,
        hashedPassword,
        associate.Email,
        associate.PhoneNumber,
        associate.Gender,
        associate.PrivateEmail,
        associate.ManagerID,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m DepartmentModel) Delete(id int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    stmt := `DELETE FROM Departments WHERE id = ?`
    _, err := m.DB.ExecContext(ctx, stmt, id)
    return err
}

func (m AssociateModel) GetByEmail(email string) (*Associate, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    query := `select id, first_name, last_name, COALESCE(email, ''), COALESCE(password, ''), title, department, COALESCE(profile_picture, ''), start_date, manager_id from Associates where email = ?`
    var associate Associate
    
    row := m.DB.QueryRowContext(ctx, query, email)
    err := row.Scan(
        &associate.ID,
        &associate.FirstName,
        &associate.LastName,
        &associate.Email,
        &associate.Password,
        &associate.Title,
        &associate.Department,
        &associate.ProfilePicture,
        &associate.StartDate,
        &associate.ManagerID,
    )
    
    if err != nil {
        return nil, err
    }
    
    return &associate, nil
}

func (m AssociateModel) PasswordMatches(plainText string, associate Associate) (bool, error) {
    err := bcrypt.CompareHashAndPassword([]byte(associate.Password), []byte(plainText))
    if err != nil {
        switch {
        case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
             return false, nil
        default:
             return false, err
        }
    }
    return true, nil
}

func (m AssociateModel) GetOne(id int) (*Associate, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    query := `SELECT id, first_name, last_name, title, department, office, status, start_date, empl_status, salary, dob, COALESCE(profile_picture, ''), COALESCE(email, ''), COALESCE(password, ''), COALESCE(phone_number, ''), COALESCE(gender, ''), COALESCE(private_email, ''), manager_id
    FROM Associates WHERE id = ?`

    var a Associate
    row := m.DB.QueryRowContext(ctx, query, id)
    err := row.Scan(
        &a.ID,
        &a.FirstName,
        &a.LastName,
        &a.Title,
        &a.Department,
        &a.Office,
        &a.Status,
        &a.StartDate,
        &a.EmplStatus,
        &a.Salary,
        &a.DOB,
        &a.ProfilePicture,
        &a.Email,
        &a.Password,
        &a.PhoneNumber,
        &a.Gender,
        &a.PrivateEmail,
        &a.ManagerID,
    )

    if err != nil {
        return nil, err
    }
    
    return &a, nil
}

func (m AssociateModel) GetAll() ([]Associate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, first_name, last_name, title, department, office, status, start_date, empl_status, salary, dob, COALESCE(profile_picture, ''), COALESCE(email, ''), COALESCE(phone_number, ''), COALESCE(gender, ''), COALESCE(private_email, ''), manager_id
	FROM Associates ORDER BY last_name`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var associates []Associate
	for rows.Next() {
		var a Associate
		err := rows.Scan(
			&a.ID,
			&a.FirstName,
			&a.LastName,
			&a.Title,
			&a.Department,
			&a.Office,
			&a.Status,
			&a.StartDate,
			&a.EmplStatus,
			&a.Salary,
			&a.DOB,
			&a.ProfilePicture,
            &a.Email,
            &a.PhoneNumber,
            &a.Gender,
            &a.PrivateEmail,
            &a.ManagerID,
		)
		if err != nil {
			return nil, err
		}
		associates = append(associates, a)
	}

	return associates, nil
}

func (m AssociateModel) Update(id int, associate Associate) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    var query string
    var args []interface{}
    
    if associate.Password != "" {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(associate.Password), 12)
        if err != nil {
            return err
        }
        query = `UPDATE Associates SET first_name=?, last_name=?, title=?, department=?, office=?, status=?, start_date=?, empl_status=?, salary=?, dob=?, email=?, password=?, phone_number=?, gender=?, private_email=?, manager_id=? WHERE id=?`
        args = []interface{}{
            associate.FirstName, associate.LastName, associate.Title, associate.Department, associate.Office, associate.Status, 
            associate.StartDate, associate.EmplStatus, associate.Salary, associate.DOB, associate.Email, hashedPassword, 
            associate.PhoneNumber, associate.Gender, associate.PrivateEmail, associate.ManagerID, id,
        }
    } else {
        query = `UPDATE Associates SET first_name=?, last_name=?, title=?, department=?, office=?, status=?, start_date=?, empl_status=?, salary=?, dob=?, email=?, phone_number=?, gender=?, private_email=?, manager_id=? WHERE id=?`
        args = []interface{}{
            associate.FirstName, associate.LastName, associate.Title, associate.Department, associate.Office, associate.Status, 
            associate.StartDate, associate.EmplStatus, associate.Salary, associate.DOB, associate.Email, 
            associate.PhoneNumber, associate.Gender, associate.PrivateEmail, associate.ManagerID, id,
        }
    }

    _, err := m.DB.ExecContext(ctx, query, args...)
    return err
}

func (m AssociateModel) UpdatePassword(id int, password string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil {
        return err
    }

    query := `UPDATE Associates SET password = ? WHERE id = ?`
    _, err = m.DB.ExecContext(ctx, query, hashedPassword, id)
    return err
}

func (m AssociateModel) Delete(id int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    query := `DELETE FROM Associates WHERE id = ?`
    _, err := m.DB.ExecContext(ctx, query, id)
    return err
}

type Office struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Department struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (m OfficeModel) GetAll() ([]Office, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT id, name FROM Offices ORDER BY name`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var offices []Office
	for rows.Next() {
		var o Office
		if err := rows.Scan(&o.ID, &o.Name); err != nil {
			return nil, err
		}
		offices = append(offices, o)
	}
	return offices, nil
}

func (m OfficeModel) Insert(office Office) (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    stmt := `INSERT INTO Offices (name) VALUES (?)`
    result, err := m.DB.ExecContext(ctx, stmt, office.Name)
    if err != nil {
        return 0, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }

    return int(id), nil
}

func (m OfficeModel) Delete(id int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    stmt := `DELETE FROM Offices WHERE id = ?`
    _, err := m.DB.ExecContext(ctx, stmt, id)
    return err
}

func (m DepartmentModel) GetAll() ([]Department, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT id, name FROM Departments ORDER BY name`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var departments []Department
	for rows.Next() {
		var d Department
		if err := rows.Scan(&d.ID, &d.Name); err != nil {
			return nil, err
		}
		departments = append(departments, d)
	}
	return departments, nil
}

func (m DepartmentModel) Insert(dept Department) (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    stmt := `INSERT INTO Departments (name) VALUES (?)`
    result, err := m.DB.ExecContext(ctx, stmt, dept.Name)
    if err != nil {
        return 0, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }

    return int(id), nil
}



