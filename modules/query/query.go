// Package query is the collection of database facing functions used by
// zedlist.
package query

import (
	"fmt"
	sysLog "log"
	"time"

	"github.com/gernest/zedlist/modules/forms"
	"github.com/gernest/zedlist/modules/settings"
	"github.com/jinzhu/gorm"

	"github.com/drhodes/golorem"
	"github.com/gernest/zedlist/models"
	"github.com/oxtoacart/bpool"
	"golang.org/x/crypto/bcrypt"
)

//
//
//		JOB
//
//

// CreateJob saves the job j into database
func CreateJob(conn *gorm.DB, j *models.Job) error {
	return Create(conn, j)
}

// GetJobByID retrieves a job record from the database by id.
func GetJobByID(conn *gorm.DB, id int64) (*models.Job, error) {
	jb := &models.Job{}
	q := conn.First(jb, id)
	if q.Error != nil {
		return nil, q.Error
	}
	return jb, nil
}

// GetALLJobs returns an ordered slice of all jobs reocrds. The order is by creation
// date and id in descending order.
// TODO remove id?
func GetALLJobs(conn *gorm.DB) ([]*models.Job, error) {
	jbs := []*models.Job{}
	q := conn.Order("created_at desc").Find(&jbs)
	if q.Error != nil {
		return nil, q.Error
	}
	return jbs, nil
}

// GetLatestJobs returns latest jobs limited by settings.MazListLimit.
func GetLatestJobs(conn *gorm.DB) ([]*models.Job, error) {
	jobs := []*models.Job{}
	q := conn.Order("created_at desc,id").Limit(settings.MaxListLimit).Find(&jobs)
	if q.Error != nil {
		return nil, q.Error
	}
	return jobs, nil

}

//
//
//		MIGRATIONS
//
//

var bufPool = bpool.NewBufferPool(2)

func composeJob() *models.Job {
	j := &models.Job{}
	j.Title = lorem.Sentence(5, 10)
	buf := bufPool.Get()
	defer bufPool.Put(buf)

	for _ = range make([]struct{}, 2) {
		buf.WriteString(lorem.Paragraph(4, 5))
		buf.WriteString("\n")
	}
	j.Description = buf.String()
	return j
}

// PopulateDB polutates the database with dummy data( for tests only )
func PopulateDB(conn *gorm.DB) {
	sysLog.Print("Populating database...")
	var migrationData = []struct {
		name, short string
	}{
		{"mwanza", "mza"},
		{"dar es salaam", "dar"},
	}
	deadline := time.Now().Add(time.Hour)

	q := conn.First(&models.Region{})
	qj := conn.First(&models.Job{})
	qu := conn.First(&models.User{})
	if q.Error != nil && qj.Error != nil && qu.Error != nil {
		for _, v := range migrationData {
			conn.FirstOrCreate(&models.Region{}, models.Region{Name: v.name, Short: v.short})
		}
		regs := []models.Region{}
		conn.Find(&regs)
		for _ = range make([]struct{}, 40) {
			for _, v := range regs {
				j := composeJob()
				j.Deadline = deadline
				j.Region = v
				conn.Create(j)
			}
		}

		// create a sample user
		err := SampleUser(conn)
		if err != nil {
			sysLog.Println(err)
		}
	}
	sysLog.Printf("Done. \n")
}

// MigrateSession creates sessions database table if it does not exist.
func MigrateSession(conn *gorm.DB) {
	conn.AutoMigrate(&models.Session{})
}

// DropSession a helper for droping the sessions table. useful in tests.
func DropSession(conn *gorm.DB) {
	conn.DropTableIfExists(&models.Session{})
}

// SampleUser creates a sample user.
func SampleUser(conn *gorm.DB) error {
	// create a  sample user( to speed up development)
	regForm := forms.Register{
		Email:           "root@home.com",
		Password:        "superroot",
		ConfirmPassword: "superroot",
	}
	u, err := CreateNewUser(conn, regForm)
	if err != nil {
		return err
	}
	if u != nil {
		sysLog.Printf("%s", fmt.Sprintf("created sample user email %s password %s", u.Email, regForm.Password))
	}
	return nil
}

//
//
//		USER
//
//

