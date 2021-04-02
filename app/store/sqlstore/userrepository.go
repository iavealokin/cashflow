package sqlstore

import (
	"database/sql"
	"errors"
	"log"
	//"fmt"
	"github.com/iavealokin/cashflow/app/model"
)

//UserRepository struct
type UserRepository struct {
	store *Store
}
//UserLogin ...
func (r *UserRepository) UserLogin(login string, password string) (*model.User, error){
	var errorDrop error
	u := new(model.User)
	sqlStatement := `
	select count(*) cnt,username, surname, userid from users
	WHERE userlogin = $1 and userpwd=$2
	group by username,surname,userid;`
var cnt int
	rows, err := r.store.db.Query(sqlStatement, login, password)
		if err != nil {
			panic(err)
		}
defer rows.Close()
for rows.Next(){
	err:=rows.Scan(&cnt,&u.Username,&u.Surname, &u.ID)
	if err !=nil{
		log.Fatal(err)
	}
}
		if cnt == 0 {
			errorDrop = errors.New("Incorrect login or password")
			
		}
		return u,errorDrop

}

//Create user... 
func (r *UserRepository) Create(o *model.Operation) (error) {

	return r.store.db.QueryRow(
		"INSERT INTO operations " +
			"(billing_id,user_id,amount,direction,operation_comment,operation_date) values($1,$2,$3,$4,$5,current_timestamp) RETURNING operation_id",
	"1",
	"1",
  	&o.Amount,
   	&o.Direction,
   	&o.Comment,
	).Scan(&o.ID)
}

//Drop user ...
func (r *UserRepository) Drop(u *model.User) (error){
	var errorDrop error
	sqlStatement := `
DELETE FROM users
WHERE id = $1;`
	if u.ID == 1 {
		errorDrop = errors.New("Permission denied - delete user Admin")
	} else {
		res, err := r.store.db.Exec(sqlStatement, u.ID)
		if err != nil {
			panic(err)
		}
		count, err := res.RowsAffected()
		if err != nil {
			panic(err)
		}
		if count == 0 {
			errorDrop = errors.New("No records delete")
		} 
	}
	
	return errorDrop
	}



//Update user ...
func (r *UserRepository) Update(u *model.User) (error){
	var errorDrop error
	var err error
	var res sql.Result
	if u.Password==""{
		sqlStatement := `
		UPDATE users
		SET login = $1,
		name = $2,
		surname = $3,
		birthday = $4
		where id = $5;`
			res, err = r.store.db.Exec(sqlStatement, u.Login, u.Username, u.Surname, u.Birthday, u.ID)
			if err != nil {
				panic(err)
			}
	}else{
		sqlStatement := `
		UPDATE users
		SET login = $1,
		name = $2,
		surname = $3,
		birthday = $4,
		password = $5
		where id = $6;`
			res, err = r.store.db.Exec(sqlStatement, u.Login, u.Username, u.Surname, u.Birthday, u.Password, u.ID)
			if err != nil {
				panic(err)
			}
	}
	
		count, err := res.RowsAffected()
		if err != nil {
			panic(err)
		}
		if count == 0 {
			errorDrop = errors.New("No records for update")
		} 
		return errorDrop
	}
	//Get users list
	func (r *UserRepository) GetOperations(userid int) ([] model.Operation, error){
		
		sqlStatement := `select 
		coalesce(to_char(operation_id,'99999999'),' ') oper_id
		, coalesce(to_char(amount,'9999999.99'),'0.00') amount
		, case when trim(coalesce(to_char(direction,'99'),' '))='2' then 'Расход' else 'Доход' end direction
		, coalesce(operation_comment,' ')
		, coalesce(to_char(operation_date, 'YYYY-MM-DD HH24:MI:SS'),' ') as operation_date
		from operations where cast(operation_date as date)=cast(current_timestamp as date)
		and user_id=$1
		union all
		select '',
case when plus>minus then to_char(plus- minus,'999999999.99')
when plus<minus then '-'||to_char(minus-plus,'99999999999.99')
when coalesce(plus,0)=coalesce(minus,0) then '0.00'
end
,case when plus =0 and minus =0
then '0.00'
when plus = 0 and minus>0
then '0.00/'||to_char(minus,'999999999.99')
else to_char(plus,'99999999.99')||'/'||to_char(coalesce(minus,0),'999999999.99')
end
, 'Итого',''
		from
		(
select coalesce(sum(amount),0) plus,(select coalesce(sum(amount),0)
from operations
where cast(operation_date as date)=cast(current_timestamp as date) and direction=2 and user_id=$1)minus
from
		operations
		where cast(operation_date as date)=cast(current_timestamp as date) and direction=1 and user_id=$1
		)t;`
		var trns [] model.Operation
		rows,err := r.store.db.Query(sqlStatement,userid)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()
			for rows.Next() {
				trn := new(model.Operation)
				err := rows.Scan(&trn.ID, &trn.Amount, &trn.Direction, &trn.Comment, &trn.Date)
				if err != nil {
					log.Fatal(err)
				}
				trnarr := model.Operation{ID: trn.ID, Amount: trn.Amount, Direction: trn.Direction, Comment: trn.Comment, Date: trn.Date}
				trns = append(trns, trnarr)
			}
			if err = rows.Err(); err != nil {
				log.Fatal(err)
			}
			return trns, err
		}

//Get balance, limits,cashflow,etc..
		func (r *UserRepository) GetUserData(userid int) (*model.UserData, error){
			ud := new(model.UserData)
			sqlStatement := `select coalesce(income,0),coalesce(outcome,0),coalesce(income,0)-coalesce(outcome,0), case when income-outcome>0 then 1 else 0 end flag from
			(
			select sum(coalesce(amount,0)) income,(select sum(coalesce(amount,0)) from operations 
			where user_id=$1 and date_trunc('month',operation_date)=date_trunc('month',current_timestamp) and direction=2) outcome
			from operations where user_id=$1 and date_trunc('month',operation_date)=date_trunc('month',current_timestamp) and direction=1
			)t;`

			rows,err := r.store.db.Query(sqlStatement,userid)
				if err != nil {
					log.Fatal(err)
				}
				defer rows.Close()
				for rows.Next() {
					//ud := new(model.UserData)
					err := rows.Scan(&ud.Income, &ud.Outcome, &ud.Difference, &ud.Flag)
					if err != nil {
						log.Fatal(err)
					}
									}
				if err = rows.Err(); err != nil {
					log.Fatal(err)
				}

				actives := GetActives(userid,r)
				passives := GetPassives(userid,r)


				ud.Actives=actives
				ud.Passives=passives
				return ud, err
			}
	


func GetActives(userid int, r *UserRepository) []model.UserActive{
		sqlStatementActives := `select 
id                 
 ,  active_name     
 ,  cost           
 ,  amount        
 ,  result         
 ,  percent         
 ,  user_id 
from actives where user_id=$1;`
	var actives []model.UserActive
				rows,err := r.store.db.Query(sqlStatementActives,userid)
				if err != nil {
					log.Fatal(err)
				}
				defer rows.Close()
				for rows.Next() {
					active := new(model.UserActive)
					err := rows.Scan(&active.ID, &active.Name, &active.Cost, &active.Amount,&active.Result,&active.Percent,&active.UserID)
					if err != nil {
						log.Fatal(err)
					}
					activearr := model.UserActive{ID: active.ID, Name: active.Name, Cost: active.Cost, Amount: active.Amount,Result: active.Result,Percent:active.Percent,UserID: active.UserID}
					actives =append(actives,activearr)
									}
				if err = rows.Err(); err != nil {
					log.Fatal(err)
				}
				rows,err = r.store.db.Query("select to_char(coalesce(sum(result),0.00),'99999999999.99') from actives where user_id=$1",userid)
				if err != nil {
					log.Fatal(err)
				}
				acv := new(model.UserActive)
				for rows.Next() {

					err := rows.Scan(&acv.Sum)
					if err != nil {

						log.Fatal(err)
					}
					}
				if err = rows.Err(); err != nil {
					log.Fatal(err)
				}

				actives[0].Sum=acv.Sum
				return actives
}
func GetCategories(userid int, r *UserRepository) []model.Category {

}

func GetPassives(userid int, r *UserRepository) []model.UserPassive{
		sqlStatementPassives := `select 
id                 
 ,  passive_name     
 ,  cost           
 ,  amount        
 ,  result         
 ,  percent         
 ,  user_id 
from passives where user_id=$1;`
	var passives []model.UserPassive
				rows,err := r.store.db.Query(sqlStatementPassives,userid)
				if err != nil {
					log.Fatal(err)
				}
				defer rows.Close()
				for rows.Next() {
					passive := new(model.UserPassive)
					err := rows.Scan(&passive.ID, &passive.Name, &passive.Cost, &passive.Amount,&passive.Result,&passive.Percent,&passive.UserID)
					if err != nil {
						log.Fatal(err)
					}
					passivearr := model.UserPassive{ID: passive.ID, Name: passive.Name, Cost: passive.Cost, Amount: passive.Amount,Result: passive.Result,Percent:passive.Percent,UserID: passive.UserID}
					passives =append(passives,passivearr)
									}
				if err = rows.Err(); err != nil {
					log.Fatal(err)
				}
			rows,err = r.store.db.Query("select to_char(coalesce(sum(result),0.00),'99999999999.99') from passives where user_id=$1",userid)
				if err != nil {
					log.Fatal(err)
				}
				psv := new(model.UserPassive)
				for rows.Next() {

					err := rows.Scan(&psv.Sum)
					if err != nil {

						log.Fatal(err)
					}
					}
				if err = rows.Err(); err != nil {
					log.Fatal(err)
				}
				passives[0].Sum=psv.Sum
				return passives
}