// GetUserByID retrieves the user by the given id.
func GetUserByID(conn *gorm.DB, id int64) (*models.User, error) {
	usr := &models.User{}
	q := conn.First(usr, id)
	if q.Error != nil {
		return nil, q.Error
	}
	return usr, nil
}

// GetUserByEmail retrieves a user by email.
func GetUserByEmail(conn *gorm.DB, email string) (*models.User, error) {
	usr := &models.User{}
	q := conn.Where(&models.User{Email: email}).First(usr)
	if q.Error != nil {
		return nil, q.Error
	}
	return usr, nil
}

// GetUserByName retrieves a user by name.
func GetUserByName(conn *gorm.DB, name string) (*models.User, error) {
	usr := &models.User{}
	q := conn.Where(&models.User{Name: name}).First(usr)
	if q.Error != nil {
		return nil, q.Error
	}
	return usr, nil
}

// CreateNewUser creates a new user.
func CreateNewUser(conn *gorm.DB, reg forms.Register) (*models.User, error) {
	_, err := GetUserByEmail(conn, reg.Email)
	if err == nil {
		return nil, fmt.Errorf("email %s already taken", reg.Email)
	}
	hashedPass, err := hashPassword(reg.Password)
	if err != nil {
		return nil, err
	}
	usr := &models.User{
		Email:    reg.Email,
		Password: hashedPass,
		Name:     reg.UserName,
		Status:   models.StatusActive,
		Person: models.Person{
			Email:      reg.Email,
			ObjectType: models.ObjPerson,
			PersonName: models.PersonName{
				Name: reg.UserName,
			},
		},
	}
	err = Create(conn, usr)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

// AuthenticateUserByEmail checks if the user which matches the loginForm details exist, and is valid.
func AuthenticateUserByEmail(conn *gorm.DB, loginForm forms.Login) (*models.User, error) {
	usr, err := GetUserByEmail(conn, loginForm.Name)
	if err != nil {
		return nil, err
	}
	err = verifyPass(usr.Password, loginForm.Password)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

// AuthenticateUserByEmail checks if the user which matches the loginForm details exist, and is valid.
func AuthenticateUserByName(conn *gorm.DB, loginForm forms.Login) (*models.User, error) {
	usr, err := GetUserByName(conn, loginForm.Name)
	if err != nil {
		return nil, err
	}
	err = verifyPass(usr.Password, loginForm.Password)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

// Verfiess if the given has mathes the password. The hash must be a bcrypt encoded has.
// it uses bcrypt to compare the two passwords
func verifyPass(hash, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}

// Encrypts a given string using bcrypt library. It returns the hashed password as a string,
// or any error
func hashPassword(pass string) (string, error) {
	np, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(np), err
}

//
//
//		PERSON
//
//

// GetPersonByUserID retrieves a person by using a user's ID.
func GetPersonByUserID(conn *gorm.DB, userID int64) (*models.Person, error) {
	p := &models.Person{}
	usr := &models.User{}
	qu := conn.First(usr, userID)
	if qu.Error != nil {
		return nil, qu.Error
	}
	q := conn.Preload("PersonName").First(p, usr.PersonID)
	if q.Error != nil {
		return nil, q.Error
	}
	return p, nil
}

// PersonCreateJob creates a job associated with the person p using the given jobForm j.
func PersonCreateJob(conn *gorm.DB, p *models.Person, j forms.JobForm) error {
	job := models.Job{
		Title:       j.Title,
		Description: j.Description,
	}
	q := conn.Model(p).Association("Jobs").Append(&job)
	return q.Error
}

// PersonDeleteJob deletes the job with ID jobID whichz is associated with person p.
func PersonDeleteJob(conn *gorm.DB, p *models.Person, jobID int) error {
	job := &models.Job{}
	q := conn.First(job, jobID)
	if q.Error != nil {
		return q.Error
	}
	qn := conn.Model(p).Association("jobs").Delete(job)
	return qn.Error
}

//
//
//		REGION
//
//

// GetAllRegions retrieves all regions from the database.
func GetAllRegions(conn *gorm.DB) ([]*models.Region, error) {
	regs := []*models.Region{}
	q := conn.Model(&models.Region{}).Find(&regs)
	if q.Error != nil {
		return nil, q.Error
	}
	return regs, nil
}

// GetJobByRegionShort retrieved jobs from a given region, the region short representation
// is used eg mza( for mwanza).
func GetJobByRegionShort(conn *gorm.DB, short string) ([]*models.Job, int64, error) {
	reg := &models.Region{}
	q := conn.Where(&models.Region{Short: short}).First(reg)
	if q.Error != nil {
		return nil, 0, q.Error
	}
	var count int64
	jobs := []*models.Job{}
	qerr := conn.Model(&models.Job{}).Where(&models.Job{RegionID: reg.ID}).Count(&count).
		Order("created_at desc,id").Limit(settings.MaxListLimit).Find(&jobs)

	if qerr.Error != nil {
		return nil, 0, qerr.Error
	}
	return jobs, count, nil
}

// GetJobByRegionPaginate retrieves jobs by a given region short name, offsetting at offset
// and limiting up to limit.
func GetJobByRegionPaginate(conn *gorm.DB, short string, offset, limit int) ([]*models.Job, error) {
	reg := &models.Region{}
	q := conn.Where(&models.Region{Short: short}).First(reg)
	if q.Error != nil {
		return nil, q.Error
	}
	jobs := []*models.Job{}
	qerr := conn.Model(&models.Job{}).Where(&models.Job{RegionID: reg.ID}).
		Order("created_at desc,id").Offset(offset).Limit(limit).Find(&jobs)

	if qerr.Error != nil {
		return nil, qerr.Error
	}
	return jobs, nil
}

//
//
//		SESSION
//
//

// GetSessionByKey retrieves the session by key. the key field is the session.ID for the
// gorilla/sessions Session.
func GetSessionByKey(conn *gorm.DB, key string) (*models.Session, error) {
	sess := &models.Session{}
	querry := conn.Where(&models.Session{Key: key}).First(sess)
	if querry.Error != nil {
		return nil, querry.Error
	}
	return sess, nil
}

// UpdateSession updates the databse record for the given session.
func UpdateSession(conn *gorm.DB, sess *models.Session) error {
	ss, err := GetSessionByKey(conn, sess.Key)
	if err != nil {
		return err
	}
	ss.Data = sess.Data
	return Update(conn, ss)
}

// DeleteSession deletes session record from database.
func DeleteSession(conn *gorm.DB, key string) error {
	s, err := GetSessionByKey(conn, key)
	if err != nil {
		return err
	}
	return conn.Delete(s).Error
}

//
//
//		UTILS
//
//

// Create creates a database record for the struct v.
func Create(conn *gorm.DB, v interface{}) error {
	q := conn.Create(v)
	return q.Error
}

// Update update any registered object in the database. Note, v should be a gorm compliant
// struct, which has already been migrated.
func Update(conn *gorm.DB, v interface{}) error {
	s := conn.Save(v)
	return s.Error
}

//Delete deletes a value v fr©om the database.
func Delete(conn *gorm.DB, v interface{}) error {
	d := conn.Delete(v)
	return d.Error
}

//
//
//		TOKEN
//
//

// GetTokenByKey retrieves a token by the given key.
func GetTokenByKey(conn *gorm.DB, key string) (*models.Token, error) {
	tk := &models.Token{}
	q := conn.Where(&models.Token{Key: key}).First(tk)
	if q.Error != nil {
		return nil, q.Error
	}
	return tk, nil
}

//
//
//		RESUME
//
//

// GetResumeByID retrieves the resume by ID.
func GetResumeByID(conn *gorm.DB, id int64) (*models.Resume, error) {
	resume := &models.Resume{}
	q := conn.Preload("ResumeBasic").First(resume, id)
	if q.Error != nil {
		return nil, q.Error
	}
	return resume, nil
}

//GetAllPersonResumes retruns all resumes velonging to person.
func GetAllPersonResumes(conn *gorm.DB, p *models.Person) ([]models.Resume, error) {
	rst := []models.Resume{}
	q := conn.Model(p).Related(&rst)
	if q.Error != nil {
		return nil, q.Error
	}
	return rst, nil
}

//CreateResume creates a new resume record
func CreateResume(conn *gorm.DB, p *models.Person, r models.Resume) error {
	p.Resumes = append(p.Resumes, r)
	return Update(conn, p)
}

func DeleteUser(conn *gorm.DB, id int) error {
	usr := &models.User{}
	q := conn.Preload("Person").
		First(usr, id)
	if q.Error != nil {
		return q.Error
	}
	return Delete(conn, usr)
}
